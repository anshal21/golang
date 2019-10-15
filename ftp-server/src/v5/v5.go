package v5

import (
	"code.uber.internal/learning/golang-engucation/ftp-server/src/lib"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/operator"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/v3"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	_filePathSeparator = string(filepath.Separator)

	_channelBufferSize = 5

	_workerCount = 4
)

type v5 struct {
	operator.Operator
	fileCache map[string][]byte
	cacheHit int
	cacheMiss int
}

func New() operator.Operator {
	return v5{
		Operator: v3.New(),
	}
}


func (v v5) Delete(dirName string) (err error) {

	tasks := make(chan string, _channelBufferSize)

	var enqueueErr error

	go func() {
		defer close(tasks)
		enqueueErr = v.enqueueForDeletion(tasks, dirName)
		if enqueueErr != nil {
			err = enqueueErr
		}
	}()

	var wg sync.WaitGroup

	wg.Add(_workerCount)

	for i := 0; i < _workerCount; i++ {
		go func() {
			defer wg.Done()
			for task := range tasks {
				removeErr := lib.RemoveFile(task)
				if removeErr != nil {
					err = removeErr
				}
			}
		}()
	}

	wg.Wait()

	return err
}



func (v v5) enqueueForDeletion(tasks chan string, dirName string) error {
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
			err = v.Delete(filePath)
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

