package filepath

import (
	"os"
	"path/filepath"
	"runtime"
)

func InitRootFolder(path string) error {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), path)
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}
