package photon

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"path/filepath"
	"unsafe"
)

// FileIOMode is a file manipulation mode string, passed to fopen.
type FileIOMode string

const (
	// FileIOModeRead opens a file in read mode.
	FileIOModeRead FileIOMode = "r"
	// FileIOModeWrite opens or creates a file in write mode.
	FileIOModeWrite FileIOMode = "w"
	// FileIOModeAppend opens a file in append mode.
	FileIOModeAppend FileIOMode = "a"
	// FileIOModeReadPlus opens a file in both read and write mode.
	FileIOModeReadPlus FileIOMode = "r+"
	// FileIOModeWritePlus opens a file in both read and write mode.
	FileIOModeWritePlus FileIOMode = "w+"
	// FileIOModeAppendPlus opens a file in both read and write mode.
	FileIOModeAppendPlus FileIOMode = "a+"
)

var (
	// StdoutFile is a wrapped standard output file descriptor.
	StdoutFile = &File{fp: C.stdout}
	// StdinFile is a wrapped standard input file descriptor.
	StdinFile = &File{fp: C.stdin}
	// StderrFile is a wrapped standard error file descriptor.
	StderrFile = &File{fp: C.stderr}
)

// File is a C file descriptor wrapper.
type File struct {
	// path is the file path supplied to OpenFile, empty for standard file descriptor wrappers.
	path string
	// fp is the wrapped file descriptor.
	fp *C.FILE
}

// OpenFile opens and wraps a new file descriptor using fopen.
// path is the file path, cleaned with filepath.Clean.
// mode is the file manipulation mode.
// Returns the wrapped file or nil, if the path was empty.
func OpenFile(path string, mode FileIOMode) *File {
	if path == "" {
		return nil
	}

	cleanedPath := filepath.Clean(path)
	cPath := C.CString(cleanedPath)
	defer C.free(unsafe.Pointer(cPath))

	cMode := C.CString(string(mode))
	defer C.free(unsafe.Pointer(cMode))

	return &File{path: cleanedPath, fp: C.fopen(cPath, cMode)}
}

// Path returns the file path supplied to OpenFile, empty for standard file descriptor wrappers.
func (f *File) Path() string {
	return f.path
}

// Close closes the underlying file descriptor.
func (f *File) Close() {
	C.fclose(f.fp)
	f.fp = nil
}

// fpOrNil transforms this wrapper to a *FILE, accounting for nil values.
func (f *File) fpOrNil() *C.FILE {
	if f == nil {
		return nil
	}
	if f.fp == nil {
		panic("tried to access closed file descriptor")
	}

	return f.fp
}
