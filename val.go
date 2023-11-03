package photon

/*
#include <mir.h>
#include <photon-math.h>

#define VAL_ACCESS(_name, _type, _item) \
  MIR_val_t photon_val_new_##_name(_type _item) { \
    MIR_val_t v; \
    v._item = _item; \
    return v; \
  } \
  _type photon_val_get_##_name(MIR_val_t v) { \
    return v._item; \
  }

VAL_ACCESS(insn_code, MIR_insn_code_t, ic)
VAL_ACCESS(ptr, void*, a)
VAL_ACCESS(int, int64_t, i)
VAL_ACCESS(uint, uint64_t, u)
VAL_ACCESS(float, float, f)
VAL_ACCESS(double, double, d)

MIR_val_t photon_val_new_ldouble(photon_ldouble_t ld) {
  MIR_val_t v;
  v.ld = photon_ldouble_unwrap(ld);

  return v;
}

photon_ldouble_t photon_val_get_ldouble(MIR_val_t v) {
  return photon_ldouble_wrap(v.ld);
}
*/
import "C"
import (
	"fmt"
	"math/big"
	"reflect"
	"unsafe"
)

type Val struct {
	c C.MIR_val_t
}

func NewVal(v interface{}) *Val {
	var cVal C.MIR_val_t
	switch v0 := v.(type) {
	case InsnCode:
		cVal = C.photon_val_new_insn_code(C.MIR_insn_code_t(v0))
	case unsafe.Pointer:
		cVal = C.photon_val_new_ptr(v0)
	case int:
		cVal = C.photon_val_new_int(C.int64_t(v0))
	case int8:
		cVal = C.photon_val_new_int(C.int64_t(v0))
	case int32:
		cVal = C.photon_val_new_int(C.int64_t(v0))
	case int64:
		cVal = C.photon_val_new_int(C.int64_t(v0))
	case uint:
		cVal = C.photon_val_new_uint(C.uint64_t(v0))
	case uint8:
		cVal = C.photon_val_new_uint(C.uint64_t(v0))
	case uint32:
		cVal = C.photon_val_new_uint(C.uint64_t(v0))
	case uint64:
		cVal = C.photon_val_new_uint(C.uint64_t(v0))
	case float32:
		cVal = C.photon_val_new_float(C.float(v0))
	case float64:
		cVal = C.photon_val_new_double(C.double(v0))
	case *big.Float:
		cVal = C.photon_val_new_ldouble(newLongDouble(v0))
	default:
		panic(fmt.Sprintf("unsupported val type %s", reflect.TypeOf(v).String()))
	}

	return &Val{c: cVal}
}

func (v *Val) InsnCode() InsnCode {
	return InsnCode(C.photon_val_get_insn_code(v.c))
}

func (v *Val) Pointer() unsafe.Pointer {
	return C.photon_val_get_ptr(v.c)
}

func (v *Val) Int() int64 {
	return int64(C.photon_val_get_int(v.c))
}

func (v *Val) Uint() uint64 {
	return uint64(C.photon_val_get_uint(v.c))
}

func (v *Val) Float() float32 {
	return float32(C.photon_val_get_float(v.c))
}

func (v *Val) Double() float64 {
	return float64(C.photon_val_get_double(v.c))
}

func (v *Val) LongDouble() *big.Float {
	return newBigFloat(C.photon_val_get_ldouble(v.c))
}
