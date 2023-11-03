package photon

/*
#include <mir.h>

size_t photon_proto_get_nargs(MIR_proto_t v) {
  return VARR_LENGTH(MIR_var_t, v->args);
}
*/
import "C"
import "unsafe"

type Proto struct {
	c C.MIR_proto_t
}

func (p *Proto) Res() int {
	return int(p.c.nres)
}

func (p *Proto) Args() int {
	return int(C.photon_proto_get_nargs(p.c))
}

func (p *Proto) ResTypes() []DataType {
	resLen := p.Res()

	cResTypes := unsafe.Slice(&p.c.res_types, resLen)
	resTypes := make([]DataType, resLen)
	for i, rt := range cResTypes {
		resTypes[i] = DataType(*rt)
	}

	return resTypes
}
