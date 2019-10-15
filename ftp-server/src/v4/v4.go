package v4

import (
	"code.uber.internal/learning/golang-engucation/ftp-server/src/lib"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/operator"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/v3"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	_filePathSeparator = string(filepath.Separator)

	_channelBufferSize = 5
)

type v4 struct {
	operator.Operator
	fileCache map[string][]byte
	cacheHit int
	cacheMiss int
}

func New() operator.Operator {
	return v4{
		Operator: v3.New(),
	}
}


func (v v4) Delete(dirName string) (err error) {

	tasks := make(chan string, _channelBufferSize)

	go v.enqueueForDeletion(tasks, dirName)

	for task := range tasks {

		go lib.RemoveFile(task)
	}

	return nil
}



func (v v4) enqueueForDeletion(tasks chan string, dirName string) error {
	dir, err := os.Open(dirName)
	if err != nil {
		return errors.New("File path doesn't exist")
	}

	files, _ := dir.Readdir(0)
	if err != nil {
		return errors.New("Internal Server Error")
	}

	for _, file := range files {
		filePath := strings.Join([]string{dirName, file.Name()}, _filePathSeparator)
		if file.IsDir() {
			err = v.enqueueForDeletion(tasks, filePath)
			if err != nil {
				return err
			}
		}
		tasks <- filePath
	}

	/*
		Potential problem here is race conditions between directory and files
		it is possible that there is an attempt to delete directory before deleting file
		TODO: Either fix this or have a single level dir
	 */
	tasks <- dirName

	return nil

}

