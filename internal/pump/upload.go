package pump

import (
	"github.com/vmihailenco/msgpack/v5"
	"helloworld/internal/pump/analytics"
	log "helloworld/pkg/logger"
	"sync"
	"sync/atomic"
	"time"
)

//const anaylticsKeyName = "job-analytics"

// UploadOptions contains configuration items related to storage.
type UploadOptions struct {
	WorkersNum              int           `json:"workers-num"                 mapstructure:"pool-size"`
	RecordsBufferSize       uint64        `json:"records-buffer-size"       mapstructure:"records-buffer-size"`
	FlushInterval           time.Duration `json:"flush-interval"            mapstructure:"flush-interval"`
	StorageExpirationTime   time.Duration `json:"storage-expiration-time"   mapstructure:"storage-expiration-time"`
	Enable                  bool          `json:"enable"                    mapstructure:"enable"`
	EnableDetailedRecording bool          `json:"enable-detailed-recording" mapstructure:"enable-detailed-recording"`
}

type UploadIns struct {
	handler                    UploadHandler
	workerNums                 int
	recordChan                 chan *analytics.AnalyticsRecord
	workerBufferSize           uint64
	recordsBufferFlushInterval time.Duration
	shouldStop                 uint32
	waitGroup                  sync.WaitGroup
}

// NewUploadIns returns a new storage instance.
func NewUploadIns(opts *UploadOptions, handler UploadHandler) *UploadIns {
	wn := opts.WorkersNum
	recordBufferSize := opts.RecordsBufferSize
	workerBufferSize := recordBufferSize / uint64(wn)
	log.Info().Uint64("workerBufferSize", workerBufferSize).Msgf("analytics pool worker buffer size")
	recordsChan := make(chan *analytics.AnalyticsRecord, recordBufferSize)
	uploadIns := &UploadIns{
		handler:                    handler,
		workerNums:                 wn,
		recordChan:                 recordsChan,
		workerBufferSize:           workerBufferSize,
		recordsBufferFlushInterval: opts.FlushInterval,
	}
	return uploadIns
}

func (u *UploadIns) Start() {
	atomic.SwapUint32(&u.shouldStop, 0)
	for i := 0; i < u.workerNums; i++ {
		u.waitGroup.Add(1)
		workerIs := i
		go u.runWorker(workerIs)
	}
}

func (u *UploadIns) runWorker(workerId int) {
	defer u.waitGroup.Done()
	log.Debug().Int("goroutinue id", workerId).Msg("worker running")

	recordsBuffer := make([][]byte, 0, u.workerBufferSize)
	lastSentTS := time.Now()
	for {
		var readyToSend bool
		select {
		case record, ok := <-u.recordChan:
			// channel 被关闭
			if !ok {
				u.handler.AppendToSetPipelined(anaylticsKeyName, recordsBuffer)
				return
			}

			// 有新的数据到达
			if encoded, err := msgpack.Marshal(record); err != nil {
				log.Err(err).Msg("error encoding analytics data")
			} else {
				recordsBuffer = append(recordsBuffer, encoded)
			}

			// 检查 buffer 里的内容是否可以发送
			readyToSend = uint64(len(recordsBuffer)) == u.workerBufferSize
		case <-time.After(u.recordsBufferFlushInterval):
			// 只要时间到达就发送
			readyToSend = true
		}
		// 发送数据并重置 buffer
		if len(recordsBuffer) > 0 && (readyToSend || time.Since(lastSentTS) >= u.recordsBufferFlushInterval) {
			u.handler.AppendToSetPipelined(anaylticsKeyName, recordsBuffer)
			recordsBuffer = recordsBuffer[:0]
			lastSentTS = time.Now()
		}
	}

}

type UploadHandler interface {
	Connect() bool
	AppendToSetPipelined(string, [][]byte)
}
