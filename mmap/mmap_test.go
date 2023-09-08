package mmap

import (
	"testing"
)

func TestNewMmap(t *testing.T) {
	mp, err := New("tmp.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer mp.Close()
	for i := 0; i < int(BatchMemSize); i++ {
		mp.Write([]byte("abcd\n"))
	}
}
