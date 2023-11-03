package photon

/*
#include <mir.h>
*/
import "C"
import (
	"fmt"
	"strconv"
)

type ErrorFunc func(type_ ErrorType, msg string)

var DefaultErrorFunc ErrorFunc = func(type_ ErrorType, msg string) {
	name := strconv.Itoa(int(type_))
	if errName := type_.Name(); errName != "" {
		name = errName
	}

	panic(fmt.Sprintf("error %s: %s", name, msg))
}

var CurrentErrorHandler = DefaultErrorFunc

//export photon_handle_error
func photon_handle_error(type_ C.MIR_error_type_t, msg *C.char) {
	CurrentErrorHandler(ErrorType(type_), C.GoString(msg))
}
