package cache

import (
	"bytes"
	"compress/gzip"
)

type GZipCacheItemStorageProvider struct {
	parent CacheItemStorageProvider
}

func NewGZipCacheItemStorageProvider(parent CacheItemStorageProvider) *GZipCacheItemStorageProvider {
	return &GZipCacheItemStorageProvider{parent: parent}
}

// Get implements CacheItemStorageProvider.
func (g *GZipCacheItemStorageProvider) Get(key string) ([]byte, error) {
	compressedData, err := g.parent.Get(key)
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(gr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Set implements CacheItemStorageProvider.
func (g *GZipCacheItemStorageProvider) Set(key string, value []byte) error {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(value)
	if err != nil {
		return err
	}
	err = gw.Close()
	if err != nil {
		return err
	}
	return g.parent.Set(key, buf.Bytes())
}

var _ CacheItemStorageProvider = (*GZipCacheItemStorageProvider)(nil)
