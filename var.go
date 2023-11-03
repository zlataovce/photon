package photon

/*
#include <mir.h>
*/
import "C"

type Var struct {
	type_ DataType
	name  string
	size  int
}

func NewVar(type_ DataType, name string) *Var {
	return NewVarWithSize(type_, name, 0)
}

func NewVarWithSize(type_ DataType, name string, size int) *Var {
	return &Var{type_: type_, name: name, size: size}
}

func (v *Var) Type() DataType {
	return v.type_
}

func (v *Var) Name() string {
	return v.name
}

func (v *Var) Size() int {
	return v.size
}
