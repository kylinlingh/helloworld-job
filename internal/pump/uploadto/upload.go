package uploadto

import (
	"github.com/dgraph-io/ristretto"
	"github.com/vmihailenco/msgpack/v5"
	"helloworld/internal/pump/analytics"
	log "helloworld/pkg/logger"
	"sync"
	"sync/atomic"
	"time"
)

const anaylticsKeyName = "job-analytics"

var uploadIns *UploadIns

// UploadOptions defines options for upload-storage
type UploadOptions struct {
	Enable            bool          `json:"enable"                    mapstructure:"enable"`
	Storage           string        `json:"storage"                    mapstructure:"storage"`
	WorkersNum        int           `json:"workers-num"                 mapstructure:"pool-size"`
	RecordsBufferSize uint64        `json:"records-buffer-size"       mapstructure:"records-buffer-size"`
	FlushInterval     time.Duration `json:"flush-interval"            mapstructure:"flush-interval"`
	//StorageExpirationTime   time.Duration `json:"storage-expiration-time"   mapstructure:"storage-expiration-time"`
	EnableDetailedRecording bool `json:"enable-detailed-recording" mapstructure:"enable-detailed-recording"`
}

type UploadIns struct {
	uploadStorage              UploadStorage
	localCache                 *ristretto.Cache
	workerNums                 int
	recordChan                 chan *analytics.AnalyticsRecord
	workerBufferSize           uint64
	recordsBufferFlushInterval time.Duration
	shouldStop                 uint32
	waitGroup                  sync.WaitGroup
}

// NewUploadIns returns a new downloadfrom instance.
func NewUploadIns(opts *UploadOptions, handler UploadStorage) *UploadIns {
	wn := opts.WorkersNum
	recordBufferSize := opts.RecordsBufferSize
	workerBufferSize := recordBufferSize / uint64(wn)
	log.Info().Uint64("workerBufferSize", workerBufferSize).Msgf("analytics pool worker buffer size")
	recordsChan := make(chan *analytics.AnalyticsRecord, recordBufferSize)
	uploadIns = &UploadIns{
		uploadStorage:              handler,
		workerNums:                 wn,
		recordChan:                 recordsChan,
		workerBufferSize:           workerBufferSize,
		recordsBufferFlushInterval: opts.FlushInterval,
	}
	return uploadIns
}

func (u *UploadIns) Start() {
	u.uploadStorage.Connect()

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
				u.uploadStorage.AppendToSetPipelined(anaylticsKeyName, recordsBuffer)
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
			log.Trace().Msg("record can be uploaded")
			u.uploadStorage.AppendToSetPipelined(anaylticsKeyName, recordsBuffer)
			recordsBuffer = recordsBuffer[:0]
			lastSentTS = time.Now()
		}
	}

}

func (u *UploadIns) UploadRecord(record *analytics.AnalyticsRecord) error {
	// 检查信号
	if atomic.LoadUint32(&u.shouldStop) > 0 {
		return nil
	}

	u.recordChan <- record
	log.Trace().Msg("record has been uploaded")
	return nil
}

func GetUploadInstance() *UploadIns {
	return uploadIns
}

type UploadStorage interface {
	Connect() bool
	AppendToSetPipelined(string, [][]byte)
	GetStorage() UploadStorage
}
