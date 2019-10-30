package sys

import (
	"os"
	"path/filepath"
)

func GetProcessDir()string{
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}
