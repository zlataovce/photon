package photon

/*
#cgo CFLAGS: -Imir -Imir/c2mir -Iinternal -Wno-psabi
#cgo LDFLAGS: -Lmir -lmir

#include <mir.h>
#include <mir-gen.h>
#include <stdarg.h>
#include <photon-math.c> // include the .c file here to compile it
#include <photon-dlist.h>

DLIST_ACCUMULATOR(modules, MIR_module_t, MIR_context_t, *MIR_get_module_list(item))

extern void photon_handle_error(MIR_error_type_t type, char* msg);
MIR_NO_RETURN void photon_handle_raw_error(MIR_error_type_t type, const char* format, ...) {
  va_list args, args_copy;
  va_start(args, format);
  va_copy(args_copy, args);

  size_t needed = vsnprintf(NULL, 0, format, args) + 1;
  va_end(args);

  char* buffer = malloc(needed);
  if (buffer == NULL) {
    fprintf(stderr, "malloc failed\n");
    exit(-1);
  }

  va_start(args_copy, format);
  vsprintf(buffer, format, args_copy);
  va_end(args_copy);

  photon_handle_error(type, buffer);
  free(buffer);
  exit(-1);
}

void photon_set_error_handler(MIR_context_t ctx) {
  MIR_set_error_func(ctx, photon_handle_raw_error);
}

MIR_op_t photon_new_ldouble_op(MIR_context_t ctx, photon_ldouble_t v) {
  return MIR_new_ldouble_op(ctx, photon_ldouble_unwrap(v));
}
*/
import "C"
import (
	"fmt"
	"math/big"
	"reflect"
	"unsafe"
)

type LinkFunc *[0]byte // void (MIR_item_t)

var (
	LinkFuncInterpInterface      LinkFunc = (*[0]byte)(C.MIR_set_interp_interface)
	LinkFuncGenInterface         LinkFunc = (*[0]byte)(C.MIR_set_gen_interface)
	LinkFuncParallelGenInterface LinkFunc = (*[0]byte)(C.MIR_set_parallel_gen_interface)
	LinkFuncLazyGenInterface     LinkFunc = (*[0]byte)(C.MIR_set_lazy_gen_interface)
)

type ImportResolverFunc *[0]byte // void* (const char*)

// Context is the MIR context.
type Context struct {
	c      C.MIR_context_t
	gen    bool
	closed bool
	used   bool
}

// NewContext creates a new MIR context.
func NewContext() *Context {
	mirCtx := C.MIR_init()
	C.photon_set_error_handler(mirCtx)

	return &Context{c: mirCtx}
}

func (c *Context) checkClosed() {
	if c.closed {
		panic("context is closed")
	}
}

func (c *Context) checkModule() {
	if c.used {
		panic("already building a module")
	}

	c.used = true
}

func (c *Context) Closed() bool {
	return c.closed
}

// Close frees up all memory needed for IR and compilation.
func (c *Context) Close() {
	if !c.closed {
		if c.used {
			panic("couldn't close context, building a module")
		}
		if c.gen {
			C.MIR_gen_finish(c.c)
		}

		C.MIR_finish(c.c)
	}

	c.closed = true
}

func (c *Context) Link(f LinkFunc, r ImportResolverFunc) {
	c.checkClosed()
	if !c.gen && (f == LinkFuncGenInterface || f == LinkFuncParallelGenInterface || f == LinkFuncLazyGenInterface) {
		panic("context is not generation capable")
	}

	C.MIR_link(c.c, f, r)
}

func (c *Context) LoadModule(m *Module) {
	c.checkClosed()
	C.MIR_load_module(c.c, m.c)
}

func (c *Context) LoadExternal(name string, addr unsafe.Pointer) {
	c.checkClosed()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	C.MIR_load_external(c.c, cName, addr)
}

// ModuleList lists all modules in the context.
func (c *Context) ModuleList() []*Module {
	c.checkClosed()
	pModules := C.photon_list_modules(c.c)
	if pModules == nil {
		panic("memory allocation failed")
	}

	defer C.free(unsafe.Pointer(pModules.items))
	defer C.free(unsafe.Pointer(pModules))

	cModules := unsafe.Slice(&pModules.items, int(pModules.length))
	modules := make([]*Module, pModules.length)
	for i, module := range cModules {
		modules[i] = &Module{c: *module}
	}
	return modules
}

func (c *Context) SwapContext(m *Module, newCtx *Context) {
	c.checkClosed()
	newCtx.checkClosed()

	C.MIR_change_module_ctx(c.c, m.c, newCtx.c)
}

// ReadBinary reads a binary MIR representation from a file into the context.
func (c *Context) ReadBinary(f *File) {
	c.checkClosed()
	C.MIR_read(c.c, f.fpOrNil())
}

// WriteBinary writes the context data as a binary MIR representation into a file.
func (c *Context) WriteBinary(f *File) {
	c.checkClosed()
	C.MIR_write(c.c, f.fpOrNil())
}

// WriteBinaryModule writes the module data as a binary MIR representation into a file.
func (c *Context) WriteBinaryModule(f *File, m *Module) {
	c.checkClosed()
	C.MIR_write_module(c.c, f.fpOrNil(), m.c)
}

// ReadTextual reads a textual MIR representation from a string into the context.
func (c *Context) ReadTextual(data string) {
	c.checkClosed()
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.MIR_scan_string(c.c, cData)
}

// WriteTextual writes the context data as a textual MIR representation into a file.
func (c *Context) WriteTextual(f *File) {
	c.checkClosed()
	C.MIR_output(c.c, f.fpOrNil())
}

// WriteTextualModule writes the module data as a textual MIR representation into a file.
func (c *Context) WriteTextualModule(f *File, m *Module) {
	c.checkClosed()
	C.MIR_output_module(c.c, f.fpOrNil(), m.c)
}

// WriteTextualItem writes the item data as a textual MIR representation into a file.
func (c *Context) WriteTextualItem(f *File, i *Item) {
	c.checkClosed()
	C.MIR_output_item(c.c, f.fpOrNil(), i.c)
}

// WriteTextualInsn writes the instruction data as a textual MIR representation into a file.
func (c *Context) WriteTextualInsn(f *File, insn *Insn, func_ *Func, newline bool) {
	c.checkClosed()
	C.MIR_output_insn(c.c, f.fpOrNil(), insn.c, func_.c, C.int(btoi(newline)))
}

// WriteTextualOp writes the operation data as a textual MIR representation into a file.
func (c *Context) WriteTextualOp(f *File, op *Op, func_ *Func) {
	c.checkClosed()
	C.MIR_output_op(c.c, f.fpOrNil(), op.c, func_.c)
}

func (c *Context) TypeName(type_ DataType) string {
	c.checkClosed()
	return C.GoString(C.MIR_type_str(c.c, C.MIR_type_t(type_)))
}

func (c *Context) InsnName(code InsnCode) string {
	c.checkClosed()
	return C.GoString(C.MIR_insn_name(c.c, C.MIR_insn_code_t(code)))
}

func (c *Context) Interpret(item *Item, vals ...*Val) []*Val {
	c.checkClosed()

	fl, ok := item.Item().(FuncLike)
	if !ok {
		panic("item is not a function or a prototype")
	}

	var (
		resNum     = fl.Res()
		argNum     = fl.Args()
		realArgNum = len(vals)
	)
	if argNum != realArgNum {
		panic(fmt.Sprintf("function expected %d arg(s), got %d", argNum, realArgNum))
	}

	var cValsPtr *C.MIR_val_t
	if len(vals) != 0 {
		cVals := make([]C.MIR_val_t, len(vals))
		for idx, v := range vals {
			cVals[idx] = v.c
		}

		cValsPtr = (*C.MIR_val_t)(unsafe.Pointer(unsafe.SliceData(cVals)))
	}

	cResults := make([]C.MIR_val_t, resNum)
	C.MIR_interp_arr(c.c, item.c, (*C.MIR_val_t)(unsafe.Pointer(unsafe.SliceData(cResults))), C.size_t(len(vals)), cValsPtr)

	results := make([]*Val, resNum)
	for idx, r := range cResults {
		results[idx] = &Val{c: r}
	}

	return results
}

// ModuleContext is a MIR context in the process of building a module.
type ModuleContext struct {
	*Context
	*Module

	closed bool
	used   bool
}

func (mc *ModuleContext) checkClosed() {
	mc.Context.checkClosed()
	if mc.closed {
		panic("module context is closed")
	}
}

func (mc *ModuleContext) checkFunc() {
	if mc.used {
		panic("already building a function")
	}

	mc.used = true
}

// NewModuleContext creates a new module-building context.
// name is the module name.
func (c *Context) NewModuleContext(name string) *ModuleContext {
	c.checkClosed()
	c.checkModule()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return &ModuleContext{
		Context: c,
		Module:  &Module{c: C.MIR_new_module(c.c, cName)},
	}
}

func (mc *ModuleContext) ItemName(item *Item) string {
	mc.checkClosed()
	return C.GoString(C.MIR_item_name(mc.Context.c, item.c))
}

func (mc *ModuleContext) ItemFunc(item *Item) *Func {
	mc.checkClosed()
	return &Func{c: C.MIR_get_item_func(mc.Context.c, item.c)}
}

func (mc *ModuleContext) NewImport(name string) *Item {
	mc.checkClosed()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return &Item{c: C.MIR_new_import(mc.Context.c, cName)}
}

func (mc *ModuleContext) NewExport(name string) *Item {
	mc.checkClosed()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return &Item{c: C.MIR_new_export(mc.Context.c, cName)}
}

func (mc *ModuleContext) NewForward(name string) *Item {
	mc.checkClosed()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return &Item{c: C.MIR_new_forward(mc.Context.c, cName)}
}

func (mc *ModuleContext) NewBSS(name string, len int) *Item {
	mc.checkClosed()
	var cName *C.char
	if name != "" { // nil values are acceptable
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
	}

	return &Item{c: C.MIR_new_bss(mc.Context.c, cName, C.size_t(len))}
}

func allocEls(elType DataType, els []interface{}) (cEls unsafe.Pointer) {
	if len(els) != 0 {
		switch elType {
		case DataTypeI8:
			cInt8s := make([]C.int8_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(int8); ok {
					cInt8s[idx] = C.int8_t(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the int8 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cInt8s))
		case DataTypeU8:
			cUint8s := make([]C.uint8_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(uint8); ok {
					cUint8s[idx] = C.uint8_t(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the uint8 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cUint8s))
		case DataTypeI16:
			cInt16s := make([]C.int16_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(int16); ok {
					cInt16s[idx] = C.int16_t(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the int16 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cInt16s))
		case DataTypeU16:
			cUint16s := make([]C.uint16_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(uint16); ok {
					cUint16s[idx] = C.uint16_t(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the uint16 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cUint16s))
		case DataTypeI32:
			cInt32s := make([]C.int32_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(int32); ok {
					cInt32s[idx] = C.int32_t(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the int32 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cInt32s))
		case DataTypeU32:
			cUint32s := make([]C.uint32_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(uint32); ok {
					cUint32s[idx] = C.uint32_t(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the uint32 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cUint32s))
		case DataTypeI64:
			cInt64s := make([]C.int64_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(int); ok {
					cInt64s[idx] = C.int64_t(i0)
				} else if i1, ok := i.(int64); ok {
					cInt64s[idx] = C.int64_t(i1)
				} else {
					panic(fmt.Sprintf("expected all items to be of the int or int64 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cInt64s))
		case DataTypeU64:
			cUint64s := make([]C.uint64_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(uint); ok {
					cUint64s[idx] = C.uint64_t(i0)
				} else if i1, ok := i.(uint64); ok {
					cUint64s[idx] = C.uint64_t(i1)
				} else {
					panic(fmt.Sprintf("expected all items to be of the uint or uint64 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cUint64s))
		case DataTypeF:
			cFloats := make([]C.float, len(els))
			for idx, i := range els {
				if i0, ok := i.(float32); ok {
					cFloats[idx] = C.float(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the float32 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cFloats))
		case DataTypeD:
			cDoubles := make([]C.double, len(els))
			for idx, i := range els {
				if i0, ok := i.(float64); ok {
					cDoubles[idx] = C.double(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the float64 type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cDoubles))
		case DataTypeLD:
			cLDoubles := make([]C.photon_ldouble_t, len(els))
			for idx, i := range els {
				if i0, ok := i.(*big.Float); ok {
					cLDoubles[idx] = newLongDouble(i0)
				} else {
					panic(fmt.Sprintf("expected all items to be of the *big.Float type, found %s", reflect.TypeOf(i).String()))
				}
			}
			cEls = C.photon_ldouble_unwrap_arr((*C.photon_ldouble_t)(unsafe.Pointer(unsafe.SliceData(cLDoubles))), C.size_t(len(els)))
		case DataTypeP:
			cPtrs := make([]unsafe.Pointer, len(els))
			for idx, i := range els {
				if i0, ok := i.(unsafe.Pointer); ok {
					cPtrs[idx] = i0
				} else {
					cPtrs[idx] = unsafe.Pointer(&i)
				}
			}
			cEls = unsafe.Pointer(unsafe.SliceData(cPtrs))
		default:
			panic(fmt.Sprintf("unsupported data type %d", elType))
		}
	}

	return nil
}

func (mc *ModuleContext) NewData(name string, elType DataType, els ...interface{}) *Item {
	mc.checkClosed()
	var cName *C.char
	if name != "" { // nil values are acceptable
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
	}

	cEls := allocEls(elType, els)
	return &Item{c: C.MIR_new_data(mc.Context.c, cName, C.MIR_type_t(elType), C.size_t(len(els)), cEls)}
}

func (mc *ModuleContext) NewStringData(name string, str string) *Item {
	mc.checkClosed()
	var cName *C.char
	if name != "" { // nil values are acceptable
		cName = C.CString(name)
	}

	mirStr := C.MIR_str_t{
		len: C.size_t(len(str)),
		s:   C.CString(str),
	}
	defer C.free(unsafe.Pointer(mirStr.s))

	return &Item{c: C.MIR_new_string_data(mc.Context.c, cName, mirStr)}
}

func (mc *ModuleContext) NewRefData(name string, item *Item, disp int) *Item {
	mc.checkClosed()
	var cName *C.char
	if name != "" { // nil values are acceptable
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
	}

	return &Item{c: C.MIR_new_ref_data(mc.Context.c, cName, item.c, C.int64_t(disp))}
}

func (mc *ModuleContext) NewExprData(name string, exprItem *Item) *Item {
	mc.checkClosed()
	var cName *C.char
	if name != "" { // nil values are acceptable
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
	}

	return &Item{c: C.MIR_new_expr_data(mc.Context.c, cName, exprItem.c)}
}

// FunctionContext is a MIR context in the process of building a function.
type FunctionContext struct {
	*ModuleContext
	*Item

	closed bool
}

func (fc *FunctionContext) checkClosed() {
	fc.ModuleContext.checkClosed()
	if fc.closed {
		panic("function context is closed")
	}
}

func (mc *ModuleContext) NewProto(name string, resType DataType, vars ...*Var) *FunctionContext {
	return mc.NewProtoMultiRes(name, []DataType{resType}, vars...)
}

func (mc *ModuleContext) NewProtoMultiRes(name string, resTypes []DataType, vars ...*Var) *FunctionContext {
	mc.checkClosed()
	mc.checkFunc()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var (
		cResTypesPtr *C.MIR_type_t
		cVarsPtr     *C.MIR_var_t
	)

	if len(resTypes) != 0 {
		cResTypes := make([]C.MIR_type_t, len(resTypes))
		for idx, dt := range resTypes {
			cResTypes[idx] = C.MIR_type_t(dt)
		}

		cResTypesPtr = (*C.MIR_type_t)(unsafe.Pointer(unsafe.SliceData(cResTypes)))
	}

	if len(vars) != 0 {
		cVars := make([]C.MIR_var_t, len(vars))
		for idx, v := range vars {
			mirVar := C.MIR_var_t{
				_type: C.MIR_type_t(v.type_),
				name:  C.CString(v.name),
				size:  C.size_t(v.size),
			}
			defer C.free(unsafe.Pointer(mirVar.name))

			cVars[idx] = mirVar
		}

		cVarsPtr = (*C.MIR_var_t)(unsafe.Pointer(unsafe.SliceData(cVars)))
	}

	return &FunctionContext{
		ModuleContext: mc,
		Item: &Item{
			c: C.MIR_new_proto_arr(
				mc.Context.c,
				cName,
				C.size_t(len(resTypes)),
				cResTypesPtr,
				C.size_t(len(vars)),
				cVarsPtr,
			),
		},
	}
}

func (mc *ModuleContext) NewVarargProto(name string, resType DataType, vars ...*Var) *FunctionContext {
	return mc.NewVarargProtoMultiRes(name, []DataType{resType}, vars...)
}

func (mc *ModuleContext) NewVarargProtoMultiRes(name string, resTypes []DataType, vars ...*Var) *FunctionContext {
	mc.checkClosed()
	mc.checkFunc()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var (
		cResTypesPtr *C.MIR_type_t
		cVarsPtr     *C.MIR_var_t
	)

	if len(resTypes) != 0 {
		cResTypes := make([]C.MIR_type_t, len(resTypes))
		for idx, dt := range resTypes {
			cResTypes[idx] = C.MIR_type_t(dt)
		}

		cResTypesPtr = (*C.MIR_type_t)(unsafe.Pointer(unsafe.SliceData(cResTypes)))
	}

	if len(vars) != 0 {
		cVars := make([]C.MIR_var_t, len(vars))
		for idx, v := range vars {
			mirVar := C.MIR_var_t{
				_type: C.MIR_type_t(v.type_),
				name:  C.CString(v.name),
				size:  C.size_t(v.size),
			}
			defer C.free(unsafe.Pointer(mirVar.name))

			cVars[idx] = mirVar
		}

		cVarsPtr = (*C.MIR_var_t)(unsafe.Pointer(unsafe.SliceData(cVars)))
	}

	return &FunctionContext{
		ModuleContext: mc,
		Item: &Item{
			c: C.MIR_new_vararg_proto_arr(
				mc.Context.c,
				cName,
				C.size_t(len(resTypes)),
				cResTypesPtr,
				C.size_t(len(vars)),
				cVarsPtr,
			),
		},
	}
}

func (mc *ModuleContext) NewFunc(name string, resType DataType, vars ...*Var) *FunctionContext {
	return mc.NewFuncMultiRes(name, []DataType{resType}, vars...)
}

func (mc *ModuleContext) NewFuncMultiRes(name string, resTypes []DataType, vars ...*Var) *FunctionContext {
	mc.checkClosed()
	mc.checkFunc()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var (
		cResTypesPtr *C.MIR_type_t
		cVarsPtr     *C.MIR_var_t
	)

	if len(resTypes) != 0 {
		cResTypes := make([]C.MIR_type_t, len(resTypes))
		for idx, dt := range resTypes {
			cResTypes[idx] = C.MIR_type_t(dt)
		}

		cResTypesPtr = (*C.MIR_type_t)(unsafe.Pointer(unsafe.SliceData(cResTypes)))
	}

	if len(vars) != 0 {
		cVars := make([]C.MIR_var_t, len(vars))
		for idx, v := range vars {
			mirVar := C.MIR_var_t{
				_type: C.MIR_type_t(v.type_),
				name:  C.CString(v.name),
				size:  C.size_t(v.size),
			}
			defer C.free(unsafe.Pointer(mirVar.name))

			cVars[idx] = mirVar
		}

		cVarsPtr = (*C.MIR_var_t)(unsafe.Pointer(unsafe.SliceData(cVars)))
	}

	return &FunctionContext{
		ModuleContext: mc,
		Item: &Item{
			c: C.MIR_new_func_arr(
				mc.Context.c,
				cName,
				C.size_t(len(resTypes)),
				cResTypesPtr,
				C.size_t(len(vars)),
				cVarsPtr,
			),
		},
	}
}

func (mc *ModuleContext) NewVarargFunc(name string, resType DataType, vars ...*Var) *FunctionContext {
	return mc.NewVarargFuncMultiRes(name, []DataType{resType}, vars...)
}

func (mc *ModuleContext) NewVarargFuncMultiRes(name string, resTypes []DataType, vars ...*Var) *FunctionContext {
	mc.checkClosed()
	mc.checkFunc()

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var (
		cResTypesPtr *C.MIR_type_t
		cVarsPtr     *C.MIR_var_t
	)

	if len(resTypes) != 0 {
		cResTypes := make([]C.MIR_type_t, len(resTypes))
		for idx, dt := range resTypes {
			cResTypes[idx] = C.MIR_type_t(dt)
		}

		cResTypesPtr = (*C.MIR_type_t)(unsafe.Pointer(unsafe.SliceData(cResTypes)))
	}

	if len(vars) != 0 {
		cVars := make([]C.MIR_var_t, len(vars))
		for idx, v := range vars {
			mirVar := C.MIR_var_t{
				_type: C.MIR_type_t(v.type_),
				name:  C.CString(v.name),
				size:  C.size_t(v.size),
			}
			defer C.free(unsafe.Pointer(mirVar.name))

			cVars[idx] = mirVar
		}

		cVarsPtr = (*C.MIR_var_t)(unsafe.Pointer(unsafe.SliceData(cVars)))
	}

	return &FunctionContext{
		ModuleContext: mc,
		Item: &Item{
			c: C.MIR_new_vararg_func_arr(
				mc.Context.c,
				cName,
				C.size_t(len(resTypes)),
				cResTypesPtr,
				C.size_t(len(vars)),
				cVarsPtr,
			),
		},
	}
}

func (fc *FunctionContext) NewFuncReg(func_ *Func, type_ DataType, name string) Reg {
	fc.checkClosed()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return Reg(C.MIR_new_func_reg(fc.ModuleContext.Context.c, func_.c, C.MIR_type_t(type_), cName))
}

func allocOps(ops []*Op) *C.MIR_op_t {
	if len(ops) == 0 {
		return nil
	}

	cOps := make([]C.MIR_op_t, len(ops))
	for idx, op := range ops {
		cOps[idx] = op.c
	}

	return (*C.MIR_op_t)(unsafe.Pointer(unsafe.SliceData(cOps)))
}

func (fc *FunctionContext) NewInsn(code InsnCode, ops ...*Op) *Insn {
	fc.checkClosed()
	cOps := allocOps(ops)

	return &Insn{
		c: C.MIR_new_insn_arr(
			fc.ModuleContext.Context.c,
			C.MIR_insn_code_t(code),
			C.size_t(len(ops)),
			cOps,
		),
	}
}

func (fc *FunctionContext) CopyInsn(insn *Insn) *Insn {
	fc.checkClosed()
	return &Insn{c: C.MIR_copy_insn(fc.ModuleContext.Context.c, insn.c)}
}

func (fc *FunctionContext) NewLabel() *Insn {
	fc.checkClosed()
	return &Insn{c: C.MIR_new_label(fc.ModuleContext.Context.c)}
}

func (fc *FunctionContext) Reg(regName string, func_ *Func) Reg {
	fc.checkClosed()
	cRegName := C.CString(regName)
	defer C.free(unsafe.Pointer(cRegName))

	return Reg(C.MIR_reg(fc.ModuleContext.Context.c, cRegName, func_.c))
}

func (fc *FunctionContext) RegType(reg Reg, func_ *Func) DataType {
	fc.checkClosed()
	return DataType(C.MIR_reg_type(fc.ModuleContext.Context.c, C.MIR_reg_t(reg), func_.c))
}

func (fc *FunctionContext) RegName(reg Reg, func_ *Func) string {
	fc.checkClosed()
	return C.GoString(C.MIR_reg_name(fc.ModuleContext.Context.c, C.MIR_reg_t(reg), func_.c))
}

func (fc *FunctionContext) NewRegOp(reg Reg) *Op {
	fc.checkClosed()
	return &Op{c: C.MIR_new_reg_op(fc.ModuleContext.Context.c, C.MIR_reg_t(reg))}
}

func (fc *FunctionContext) NewIntOp(v int) *Op {
	fc.checkClosed()
	return &Op{c: C.MIR_new_int_op(fc.ModuleContext.Context.c, C.int64_t(v))}
}

func (fc *FunctionContext) NewUintOp(v uint) *Op {
	fc.checkClosed()
	return &Op{c: C.MIR_new_uint_op(fc.ModuleContext.Context.c, C.uint64_t(v))}
}

func (fc *FunctionContext) NewFloatOp(v float32) *Op {
	fc.checkClosed()
	return &Op{c: C.MIR_new_float_op(fc.ModuleContext.Context.c, C.float(v))}
}

func (fc *FunctionContext) NewDoubleOp(v float64) *Op {
	fc.checkClosed()
	return &Op{c: C.MIR_new_double_op(fc.ModuleContext.Context.c, C.double(v))}
}

func (fc *FunctionContext) NewLongDoubleOp(ld *big.Float) *Op {
	fc.checkClosed()
	return &Op{c: C.photon_new_ldouble_op(fc.ModuleContext.Context.c, newLongDouble(ld))}
}

func (fc *FunctionContext) NewRefOp(item *Item) *Op {
	fc.checkClosed()
	return &Op{c: C.MIR_new_ref_op(fc.ModuleContext.Context.c, item.c)}
}

func (fc *FunctionContext) NewStrOp(str string) *Op {
	fc.checkClosed()
	mirStr := C.MIR_str_t{
		len: C.size_t(len(str)),
		s:   C.CString(str),
	}
	defer C.free(unsafe.Pointer(mirStr.s))

	return &Op{c: C.MIR_new_str_op(fc.ModuleContext.Context.c, mirStr)}
}

func (fc *FunctionContext) NewLabelOp(label Label) *Op {
	fc.checkClosed()
	return &Op{c: C.MIR_new_label_op(fc.ModuleContext.Context.c, C.MIR_label_t(label.c))}
}

func (fc *FunctionContext) AppendInsn(insn *Insn) {
	fc.checkClosed()
	C.MIR_append_insn(fc.ModuleContext.Context.c, fc.Item.c, insn.c)
}

func (fc *FunctionContext) PrependInsn(insn *Insn) {
	fc.checkClosed()
	C.MIR_prepend_insn(fc.ModuleContext.Context.c, fc.Item.c, insn.c)
}

func (fc *FunctionContext) InsertInsnAfter(after *Insn, insn *Insn) {
	fc.checkClosed()
	C.MIR_insert_insn_after(fc.ModuleContext.Context.c, fc.Item.c, after.c, insn.c)
}

func (fc *FunctionContext) InsertInsnBefore(before *Insn, insn *Insn) {
	fc.checkClosed()
	C.MIR_insert_insn_before(fc.ModuleContext.Context.c, fc.Item.c, before.c, insn.c)
}

func (fc *FunctionContext) RemoveInsn(insn *Insn) {
	fc.checkClosed()
	C.MIR_remove_insn(fc.ModuleContext.Context.c, fc.Item.c, insn.c)
}

// Close finishes building the function.
func (fc *FunctionContext) Close() {
	if !fc.closed {
		C.MIR_finish_func(fc.ModuleContext.Context.c)
	}

	fc.closed = true
	fc.ModuleContext.used = false
}

// Close finishes building the used.
func (mc *ModuleContext) Close() {
	if !mc.closed {
		if mc.used {
			panic("couldn't close context, building a function")
		}

		C.MIR_finish_module(mc.Context.c)
	}

	mc.closed = true
	mc.Context.used = false
}
