package v3

import (
	"bytes"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/operator"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/v2"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	_filePathSeparator = string(filepath.Separator)

	_fileDelimiter = "\n"
)

type v3 struct {

	operator.Operator

	fileCache map[string][]byte
	cacheHit int
	cacheMiss int
}

func New() operator.Operator {
	return v3{
		Operator: v2.New(),
		fileCache: make(map[string][]byte),
	}
}

func (v v3) Get(fileName string) (io.Reader, error) {
	if file, ok := v.fileCache[fileName]; ok {
		v.cacheHit++
		fmt.Printf("Info: Cache Hit For File: %v,Total  Cache Hit Count: %v\n", fileName, v.cacheHit)
		return bytes.NewReader(file), nil
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	v.cacheMiss++
	buff := make([]byte, 1024*1024)
	n, err := f.Read(buff)

	fmt.Printf("Info: Cache Miss For File: %v, Total Cache Miss Count: %v\n", fileName, v.cacheMiss)

	v.fileCache[fileName] = buff[:n]

	return bytes.NewReader(buff[:n]), nil


}