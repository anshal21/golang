package v1

import (
	"bytes"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/operator"
	"io"
	"os"
	"strings"
	"path/filepath"

)

type v1 struct {}

const (
	_fileDelimiter = "\n"

	_filePathSeparator = string(filepath.Separator)
)


func New() operator.Operator {
	return v1{}
}

func (v v1) List(dirName string) (io.Reader, error) {
	dir, _ := os.Open(dirName)
	files, _ := dir.Readdir(0)

	fileNames := make([]string, 0)

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	output := strings.Join(fileNames, _fileDelimiter)
	return bytes.NewReader([]byte(output)), nil
}



func (v v1) Delete(dirName string) error {

	dir, _ := os.Open(dirName)
	files, _ := dir.Readdir(0)

	for _, file := range files {
		filePath := strings.Join([]string{dirName, file.Name()}, _filePathSeparator)
		if file.IsDir() {
			v.Delete(filePath)
		}else{
			os.Remove(filePath)
		}

	}

	os.Remove(dirName)

	return nil
}

func (v v1) Get(fileName string) (io.Reader, error) {
	return nil, nil
}

