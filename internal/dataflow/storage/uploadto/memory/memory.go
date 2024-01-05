package memory

import (
	"errors"
	"github.com/dgraph-io/ristretto"
	"helloworld/internal/dataflow/datastructure"
	log "helloworld/pkg/logger"
	"sync"
)

type UploadMemoryStorage struct {
	cache *ristretto.Cache
}

func New(c *ristretto.Cache) *UploadMemoryStorage {
	return &UploadMemoryStorage{
		cache: c,
	}
}

var once sync.Once

func (u *UploadMemoryStorage) Connect() bool {
	once.Do(func() {
		var err error
		c := &ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
			Cost:        nil,
		}
		u.cache, err = ristretto.NewCache(c)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to call ristretto.NewCache()")
		}
	})
	return true
}

var uonce sync.Once

func (u *UploadMemoryStorage) AppendToSetPipelined(key string, values [][]byte) {

	uonce.Do(func() {
		u.cache.Set(key, &datastructure.MessageList{}, 1)
		// 必须等待写入成功
		u.cache.Wait()
	})
	cval, ok := u.cache.Get(key)
	if cval == nil || !ok {
		log.Error().Err(errors.New("failed to get value from ristretto cache"))
	}
	var ml *datastructure.MessageList
	ml = cval.(*datastructure.MessageList)

	ml.Mutext.Lock()
	defer ml.Mutext.Unlock()
	for _, val := range values {
		ml.ValList = append(ml.ValList, val)
	}
	log.Info().Int("count", len(ml.ValList)).Msg("records has been uploaded to memory cache.")

	//val, ok := u.cache.Get(key)
	//if !ok {
	//	log.Error().Msg("cannot get items from cache")
	//}
	//mv := val.(*datastructure.MessageList)
	//log.Info().Int("count", len(mv.ValList)).Msg("records has been uploaded to memory cache.")
}

func (u *UploadMemoryStorage) AppendToSetPipelined_out(key string, values [][]byte) {
	var ml *datastructure.MessageList
	cval, ok := u.cache.Get(key)
	if cval == nil || !ok {

		ml = &datastructure.MessageList{
			ValList: make([][]byte, 0),
		}
		ok = u.cache.Set(key, ml, 1)
		if !ok {
			log.Error().Err(errors.New("failed to set ristretto cache"))
		}
		// 必须等待写入成功
		u.cache.Wait()
		cval, _ = u.cache.Get(key)
	}
	ml = cval.(*datastructure.MessageList)

	ml.Mutext.Lock()
	defer ml.Mutext.Unlock()
	for _, val := range values {
		ml.ValList = append(ml.ValList, val)
	}
	//ml.Count += len(values)

	val, ok := u.cache.Get(key)
	if !ok {
		log.Error().Msg("cannot get items from cache")
	}
	mv := val.(*datastructure.MessageList)
	log.Info().Int("count", len(mv.ValList)).Msg("records has been uploaded to memory cache.")
}

func (u UploadMemoryStorage) GetStorage() *ristretto.Cache {
	return u.cache
}
