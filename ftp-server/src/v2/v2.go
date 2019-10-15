package v2

import (
	"bytes"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/lib"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/operator"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

const (
	_filePathSeparator = string(filepath.Separator)

	_fileDelimiter = "\n"
)

type v2 struct {}

func New() operator.Operator {
	return v2{}
}

func (v v2) Delete(dirName string) error {

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
		} else{
			err := lib.RemoveFile(filePath)
			if err != nil {
				return err
			}
			log("DELETE", filePath)
		}

	}

	log("DELETE", dirName, "type", "directory")
	err = lib.RemoveFile(dirName)
	if err != nil {
		return err
	}


	return nil
}



func (v v2) List(dirName string) (io.Reader, error) {
	dir, err := os.Open(dirName)
	if err != nil {
		return nil, errors.New("File path doesn't exist")
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}

	fileNames := make([]string, 0)

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	output := strings.Join(fileNames, _fileDelimiter)
	return bytes.NewReader([]byte(output)), nil
}

func (v v2) Get(fileName string) (io.Reader, error) {
	return nil, nil
}

func log(operation string, fileName string, metaDataKeyValue ...string) {
	if len(metaDataKeyValue)%2 != 0  || len(metaDataKeyValue) == 0 {
		fmt.Printf("%v: %v\n", operation, fileName)
		return
	}

	logContext := make(map[string]string)
	for i := 0; i < len(metaDataKeyValue); i += 2 {
		logContext[metaDataKeyValue[i]] = metaDataKeyValue[i+1]
	}
	jsonContext, _ := json.Marshal(&logContext)
	fmt.Printf("%v: %v, meta: %v\n", operation, fileName, string(jsonContext))
}