package logrotate

import (
	"fmt"
	"os"
)

// If the folder does not exist, create the folder
func createDir(dir string) error {
	if dir == "" {
		return fmt.Errorf("dir is empty")
	}
	return os.MkdirAll(dir, 0755)
}
