package photon

type ErrorType uint

const (
	ErrorTypeNone ErrorType = iota
	ErrorTypeSyntax
	ErrorTypeBinaryIo
	ErrorTypeAlloc
	ErrorTypeFinish
	ErrorTypeNoModule
	ErrorTypeNestedModule
	ErrorTypeNoFunc
	ErrorTypeFunc
	ErrorTypeVarargFunc
	ErrorTypeNestedFunc
	ErrorTypeWrongParamValue
	ErrorTypeReservedName
	ErrorTypeImportExport
	ErrorTypeUndeclaredFuncReg
	ErrorTypeRepeatedDecl
	ErrorTypeRegType
	ErrorTypeWrongType
	ErrorTypeUniqueReg
	ErrorTypeUndeclaredOpRef
	ErrorTypeOpsNum
	ErrorTypeCallOp
	ErrorTypeUnspecOp
	ErrorTypeRet
	ErrorTypeOpMode
	ErrorTypeOutOp
	ErrorTypeInvalidInsn
	ErrorTypeCtxChange
	ErrorTypeParallel
)

const ErrorTypeUnknown ErrorType = 32767

var (
	enumToErrorType = map[ErrorType]string{
		ErrorTypeNone:              "none",
		ErrorTypeSyntax:            "syntax",
		ErrorTypeBinaryIo:          "binary_io",
		ErrorTypeAlloc:             "alloc",
		ErrorTypeFinish:            "finish",
		ErrorTypeNoModule:          "no_module",
		ErrorTypeNestedModule:      "nested_module",
		ErrorTypeNoFunc:            "no_func",
		ErrorTypeFunc:              "func",
		ErrorTypeVarargFunc:        "vararg_func",
		ErrorTypeNestedFunc:        "nested_func",
		ErrorTypeWrongParamValue:   "wrong_param_value",
		ErrorTypeReservedName:      "reserved_name",
		ErrorTypeImportExport:      "import_export",
		ErrorTypeUndeclaredFuncReg: "undeclared_func_reg",
		ErrorTypeRepeatedDecl:      "repeated_decl",
		ErrorTypeRegType:           "reg_type",
		ErrorTypeWrongType:         "wrong_type",
		ErrorTypeUniqueReg:         "unique_reg",
		ErrorTypeUndeclaredOpRef:   "undeclared_op_ref",
		ErrorTypeOpsNum:            "ops_num",
		ErrorTypeCallOp:            "call_op",
		ErrorTypeUnspecOp:          "unspec_op",
		ErrorTypeRet:               "ret",
		ErrorTypeOpMode:            "op_mode",
		ErrorTypeOutOp:             "out_op",
		ErrorTypeInvalidInsn:       "invalid_insn",
		ErrorTypeCtxChange:         "ctx_change",
		ErrorTypeParallel:          "parallel",
	}

	// reverse enumToErrorType
	errorTypeToEnum = make(map[string]ErrorType, len(enumToErrorType))
)

func init() {
	for k, v := range enumToErrorType {
		errorTypeToEnum[v] = k
	}
}

func NewErrorType(name string) ErrorType {
	if errType, ok := errorTypeToEnum[name]; ok {
		return errType
	}

	return ErrorTypeUnknown
}

func (et ErrorType) Name() string {
	if name, ok := enumToErrorType[et]; ok {
		return name
	}

	return ""
}
