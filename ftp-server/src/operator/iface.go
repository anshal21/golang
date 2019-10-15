package operator

import (
	"io"
)

// Operator is an interface with allowed FTP methods
type Operator interface {
	// Delete recursively deletes a given directory
	Delete(dirName string) error

	// List lists all the files in a given directory
	List(dirName string) (io.Reader, error)

	// Get returns the requested file
	Get(fileName string) (io.Reader, error)
}