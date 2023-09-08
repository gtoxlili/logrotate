package mmap

import (
	"os"
	"syscall"
)

var (
	// BatchMemSize is the memory size mapped each time.
	BatchMemSize = int64(os.Getpagesize())
)

type Writer struct {
	file       *os.File
	data       []byte
	pageBatch  int64
	pageOffset int64
}

// New creates a new mmap Writer.
func New(filePath string) (*Writer, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	info, _ := file.Stat()
	mw := &Writer{
		file:       file,
		pageBatch:  info.Size() / BatchMemSize,
		pageOffset: info.Size() % BatchMemSize,
	}
	if err := mw.mmap(); err != nil {
		return nil, err
	}
	return mw, nil
}

func (mw *Writer) mmap() error {
	if mw.data != nil {
		if err := syscall.Munmap(mw.data); err != nil {
			return err
		}
	}
	// Increase the file size.
	if err := mw.file.Truncate(BatchMemSize * (mw.pageBatch + 1)); err != nil {
		return err
	}
	b, err := syscall.Mmap(int(mw.file.Fd()), mw.pageBatch*BatchMemSize, int(BatchMemSize), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
	mw.data = b
	return nil
}

// Close closes the mmap Writer.
func (mw *Writer) Close() error {
	if err := syscall.Munmap(mw.data); err != nil {
		return err
	}
	if err := mw.file.Truncate(mw.pageBatch*BatchMemSize + mw.pageOffset); err != nil {
		return err
	}
	return mw.file.Close()
}

// Write writes data to the mmap.
func (mw *Writer) Write(p []byte) (int, error) {
	if len(p) > int(BatchMemSize) {
		return 0, os.ErrInvalid
	}
	var n int
	// Check if expansion is needed.
	if mw.pageOffset+int64(len(p)) > BatchMemSize {
		n += copy(mw.data[mw.pageOffset:], p[:BatchMemSize-mw.pageOffset])
		p = p[BatchMemSize-mw.pageOffset:]
		mw.pageBatch++
		if err := mw.mmap(); err != nil {
			return 0, err
		}
		mw.pageOffset = 0
	}
	nn := copy(mw.data[mw.pageOffset:], p)
	mw.pageOffset += int64(nn)
	return n + nn, nil
}
