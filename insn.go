package photon

/*
#include <mir.h>
*/
import "C"

type Insn struct {
	c C.MIR_insn_t
}
