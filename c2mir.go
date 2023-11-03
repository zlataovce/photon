package photon

/*
#include <c2mir.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

var ErrCompileFailed = errors.New("compilation failed")

// MacroCommandType is a type of preprocessor macro command.
type MacroCommandType int

const (
	// MacroCommandTypeUndef defines a #undef preprocessor macro command.
	MacroCommandTypeUndef MacroCommandType = iota
	// MacroCommandTypeDefine defines a #define preprocessor macro command.
	MacroCommandTypeDefine
)

// MacroCommand is a preprocessor macro declaration command.
type MacroCommand struct {
	// Type is the macro command type.
	Type MacroCommandType
	// Name is the macro name.
	Name string
	// Def is the macro definition content, used only when Type is MacroCommandTypeDefine.
	Def string
}

// Options are C2MIR compilation options.
type Options struct {
	MessageFile      *File
	Debug            bool
	Verbose          bool
	IgnoreWarnings   bool
	NoPrepro         bool
	PreproOnly       bool
	SyntaxOnly       bool
	Pedantic         bool
	Asm              bool
	Object           bool
	ModuleNum        int
	PreproOutputFile *File // non-nil for prepro_only_p
	OutputFileName   string
	MacroCommands    []*MacroCommand
	IncludeDirs      []string
}

// C2MIRContext is a C2MIR compilation capable Context.
type C2MIRContext struct {
	*Context

	closed bool
}

func (c2m *C2MIRContext) checkClosed() {
	c2m.Context.checkClosed()
	if c2m.closed {
		panic("c2mir context is closed")
	}
}

// NewC2MIRContext initializes the C2MIR compiler in the supplied context and wraps it.
func NewC2MIRContext(c *Context) *C2MIRContext {
	C.c2mir_init(c.c)

	return &C2MIRContext{Context: c}
}

// Compile compiles C source code into MIR.
func (c2m *C2MIRContext) Compile(
	ops *Options,
	data Input,
	outputFile *File,
) error {
	defer data.Close()

	c2m.checkClosed()
	cOps := C.struct_c2mir_options{
		message_file:       ops.MessageFile.fpOrNil(),
		debug_p:            C.int(btoi(ops.Debug)),
		verbose_p:          C.int(btoi(ops.Verbose)),
		ignore_warnings_p:  C.int(btoi(ops.IgnoreWarnings)),
		no_prepro_p:        C.int(btoi(ops.NoPrepro)),
		prepro_only_p:      C.int(btoi(ops.PreproOnly)),
		syntax_only_p:      C.int(btoi(ops.SyntaxOnly)),
		pedantic_p:         C.int(btoi(ops.Pedantic)),
		asm_p:              C.int(btoi(ops.Asm)),
		object_p:           C.int(btoi(ops.Object)),
		module_num:         C.size_t(ops.ModuleNum),
		prepro_output_file: ops.PreproOutputFile.fpOrNil(),
		output_file_name:   C.CString(ops.OutputFileName),
		macro_commands_num: C.size_t(len(ops.MacroCommands)),
		include_dirs_num:   C.size_t(len(ops.IncludeDirs)),
	}

	defer C.free(unsafe.Pointer(cOps.output_file_name))

	if len(ops.MacroCommands) != 0 {
		cMacroCommands := make([]C.struct_c2mir_macro_command, len(ops.MacroCommands))
		for idx, mc := range ops.MacroCommands {
			cMc := C.struct_c2mir_macro_command{
				def_p: C.int(int(mc.Type)),
				name:  C.CString(mc.Name),
			}
			defer C.free(unsafe.Pointer(cMc.name))

			if mc.Type == MacroCommandTypeDefine {
				cMc.def = C.CString(mc.Def)
				defer C.free(unsafe.Pointer(cMc.def))
			}

			cMacroCommands[idx] = cMc
		}

		cOps.macro_commands = (*C.struct_c2mir_macro_command)(unsafe.Pointer(unsafe.SliceData(cMacroCommands)))
	}

	if len(ops.IncludeDirs) != 0 {
		cIncludeDirs := make([]*C.char, len(ops.IncludeDirs))
		for idx, dir := range ops.IncludeDirs {
			cDir := C.CString(dir)
			defer C.free(unsafe.Pointer(cDir))

			cIncludeDirs[idx] = cDir
		}

		cOps.include_dirs = (**C.char)(unsafe.Pointer(unsafe.SliceData(cIncludeDirs)))
	}

	cSourceName := C.CString(data.SourceName())
	defer C.free(unsafe.Pointer(cSourceName))

	// non-zero is returned for successful compilation
	if itob(int(C.c2mir_compile(c2m.c, &cOps, data.ReadFunc(), data.Data(), cSourceName, outputFile.fpOrNil()))) {
		return nil
	}

	return ErrCompileFailed
}

func (c2m *C2MIRContext) Close() {
	if !c2m.closed {
		C.c2mir_finish(c2m.c)
	}

	c2m.closed = true
}
