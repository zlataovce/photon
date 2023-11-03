package photon

import "testing"

func Test_photon_handle_error(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	mirCtx := NewContext()
	moduleCtx := mirCtx.NewModuleContext("m")
	funcCtx := moduleCtx.NewFunc("loop", DataTypeI64)

	funcCtx.NewFuncReg(funcCtx.Func(), DataTypeI64, "count")
	funcCtx.Reg("arg1", funcCtx.Func()) // undeclared_func_reg error

	funcCtx.Close()
	moduleCtx.Close()
	mirCtx.Close() // deferring this call makes it run even after panicking = segmentation fault, because of a dirty state
}
