package memory

import (
	"github.com/dgraph-io/ristretto"
	"helloworld/internal/dataflow/datastructure"
	log "helloworld/pkg/logger"
)

type DownloadMemoryStorage struct {
	cache *ristretto.Cache
}

func NewDownloadMemStorage(c *ristretto.Cache) *DownloadMemoryStorage {
	return &DownloadMemoryStorage{cache: c}
}

func (d *DownloadMemoryStorage) GetName() string {
	return "memory"
}

func (d *DownloadMemoryStorage) Connect() bool {
	return true
}

func (d *DownloadMemoryStorage) GetAndDeleteSet(keyName string) []interface{} {
	log.Trace().Msgf("Getting raw key set: %s", keyName)
	val, ok := d.cache.Get(keyName)
	if !ok {
		log.Trace().Msg("no data in memory storage")
		return nil
	}
	ml := val.(*datastructure.MessageList)
	ml.Mutext.Lock()
	defer ml.Mutext.Unlock()

	result := make([]interface{}, len(ml.ValList))
	for i, v := range ml.ValList {
		result[i] = v
	}

	// 清空缓存
	ml.ValList = ml.ValList[:0]
	return result
}
