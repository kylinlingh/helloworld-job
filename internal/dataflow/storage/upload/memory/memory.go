package memory

import "github.com/dgraph-io/ristretto"

type UploadMemoryStorage struct {
	cache *ristretto.Cache
}

func New(c *ristretto.Cache) *UploadMemoryStorage {
	return &UploadMemoryStorage{
		cache: c,
	}
}
