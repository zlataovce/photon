package photon

/*
#include <mir-gen.h>
*/
import "C"
import (
	"unsafe"
)

type OptLevel uint

const (
	OptLevelZero OptLevel = iota
	OptLevelOne
	OptLevelTwo
	OptLevelThree
)

// GenContext is a generation capable MIR context.
type GenContext struct {
	*Context

	gensNum int
}

func NewGenContext(gensNum int) *GenContext {
	c := NewContext()
	c.gen = true

	C.MIR_gen_init(c.c, C.int(gensNum))
	return &GenContext{Context: c, gensNum: gensNum}
}

func (gc *GenContext) NumGens() int {
	if gc.gensNum < 1 {
		return 1
	}

	return gc.gensNum
}

func (gc *GenContext) Generate(genNum int, funcItem *Item) unsafe.Pointer {
	gc.Context.checkClosed()

	return C.MIR_gen(gc.Context.c, C.int(genNum), funcItem.c)
}

func (gc *GenContext) DebugFile(genNum int, f *File) {
	gc.Context.checkClosed()

	C.MIR_gen_set_debug_file(gc.Context.c, C.int(genNum), f.fpOrNil())
}

func (gc *GenContext) DebugLevel(genNum int, level int) {
	gc.Context.checkClosed()

	C.MIR_gen_set_debug_level(gc.Context.c, C.int(genNum), C.int(level))
}

func (gc *GenContext) OptLevel(genNum int, level OptLevel) {
	gc.Context.checkClosed()

	C.MIR_gen_set_optimize_level(gc.Context.c, C.int(genNum), C.uint(level))
}
