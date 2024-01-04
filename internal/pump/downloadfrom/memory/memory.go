package memory

import (
	"github.com/dgraph-io/ristretto"
	"helloworld/internal/pump/uploadto/memory"
	log "helloworld/pkg/logger"
	"sync"
)

/*
从哪里下载数据：本地缓存的 channel
*/

type MemoryStorage struct {
	cache *ristretto.Cache
}

var once sync.Once
var m MemoryStorage

func (m *MemoryStorage) Init(config interface{}) error {
	//if storage != nil {
	//	//m.cache = storage.GetStorage()
	//	return nil
	//}
	once.Do(func() {
		var err error
		c := &ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
			Cost:        nil,
		}
		m.cache, err = ristretto.NewCache(c)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize memory storage")
		}
	})
	return nil
}

func (m *MemoryStorage) GetName() string {
	return "memory"
}
func (m *MemoryStorage) Connect() bool {
	return true
}
func (m *MemoryStorage) GetAndDeleteSet(keyName string) []interface{} {
	log.Trace().Msgf("Getting raw key set: %s", keyName)
	val, ok := m.cache.Get(keyName)
	if !ok {
		log.Trace().Msg("no data in memory storage")
		return nil
	}
	ml := val.(*memory.MessageList)
	ml.Mutext.Lock()
	defer ml.Mutext.Unlock()

	result := make([]interface{}, ml.Count)
	for i, v := range ml.ValList {
		result[i] = v
	}
	return result
}
