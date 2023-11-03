package photon

/*
#include <mir.h>
*/
import "C"

type Data struct {
	c C.MIR_data_t
}

type RefData struct {
	c C.MIR_ref_data_t
}

type ExprData struct {
	c C.MIR_expr_data_t
}
