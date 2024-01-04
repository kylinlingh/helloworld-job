package memory

import "github.com/dgraph-io/ristretto"

type DownloadMemoryStroage struct {
	cache *ristretto.Cache
}

func New(c *ristretto.Cache) *DownloadMemoryStroage {
	return &DownloadMemoryStroage{cache: c}
}
