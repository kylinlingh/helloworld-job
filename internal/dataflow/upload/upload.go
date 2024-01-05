package upload

import (
	"github.com/dgraph-io/ristretto"
	"github.com/vmihailenco/msgpack/v5"
	"helloworld/internal/dataflow/datastructure"
	"helloworld/internal/dataflow/storage/uploadto"
	log "helloworld/pkg/logger"
	"sync"
	"sync/atomic"
	"time"
)

const anaylticsKeyName = "job-analytics"

// UploadConfig defines options for uploadto-storage
type UploadConfig struct {
	Enable            bool          `json:"enable"                    mapstructure:"enable"`
	Storage           string        `json:"storage"                    mapstructure:"storage"`
	WorkersNum        int           `json:"workers-num"                 mapstructure:"pool-size"`
	RecordsBufferSize uint64        `json:"records-buffer-size"       mapstructure:"records-buffer-size"`
	FlushInterval     time.Duration `json:"flush-interval"            mapstructure:"flush-interval"`
	//StorageExpirationTime   time.Duration `json:"storage-expiration-time"   mapstructure:"storage-expiration-time"`
	EnableDetailedRecording bool `json:"enable-detailed-recording" mapstructure:"enable-detailed-recording"`
}

type UploadService struct {
	uploadStorage              uploadto.UploadStorage
	localCache                 *ristretto.Cache
	workerNums                 int
	recordChan                 chan *datastructure.AnalyticsRecord
	workerBufferSize           uint64
	recordsBufferFlushInterval time.Duration
	shouldStop                 uint32
	waitGroup                  sync.WaitGroup
}

var uploadService *UploadService

// CreateUploadService returns a new upload instance.
func CreateUploadService(opts *UploadConfig, storage uploadto.UploadStorage) *UploadService {
	wn := opts.WorkersNum
	recordBufferSize := opts.RecordsBufferSize
	workerBufferSize := recordBufferSize / uint64(wn)
	log.Debug().Uint64("workerBufferSize", workerBufferSize).Msgf("analytics pool worker buffer size")
	recordsChan := make(chan *datastructure.AnalyticsRecord, recordBufferSize)
	uploadService = &UploadService{
		uploadStorage:              storage,
		workerNums:                 wn,
		recordChan:                 recordsChan,
		workerBufferSize:           workerBufferSize,
		recordsBufferFlushInterval: opts.FlushInterval,
	}
	return uploadService
}

func (u *UploadService) Start() {
	atomic.SwapUint32(&u.shouldStop, 0)
	for i := 0; i < u.workerNums; i++ {
		u.waitGroup.Add(1)
		workerIs := i
		go u.runWorker(workerIs)
	}
}

func (u *UploadService) runWorker(workerId int) {
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

func (u *UploadService) UploadRecord(record *datastructure.AnalyticsRecord) error {
	if u == nil {
		return nil
	}
	// 检查信号
	if atomic.LoadUint32(&u.shouldStop) > 0 {
		return nil
	}

	u.recordChan <- record
	log.Info().EmbedObject(record).Msg("record uploaded successfully.")
	return nil
}

func (u *UploadService) GetStorage() uploadto.UploadStorage {
	return u.uploadStorage
}

func GetUploadService() *UploadService {
	return uploadService
}
