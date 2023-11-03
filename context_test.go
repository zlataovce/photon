package photon

import (
	"testing"
)

func TestContext_Close(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	mirCtx := NewContext()
	defer mirCtx.Close()

	moduleCtx := mirCtx.NewModuleContext("m")
	moduleCtx.Close()

	moduleCtx.NewFunc("loop", DataTypeI64, NewVar(DataTypeI64, "arg1")) // not allowed, should panic
}

func TestContext_ModuleList(t *testing.T) {
	ctx := NewContext()
	defer ctx.Close()

	moduleCtx := ctx.NewModuleContext("test_module")
	moduleCtx.Close()

	modules := ctx.ModuleList()
	if len(modules) != 1 {
		t.Fatalf("expected 1 module, got %d modules", len(modules))
	}

	for i, module := range modules {
		t.Logf("module %s, index %d", module.Name(), i)
	}
}

func TestContext_NewModuleContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	mirCtx := NewContext()
	defer mirCtx.Close()

	mirCtx.NewModuleContext("m")
	mirCtx.NewModuleContext("m2") // already building a module, panic
}

func TestModuleContext_NewFunc(t *testing.T) {
	mirCtx := NewContext()
	defer mirCtx.Close()

	moduleCtx := mirCtx.NewModuleContext("m")

	funcCtx := moduleCtx.NewFunc("loop", DataTypeI64, NewVar(DataTypeI64, "arg1"))

	r2 := funcCtx.NewFuncReg(funcCtx.Func(), DataTypeI64, "count")
	arg1 := funcCtx.Reg("arg1", funcCtx.Func())
	fin := funcCtx.NewLabel()
	cont := funcCtx.NewLabel()

	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeMOV, funcCtx.NewRegOp(r2), funcCtx.NewIntOp(0)))
	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeBGE, funcCtx.NewLabelOp(fin), funcCtx.NewRegOp(r2), funcCtx.NewRegOp(arg1)))
	funcCtx.AppendInsn(cont)
	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeADD, funcCtx.NewRegOp(r2), funcCtx.NewRegOp(r2), funcCtx.NewIntOp(1)))
	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeBLT, funcCtx.NewLabelOp(cont), funcCtx.NewRegOp(r2), funcCtx.NewRegOp(arg1)))
	funcCtx.AppendInsn(fin)
	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeRET, funcCtx.NewRegOp(r2)))

	funcCtx.Close()
	moduleCtx.Close()

	outFile := OpenFile("new_func_test.mir", FileIOModeWrite)
	defer outFile.Close()

	mirCtx.WriteTextual(outFile)
}

func TestModuleContext_Link(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	mirCtx := NewContext()
	defer mirCtx.Close()

	mirCtx.Link(LinkFuncGenInterface, nil) // not generation capable, panic
}

func TestModuleContext_Interpret(t *testing.T) {
	mirCtx := NewContext()
	defer mirCtx.Close()

	moduleCtx := mirCtx.NewModuleContext("m")

	funcCtx := moduleCtx.NewFunc("loop", DataTypeI64, NewVar(DataTypeI64, "arg1"))

	count := funcCtx.NewFuncReg(funcCtx.Func(), DataTypeI64, "count")
	arg1 := funcCtx.Reg("arg1", funcCtx.Func())
	fin := funcCtx.NewLabel()
	cont := funcCtx.NewLabel()

	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeMOV, funcCtx.NewRegOp(count), funcCtx.NewIntOp(0)))
	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeBGE, funcCtx.NewLabelOp(fin), funcCtx.NewRegOp(count), funcCtx.NewRegOp(arg1)))
	funcCtx.AppendInsn(cont)
	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeADD, funcCtx.NewRegOp(count), funcCtx.NewRegOp(count), funcCtx.NewIntOp(1)))
	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeBLT, funcCtx.NewLabelOp(cont), funcCtx.NewRegOp(count), funcCtx.NewRegOp(arg1)))
	funcCtx.AppendInsn(fin)
	funcCtx.AppendInsn(funcCtx.NewInsn(InsnCodeRET, funcCtx.NewRegOp(count)))

	funcCtx.Close()
	moduleCtx.Close()

	t.Logf("Before simplification:")
	mirCtx.WriteTextual(StderrFile) // need to write to stderr, because go test eats stdout

	mirCtx.LoadModule(moduleCtx.Module)
	mirCtx.Link(LinkFuncInterpInterface, nil)

	t.Logf("After simplification:")
	mirCtx.WriteTextual(StderrFile) // need to write to stderr, because go test eats stdout

	v := NewVal(10000000)
	t.Logf("%d", v.Int())

	results := mirCtx.Interpret(funcCtx.Item, v)

	t.Logf("%d", results[0].Int())
}
