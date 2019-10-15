package main

import (
	"code.uber.internal/learning/golang-engucation/ftp-server/src/operator"
	"code.uber.internal/learning/golang-engucation/ftp-server/src/v2"
	"fmt"
	"net/http"
	"time"
)

const (
	_operationURLParameterKey = "operation"
	_dirURLParamterKey        = "dir"

	_FTPOperationList   = "ls"
	_FTPOperationDelete = "del"
	_FTPOperationGet = "get"

	_serverPort = ":8000"
)

var (
	helpMessage = fmt.Sprintf(`Please provide one the following option
  - %v: list down the files in a directory
  - %v: delete files in a directory
  - %v: returns a file	
`, _FTPOperationList, _FTPOperationDelete, _FTPOperationGet)
)


type handler struct {
	operator operator.Operator
}

func NewHandler() handler {
	return handler{
		operator: v2.New(),
	}
}

func (h handler) ftpHandler(w http.ResponseWriter, req *http.Request) {

	operation, ok := req.URL.Query()[_operationURLParameterKey]
	if !ok {
		w.Write([]byte(helpMessage))
		return
	}


	dir, ok := req.URL.Query()[_dirURLParamterKey]
	if !ok {
		w.Write([]byte("Missing required parameter"))
		return
	}

	//fmt.Println(operation, dir)

	response := make([]byte, 1024*1024)

	switch operation[0] {
	case _FTPOperationList:
		reader, err := h.operator.List(dir[0])
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		n, _ := reader.Read(response)
		w.Write(response[:n])

	case _FTPOperationDelete:
		err := h.operator.Delete(dir[0])
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("RemovedDirectory"))

	case _FTPOperationGet:
		resp, err  := h.operator.Get(dir[0])
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		n, _ := resp.Read(response)
		w.Write(response[:n])
	}
}

// Select

func main() {
	handler := NewHandler()

	http.HandleFunc("/ftp", timmed(handler.ftpHandler))
	http.ListenAndServe(_serverPort, nil)
}



func timmed(f func(w http.ResponseWriter, req *http.Request)) func (w http.ResponseWriter, req *http.Request) {
	return func (w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		f(w, req)
		end := time.Now()

		fmt.Println("Request Time: %v", end.Sub(start))
	}
}
