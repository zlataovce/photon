package photon

import (
	"testing"
)

func TestC2MIRContext_Compile(t *testing.T) {
	mirCtx := NewContext()
	defer mirCtx.Close()

	c2mirCtx := NewC2MIRContext(mirCtx)
	defer c2mirCtx.Close()

	input := NewStringInput("test.c", "#include <stdint.h>\nint64_t loop(int64_t arg1) { int64_t count = 0; while (count < arg1) count++; return count; }\n")
	err := c2mirCtx.Compile(&Options{
		MessageFile:      StdoutFile,
		Debug:            true,
		Verbose:          true,
		IgnoreWarnings:   false,
		NoPrepro:         false,
		PreproOnly:       false,
		SyntaxOnly:       false,
		Pedantic:         false,
		Asm:              false,
		Object:           false,
		ModuleNum:        0,
		PreproOutputFile: nil,
		OutputFileName:   "test.txt",
		MacroCommands:    nil,
		IncludeDirs:      nil,
	}, input, nil)
	if err != nil {
		t.Fatalf("expected no error %s", err.Error())
	}

	modules := mirCtx.ModuleList()
	if len(modules) != 1 {
		t.Fatalf("expected 1 module, got %d modules", len(modules))
	}

	for i, module := range modules {
		t.Logf("module %s, index %d", module.Name(), i)
	}

	mirCtx.WriteTextual(StderrFile) // need to write to stderr, because go test eats stdout
}
