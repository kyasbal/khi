package binarychunk

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type BinaryReference struct {
	Offset int `json:"offset"`
	Length int `json:"len"`
	Buffer int `json:"buffer"`
}

// LargeBinaryWriter stores text as a large binary chunk and returns BinaryReference points the buffer location.
type LargeBinaryWriter interface {
	// Check if the specified text can fit in the buffer
	CanWrite(size int) bool
	// Write the specified text and returns the BinaryReference
	Write(data []byte) (*BinaryReference, error)
	// Read buffer from a BinaryReference
	Read(ref *BinaryReference) ([]byte, error)
	// Obtain the result binary as io.Reader
	GetBinary() (io.Reader, error)
	// Free allocated resource for the writer
	Dispose() error
}

// FileSystemBinaryWriter is a basic implementation of the LargeTextWriter.
type FileSystemBinaryWriter struct {
	bufferIndex       int
	maximumBufferSize int
	currentLength     int
	disposed          bool
	file              *os.File
	fileMutex         sync.Mutex
}

var _ LargeBinaryWriter = (*FileSystemBinaryWriter)(nil)

func NewFileSystemBinaryWriter(tmpPath string, bufferIndex int, maxSize int) (*FileSystemBinaryWriter, error) {
	file, err := os.CreateTemp(tmpPath, "khi-")
	if err != nil {
		return nil, err
	}
	return &FileSystemBinaryWriter{
		bufferIndex:       bufferIndex,
		maximumBufferSize: maxSize,
		currentLength:     0,
		disposed:          false,
		file:              file,
		fileMutex:         sync.Mutex{},
	}, nil
}

func (w *FileSystemBinaryWriter) CanWrite(size int) bool {
	return !w.disposed && w.currentLength+size <= w.maximumBufferSize
}

func (w *FileSystemBinaryWriter) Write(data []byte) (*BinaryReference, error) {
	w.fileMutex.Lock()
	defer w.fileMutex.Unlock()
	if !w.CanWrite(len(data)) {
		return nil, fmt.Errorf("buffer can't write the specified length %d (current:%d,maximum:%d)", len(data), w.currentLength, w.maximumBufferSize)
	}
	_, err := w.file.Seek(int64(w.currentLength), io.SeekStart)
	if err != nil {
		return nil, err
	}
	size, err := w.file.Write(data)
	if err != nil {
		return nil, err
	}
	reference := &BinaryReference{
		Buffer: w.bufferIndex,
		Length: size,
		Offset: w.currentLength,
	}
	w.currentLength += size
	return reference, nil
}

func (w *FileSystemBinaryWriter) Read(ref *BinaryReference) ([]byte, error) {
	if ref.Buffer != w.bufferIndex {
		return nil, fmt.Errorf("invalid buffer index. it's not current buffer index")
	}
	w.fileMutex.Lock()
	defer w.fileMutex.Unlock()
	_, err := w.file.Seek(int64(ref.Offset), io.SeekStart)
	if err != nil {
		return nil, err
	}
	result := make([]byte, ref.Length)
	_, err = w.file.Read(result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (w *FileSystemBinaryWriter) GetBinary() (io.Reader, error) {
	if w.disposed {
		return nil, fmt.Errorf("instance is already disposed.")
	}
	file, err := os.Open(w.file.Name())
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (w *FileSystemBinaryWriter) Dispose() error {
	if w.disposed {
		return fmt.Errorf("instance is already disposed.")
	}
	return w.file.Close()
}
