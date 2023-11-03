package photon

/*
#include <mir.h>
*/
import "C"
import "unsafe"

type FuncLike interface {
	Res() int
	ResTypes() []DataType
	Args() int
}

type Func struct {
	c C.MIR_func_t
}

func (f *Func) Res() int {
	return int(f.c.nres)
}

func (f *Func) Args() int {
	return int(f.c.nargs)
}

func (f *Func) ResTypes() []DataType {
	resLen := f.Res()

	cResTypes := unsafe.Slice(&f.c.res_types, resLen)
	resTypes := make([]DataType, resLen)
	for i, rt := range cResTypes {
		resTypes[i] = DataType(*rt)
	}

	return resTypes
}
