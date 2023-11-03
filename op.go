package photon

/*
#include <mir.h>
*/
import "C"

type Op struct {
	c C.MIR_op_t
}
