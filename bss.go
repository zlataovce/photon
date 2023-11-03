package photon

/*
#include <mir.h>
*/
import "C"

type BSS struct {
	c C.MIR_bss_t
}
