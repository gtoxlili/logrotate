package logrotate

import (
	"fmt"
	"os"
	"sync"
)

// If the folder does not exist, create the folder
func createDir(dir string) error {
	if dir == "" {
		return fmt.Errorf("dir is empty")
	}
	return os.MkdirAll(dir, 0755)
}

type race struct {
	mu sync.Mutex
}

type Option func(*Writer)

// WithRace 竞态
func WithRace() Option {
	return func(w *Writer) {
		w.race = &race{}
	}
}
