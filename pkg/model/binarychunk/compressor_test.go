package binarychunk

import (
	"compress/gzip"
	"context"
	"crypto/rand"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFileSystemGzipCompressor(t *testing.T) {
	t.Run("Should return a reader pointing decompressable binary", func(t *testing.T) {
		c := NewFileSystemGzipCompressor("/tmp")
		sourceData := make([]byte, 100000)
		rand.Read(sourceData)
		tmpDestFile, err := os.CreateTemp("/tmp", "khi-test-")
		if err != nil {
			t.Errorf("err was not a nil:%v", err)
		}
		tmpDestFile.Write(sourceData)
		reader, err := os.Open(tmpDestFile.Name())
		if err != nil {
			t.Errorf("err was not a nil:%v", err)
		}

		compressResult, err := c.CompressAll(context.Background(), reader)
		if err != nil {
			t.Errorf("err was not a nil:%v", err)
		}

		decompressedResult, err := gzip.NewReader(compressResult)
		if err != nil {
			t.Errorf("err was not a nil:%v", err)
		}
		result, err := io.ReadAll(decompressedResult)
		if err != nil {
			t.Errorf("err was not a nil:%v", err)
		}
		if diff := cmp.Diff(sourceData, result); diff != "" {
			t.Errorf("+sourceData, -result,%s", diff)
		}
	})
}
