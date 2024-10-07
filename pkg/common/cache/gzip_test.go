package cache

import (
	"crypto/rand"
	"testing"
)

func TestGZipCacheItemStorageProvider(t *testing.T) {
	parent := newTestOnMemoryCacheItemStorageProvider()
	gzip := NewGZipCacheItemStorageProvider(parent)
	original := []byte("hello world")
	err := gzip.Set("test", original)
	if err != nil {
		t.Errorf("Set failed: %s", err)
	}
	result, err := gzip.Get("test")
	if err != nil {
		t.Errorf("Get failed: %s", err)
	}
	if string(result) != string(original) {
		t.Errorf("Unexpected result: expected:%s,actual:%s", string(original), string(result))
	}
}

func BenchmarkGZipCacheItemStorageProvider__Set(b *testing.B) {
	parent := newTestOnMemoryCacheItemStorageProvider()
	gzip := NewGZipCacheItemStorageProvider(parent)
	comp := make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		rand.Read(comp)
		gzip.Set("test", comp)
	}
}

func BenchmarkGZipCacheItemStorageProvider__Get(b *testing.B) {
	parent := newTestOnMemoryCacheItemStorageProvider()
	gzip := NewGZipCacheItemStorageProvider(parent)
	comp := make([]byte, 1024)
	rand.Read(comp)
	gzip.Set("test", comp)
	for i := 0; i < b.N; i++ {
		gzip.Get("test")
	}
}
