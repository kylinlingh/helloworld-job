package memory

import (
	"errors"
	"github.com/dgraph-io/ristretto"
	"helloworld/internal/pump/uploadto"
	log "helloworld/pkg/logger"
	"sync"
)

/*
数据上报的目的地：本地内存里的 channel
*/

type MessageList struct {
	Count   int
	ValList [][]byte
	Mutext  sync.Mutex
}

type MemoryStorage struct {
	cache *ristretto.Cache
}

func (m *MemoryStorage) GetStorage() uploadto.UploadStorage {
	return m
}

var once sync.Once

func (m *MemoryStorage) Connect() bool {
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
			log.Fatal().Err(err).Msg("failed to call ristretto.NewCache()")
		}
	})
	return true
}

func (m *MemoryStorage) AppendToSetPipelined(key string, values [][]byte) {
	var ml *MessageList
	cval, ok := m.cache.Get(key)
	if cval == nil || !ok {
		ml = &MessageList{
			Count:   0,
			ValList: make([][]byte, 0),
		}
		ok = m.cache.Set(key, ml, 1)
		if !ok {
			log.Error().Err(errors.New("failed to set ristretto cache"))
		}
		// 必须等待写入成功
		m.cache.Wait()
		cval, _ = m.cache.Get(key)
	}
	ml = cval.(*MessageList)

	ml.Mutext.Lock()
	defer ml.Mutext.Unlock()
	for _, val := range values {
		ml.ValList = append(ml.ValList, val)
	}
	ml.Count += len(values)

	val, ok := m.cache.Get(key)
	if !ok {
		log.Error().Msg("cannot get items from cache")
	}
	mv := val.(*MessageList)
	log.Trace().Int("count", len(mv.ValList)).Msg("record has been uploaded to memory cache")

}
