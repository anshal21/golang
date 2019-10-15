package lib

import (
	"time"
	"os"
	)

func RemoveFile(filePath string) error {
	time.Sleep(time.Second*5)

	return os.Remove(filePath)
}
