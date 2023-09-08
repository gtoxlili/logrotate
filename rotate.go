package logrotate

import (
	"fmt"
	"github.com/gtoxlili/logrotate/mmap"
	"path/filepath"
	"strings"
	"time"
)

type Unit int8

const (
	Minutely Unit = iota
	Hourly
	Daily
	Never
	Illegal
	// Secondly is used for testing only.
	Secondly
)

func ParseKind(unit string) (Unit, error) {
	switch strings.ToLower(unit) {
	case "minutely":
		return Minutely, nil
	case "hourly":
		return Hourly, nil
	case "daily":
		return Daily, nil
	case "never":
		return Never, nil
	case "secondly":
		return Illegal, fmt.Errorf("secondly is only for test")
	default:
		return Illegal, fmt.Errorf("invalid rotation type, must be minutely, hourly, daily, or never")
	}
}

type Writer struct {
	lastRotateAt time.Time
	cursor       *mmap.Writer
	filePath     string
	// Rotation threshold.
	threshold time.Duration
	format    string
}

// NewWriter creates a new log Writer.
func NewWriter(unit Unit, dir, prefix string) (*Writer, error) {
	if err := createDir(dir); err != nil {
		return nil, err
	}
	// Get rotation threshold and format based on unit.
	var threshold time.Duration
	var format string
	switch unit {
	case Minutely:
		threshold = time.Minute
		format = "2006-01-02-15-04"
	case Hourly:
		threshold = time.Hour
		format = "2006-01-02-15"
	case Daily:
		threshold = 24 * time.Hour
		format = "2006-01-02"
	case Secondly:
		threshold = time.Second
		format = "2006-01-02-15-04-05"
	case Never:
		// Set a large threshold to ensure no rotation occurs.
		threshold = 100 * 365 * 24 * time.Hour
	}

	return &Writer{
		filePath:  filepath.Join(dir, prefix),
		threshold: threshold,
		format:    format,
	}, nil
}

func (w *Writer) shouldRotate() bool {
	return time.Since(w.lastRotateAt) >= w.threshold
}

func (w *Writer) initRotateTime() {
	tn := time.Now()
	tru := tn.Truncate(w.threshold)
	// Ensure using a monotonic clock for calculating time intervals.
	w.lastRotateAt = tn.Add(tru.Sub(tn))
}

func (w *Writer) rotate() error {
	// Close the old file.
	if w.cursor != nil {
		if err := w.cursor.Close(); err != nil {
			return err
		}
	}

	// Initialize the rotation time.
	w.initRotateTime()

	// Create a new file.
	var err error
	fp := w.filePath
	if w.threshold != 0 {
		fp += "." + w.lastRotateAt.Format(w.format)
	}
	w.cursor, err = mmap.New(fp)
	if err != nil {
		return err
	}

	return nil
}

// Write writes the log message to the file.
func (w *Writer) Write(p []byte) (n int, err error) {
	if w.cursor == nil || w.shouldRotate() {
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	return w.cursor.Write(p)
}

// Close closes the log writer.
func (w *Writer) Close() error {
	if w.cursor != nil {
		if err := w.cursor.Close(); err != nil {
			return err
		}
	}
	return nil
}
