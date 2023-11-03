package photon

import "testing"

func TestGenContext_Generate(t *testing.T) {
	mirCtx := NewGenContext(1)
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
	mirCtx.Link(LinkFuncGenInterface, nil)

	t.Logf("After simplification:")
	mirCtx.WriteTextual(StderrFile) // need to write to stderr, because go test eats stdout

	ptr := mirCtx.Generate(0, funcCtx.Item)
	t.Logf("%p", ptr)
}
