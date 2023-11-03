package photon

/*
#include <stdlib.h>
#include <stdio.h>

struct photon_str_input {
  size_t curr_char, code_len;
  const char* code;
};

typedef struct photon_str_input* photon_str_input_t;

int photon_readc_str(void* data) {
  photon_str_input_t input = data;

  return input->curr_char >= input->code_len ? EOF : input->code[input->curr_char++];
}

int photon_readc_file(void* data) {
  FILE* fp = data;

  return fgetc(fp);
}
*/
import "C"
import (
	"path/filepath"
	"unsafe"
)

// ReadFunc is a reading function passed to c2mir_compile.
type ReadFunc *[0]byte // int (void*)

// Input is a C2MIR source input.
type Input interface {
	// SourceName returns the source file name, such as "test.c".
	SourceName() string
	// Data returns the data pointer passed to c2mir_compile.
	Data() unsafe.Pointer
	// ReadFunc returns the reading function passed to c2mir_compile.
	ReadFunc() ReadFunc
	// Close cleans up residual data after consuming the input, this is called automatically by C2MIRContext.Compile.
	Close()
}

// StringInput is an Input implementation that supplies a source code string.
type StringInput struct {
	// sourceName is the source file name passed to NewStringInput.
	sourceName string
	// data is the native data holder.
	data C.photon_str_input_t
}

// NewStringInput creates a new input that supplies a source code string.
// sourceName is the source file name passed to c2mir_compile.
// data is the source code string.
func NewStringInput(sourceName string, data string) Input {
	return &StringInput{
		sourceName: sourceName,
		data: &C.struct_photon_str_input{
			curr_char: C.size_t(0),
			code_len:  C.size_t(len(data)),
			code:      C.CString(data),
		},
	}
}

// SourceName returns the source file name, such as "test.c".
func (si *StringInput) SourceName() string {
	return si.sourceName
}

// Data returns the data pointer passed to c2mir_compile.
func (si *StringInput) Data() unsafe.Pointer {
	if si.data == nil {
		panic("tried to access closed input data")
	}

	return unsafe.Pointer(si.data)
}

// ReadFunc returns the reading function passed to c2mir_compile.
func (si *StringInput) ReadFunc() ReadFunc {
	return (*[0]byte)(C.photon_readc_str)
}

// Close cleans up residual data after consuming the input.
func (si *StringInput) Close() {
	C.free(unsafe.Pointer(si.data.code))
	si.data = nil
}

// FileInput is an Input implementation that supplies source code from a file.
type FileInput struct {
	// f is the file descriptor from which source code should be read.
	f *File
}

// NewFileInput creates a new input that supplies source code from a file.
// file is the file descriptor from which source code should be read.
func NewFileInput(f *File) Input {
	return &FileInput{f: f}
}

// SourceName returns the source file name, such as "test.c".
func (fi *FileInput) SourceName() string {
	return filepath.Base(fi.f.Path())
}

// Data returns the data pointer passed to c2mir_compile.
func (fi *FileInput) Data() unsafe.Pointer {
	return unsafe.Pointer(fi.f.fpOrNil())
}

// ReadFunc returns the reading function passed to c2mir_compile.
func (fi *FileInput) ReadFunc() ReadFunc {
	return (*[0]byte)(C.photon_readc_file)
}

// Close cleans up residual data after consuming the input.
func (fi *FileInput) Close() {
	fi.f.Close()
}
