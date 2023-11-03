package photon

/*
#include <mir.h>
*/
import "C"

// InsnCode is an instruction opcode.
type InsnCode uint

const (
	// InsnCodeMOV is a 64-bit integer move instruction, uses two operands.
	InsnCodeMOV InsnCode = iota
	// InsnCodeFMOV is a single precision floating point value move instruction, uses two operands.
	InsnCodeFMOV
	// InsnCodeDMOV is a double precision floating point value move instruction, uses two operands.
	InsnCodeDMOV
	// InsnCodeLDMOV is a long double floating point value move instruction, uses two operands.
	InsnCodeLDMOV

	// InsnCodeEXT8 is a lower 8-bit input part extension sign instruction, uses two operands.
	InsnCodeEXT8
	// InsnCodeEXT16 is a lower 16-bit input part extension sign instruction, uses two operands.
	InsnCodeEXT16
	// InsnCodeEXT32 is a lower 32-bit input part extension sign instruction, uses two operands.
	InsnCodeEXT32
	// InsnCodeUEXT8 is a lower 8-bit input part extension zero instruction, uses two operands.
	InsnCodeUEXT8
	// InsnCodeUEXT16 is a lower 16-bit input part extension zero instruction, uses two operands.
	InsnCodeUEXT16
	// InsnCodeUEXT32 is a lower 32-bit input part extension zero instruction, uses two operands.
	InsnCodeUEXT32

	InsnCodeI2F
	InsnCodeI2D
	InsnCodeI2LD
	InsnCodeUI2F
	InsnCodeUI2D
	InsnCodeUI2LD
	InsnCodeF2I
	InsnCodeD2I
	InsnCodeLD2I
	InsnCodeF2D
	InsnCodeF2LD
	InsnCodeD2F
	InsnCodeD2LD
	InsnCodeLD2F
	InsnCodeLD2D

	// InsnCodeNEG is a 64-bit integer sign change instruction, uses two operands.
	InsnCodeNEG
	// InsnCodeNEGS is a 32-bit integer sign change instruction, uses two operands.
	InsnCodeNEGS
	InsnCodeFNEG
	InsnCodeDNEG
	InsnCodeLDNEG

	// InsnCodeADD is a 64-bit integer addition instruction, uses three operands.
	InsnCodeADD
	// InsnCodeADDS is a 32-bit integer addition instruction, uses three operands.
	InsnCodeADDS
	InsnCodeFADD
	InsnCodeDADD
	InsnCodeLDADD

	// InsnCodeSUB is a 64-bit integer subtraction instruction, uses three operands.
	InsnCodeSUB
	// InsnCodeSUBS is a 32-bit integer subtraction instruction, uses three operands.
	InsnCodeSUBS
	InsnCodeFSUB
	InsnCodeDSUB
	InsnCodeLDSUB

	// InsnCodeMUL is a 64-bit integer un/signed multiplication instruction, uses three operands.
	InsnCodeMUL
	// InsnCodeMULS is a 32-bit integer un/signed multiplication instruction, uses three operands.
	InsnCodeMULS
	InsnCodeFMUL
	InsnCodeDMUL
	InsnCodeLDMUL

	// InsnCodeDIV is a 64-bit integer signed division instruction, uses three operands.
	InsnCodeDIV
	// InsnCodeDIVS is a 32-bit integer signed division instruction, uses three operands.
	InsnCodeDIVS
	// InsnCodeUDIV is a 64-bit integer unsigned division instruction, uses three operands.
	InsnCodeUDIV
	// InsnCodeUDIVS is a 32-bit integer unsigned division instruction, uses three operands.
	InsnCodeUDIVS
	InsnCodeFDIV
	InsnCodeDDIV
	InsnCodeLDDIV

	// InsnCodeMOD is a 64-bit signed modulo instruction, uses three operands.
	InsnCodeMOD
	// InsnCodeMODS is a 32-bit signed modulo instruction, uses three operands.
	InsnCodeMODS
	// InsnCodeUMOD is a 64-bit unsigned modulo instruction, uses three operands.
	InsnCodeUMOD
	// InsnCodeUMODS is a 32-bit unsigned modulo instruction, uses three operands.
	InsnCodeUMODS

	// InsnCodeAND is a 64-bit integer bitwise AND instruction, uses three operands.
	InsnCodeAND
	// InsnCodeANDS is a 32-bit integer bitwise AND instruction, uses three operands.
	InsnCodeANDS

	// InsnCodeOR is a 64-bit integer bitwise OR instruction, uses three operands.
	InsnCodeOR
	// InsnCodeORS is a 32-bit integer bitwise OR instruction, uses three operands.
	InsnCodeORS

	// InsnCodeXOR is a 64-bit integer bitwise XOR instruction, uses three operands.
	InsnCodeXOR
	// InsnCodeXORS is a 32-bit integer bitwise XOR instruction, uses three operands.
	InsnCodeXORS

	InsnCodeLSH
	InsnCodeLSHS

	InsnCodeRSH
	InsnCodeRSHS

	InsnCodeURSH
	InsnCodeURSHS

	InsnCodeEQ
	InsnCodeEQS
	InsnCodeFEQ
	InsnCodeDEQ
	InsnCodeLDEQ

	InsnCodeNE
	InsnCodeNES
	InsnCodeFNE
	InsnCodeDNE
	InsnCodeLDNE

	InsnCodeLT
	InsnCodeLTS
	InsnCodeULT
	InsnCodeULTS
	InsnCodeFLT
	InsnCodeDLT
	InsnCodeLDLT

	InsnCodeLE
	InsnCodeLES
	InsnCodeULE
	InsnCodeULES
	InsnCodeFLE
	InsnCodeDLE
	InsnCodeLDLE

	InsnCodeGT
	InsnCodeGTS
	InsnCodeUGT
	InsnCodeUGTS
	InsnCodeFGT
	InsnCodeDGT
	InsnCodeLDGT

	InsnCodeGE
	InsnCodeGES
	InsnCodeUGE
	InsnCodeUGES
	InsnCodeFGE
	InsnCodeDGE
	InsnCodeLDGE

	InsnCodeJMP

	InsnCodeBT
	InsnCodeBTS
	InsnCodeBF
	InsnCodeBFS

	InsnCodeBEQ
	InsnCodeBEQS
	InsnCodeFBEQ
	InsnCodeDBEQ
	InsnCodeLDBEQ

	InsnCodeBNE
	InsnCodeBNES
	InsnCodeFBNE
	InsnCodeDBNE
	InsnCodeLDBNE

	InsnCodeBLT
	InsnCodeBLTS
	InsnCodeUBLT
	InsnCodeUBLTS
	InsnCodeFBLT
	InsnCodeDBLT
	InsnCodeLDBLT

	InsnCodeBLE
	InsnCodeBLES
	InsnCodeUBLE
	InsnCodeUBLES
	InsnCodeFBLE
	InsnCodeDBLE
	InsnCodeLDBLE

	InsnCodeBGT
	InsnCodeBGTS
	InsnCodeUBGT
	InsnCodeUBGTS
	InsnCodeFBGT
	InsnCodeDBGT
	InsnCodeLDBGT

	InsnCodeBGE
	InsnCodeBGES
	InsnCodeUBGE
	InsnCodeUBGES
	InsnCodeFBGE
	InsnCodeDBGE
	InsnCodeLDBGE

	InsnCodeCALL
	InsnCodeINLINE
	InsnCodeSWITCH
	InsnCodeRET
	InsnCodeALLOCA
	InsnCodeBSTART
	InsnCodeBEND
	InsnCodeVaArg
	InsnCodeVaBlockArg
	InsnCodeVaStart
	InsnCodeVaEnd
	InsnCodeLABEL
	InsnCodeUNSPEC
	InsnCodePHI
	InsnCodeInvalidInsn
	InsnCodeInsnBound
)

const InsnCodeUnknown InsnCode = 32767

var (
	enumToInsnCode = map[InsnCode]string{
		InsnCodeMOV:         "MOV",
		InsnCodeFMOV:        "FMOV",
		InsnCodeDMOV:        "DMOV",
		InsnCodeLDMOV:       "LDMOV",
		InsnCodeEXT8:        "EXT8",
		InsnCodeEXT16:       "EXT16",
		InsnCodeEXT32:       "EXT32",
		InsnCodeUEXT8:       "UEXT8",
		InsnCodeUEXT16:      "UEXT16",
		InsnCodeUEXT32:      "UEXT32",
		InsnCodeI2F:         "I2F",
		InsnCodeI2D:         "I2D",
		InsnCodeI2LD:        "I2LD",
		InsnCodeUI2F:        "UI2F",
		InsnCodeUI2D:        "UI2D",
		InsnCodeUI2LD:       "UI2LD",
		InsnCodeF2I:         "F2I",
		InsnCodeD2I:         "D2I",
		InsnCodeLD2I:        "LD2I",
		InsnCodeF2D:         "F2D",
		InsnCodeF2LD:        "F2LD",
		InsnCodeD2F:         "D2F",
		InsnCodeD2LD:        "D2LD",
		InsnCodeLD2F:        "LD2F",
		InsnCodeLD2D:        "LD2D",
		InsnCodeNEG:         "NEG",
		InsnCodeNEGS:        "NEGS",
		InsnCodeFNEG:        "FNEG",
		InsnCodeDNEG:        "DNEG",
		InsnCodeLDNEG:       "LDNEG",
		InsnCodeADD:         "ADD",
		InsnCodeADDS:        "ADDS",
		InsnCodeFADD:        "FADD",
		InsnCodeDADD:        "DADD",
		InsnCodeLDADD:       "LDADD",
		InsnCodeSUB:         "SUB",
		InsnCodeSUBS:        "SUBS",
		InsnCodeFSUB:        "FSUB",
		InsnCodeDSUB:        "DSUB",
		InsnCodeLDSUB:       "LDSUB",
		InsnCodeMUL:         "MUL",
		InsnCodeMULS:        "MULS",
		InsnCodeFMUL:        "FMUL",
		InsnCodeDMUL:        "DMUL",
		InsnCodeLDMUL:       "LDMUL",
		InsnCodeDIV:         "DIV",
		InsnCodeDIVS:        "DIVS",
		InsnCodeUDIV:        "UDIV",
		InsnCodeUDIVS:       "UDIVS",
		InsnCodeFDIV:        "FDIV",
		InsnCodeDDIV:        "DDIV",
		InsnCodeLDDIV:       "LDDIV",
		InsnCodeMOD:         "MOD",
		InsnCodeMODS:        "MODS",
		InsnCodeUMOD:        "UMOD",
		InsnCodeUMODS:       "UMODS",
		InsnCodeAND:         "AND",
		InsnCodeANDS:        "ANDS",
		InsnCodeOR:          "OR",
		InsnCodeORS:         "ORS",
		InsnCodeXOR:         "XOR",
		InsnCodeXORS:        "XORS",
		InsnCodeLSH:         "LSH",
		InsnCodeLSHS:        "LSHS",
		InsnCodeRSH:         "RSH",
		InsnCodeRSHS:        "RSHS",
		InsnCodeURSH:        "URSH",
		InsnCodeURSHS:       "URSHS",
		InsnCodeEQ:          "EQ",
		InsnCodeEQS:         "EQS",
		InsnCodeFEQ:         "FEQ",
		InsnCodeDEQ:         "DEQ",
		InsnCodeLDEQ:        "LDEQ",
		InsnCodeNE:          "NE",
		InsnCodeNES:         "NES",
		InsnCodeFNE:         "FNE",
		InsnCodeDNE:         "DNE",
		InsnCodeLDNE:        "LDNE",
		InsnCodeLT:          "LT",
		InsnCodeLTS:         "LTS",
		InsnCodeULT:         "ULT",
		InsnCodeULTS:        "ULTS",
		InsnCodeFLT:         "FLT",
		InsnCodeDLT:         "DLT",
		InsnCodeLDLT:        "LDLT",
		InsnCodeLE:          "LE",
		InsnCodeLES:         "LES",
		InsnCodeULE:         "ULE",
		InsnCodeULES:        "ULES",
		InsnCodeFLE:         "FLE",
		InsnCodeDLE:         "DLE",
		InsnCodeLDLE:        "LDLE",
		InsnCodeGT:          "GT",
		InsnCodeGTS:         "GTS",
		InsnCodeUGT:         "UGT",
		InsnCodeUGTS:        "UGTS",
		InsnCodeFGT:         "FGT",
		InsnCodeDGT:         "DGT",
		InsnCodeLDGT:        "LDGT",
		InsnCodeGE:          "GE",
		InsnCodeGES:         "GES",
		InsnCodeUGE:         "UGE",
		InsnCodeUGES:        "UGES",
		InsnCodeFGE:         "FGE",
		InsnCodeDGE:         "DGE",
		InsnCodeLDGE:        "LDGE",
		InsnCodeJMP:         "JMP",
		InsnCodeBT:          "BT",
		InsnCodeBTS:         "BTS",
		InsnCodeBF:          "BF",
		InsnCodeBFS:         "BFS",
		InsnCodeBEQ:         "BEQ",
		InsnCodeBEQS:        "BEQS",
		InsnCodeFBEQ:        "FBEQ",
		InsnCodeDBEQ:        "DBEQ",
		InsnCodeLDBEQ:       "LDBEQ",
		InsnCodeBNE:         "BNE",
		InsnCodeBNES:        "BNES",
		InsnCodeFBNE:        "FBNE",
		InsnCodeDBNE:        "DBNE",
		InsnCodeLDBNE:       "LDBNE",
		InsnCodeBLT:         "BLT",
		InsnCodeBLTS:        "BLTS",
		InsnCodeUBLT:        "UBLT",
		InsnCodeUBLTS:       "UBLTS",
		InsnCodeFBLT:        "FBLT",
		InsnCodeDBLT:        "DBLT",
		InsnCodeLDBLT:       "LDBLT",
		InsnCodeBLE:         "BLE",
		InsnCodeBLES:        "BLES",
		InsnCodeUBLE:        "UBLE",
		InsnCodeUBLES:       "UBLES",
		InsnCodeFBLE:        "FBLE",
		InsnCodeDBLE:        "DBLE",
		InsnCodeLDBLE:       "LDBLE",
		InsnCodeBGT:         "BGT",
		InsnCodeBGTS:        "BGTS",
		InsnCodeUBGT:        "UBGT",
		InsnCodeUBGTS:       "UBGTS",
		InsnCodeFBGT:        "FBGT",
		InsnCodeDBGT:        "DBGT",
		InsnCodeLDBGT:       "LDBGT",
		InsnCodeBGE:         "BGE",
		InsnCodeBGES:        "BGES",
		InsnCodeUBGE:        "UBGE",
		InsnCodeUBGES:       "UBGES",
		InsnCodeFBGE:        "FBGE",
		InsnCodeDBGE:        "DBGE",
		InsnCodeLDBGE:       "LDBGE",
		InsnCodeCALL:        "CALL",
		InsnCodeINLINE:      "INLINE",
		InsnCodeSWITCH:      "SWITCH",
		InsnCodeRET:         "RET",
		InsnCodeALLOCA:      "ALLOCA",
		InsnCodeBSTART:      "BSTART",
		InsnCodeBEND:        "BEND",
		InsnCodeVaArg:       "VA_ARG",
		InsnCodeVaBlockArg:  "VA_BLOCK_ARG",
		InsnCodeVaStart:     "VA_START",
		InsnCodeVaEnd:       "VA_END",
		InsnCodeLABEL:       "LABEL",
		InsnCodeUNSPEC:      "UNSPEC",
		InsnCodePHI:         "PHI",
		InsnCodeInvalidInsn: "INVALID_INSN",
		InsnCodeInsnBound:   "INSN_BOUND",
	}

	// reverse enumToInsnCode
	insnCodeToEnum = make(map[string]InsnCode, len(enumToInsnCode))
)

func init() {
	for k, v := range enumToInsnCode {
		insnCodeToEnum[v] = k
	}
}

func NewInsnCode(name string) InsnCode {
	if insnCode, ok := insnCodeToEnum[name]; ok {
		return insnCode
	}

	return InsnCodeUnknown
}

func (ic InsnCode) Name() string {
	if name, ok := enumToInsnCode[ic]; ok {
		return name
	}

	return ""
}

func (ic InsnCode) IsCall() bool {
	return itob(int(C.MIR_call_code_p(C.MIR_insn_code_t(ic))))
}

func (ic InsnCode) IsBranch() bool {
	return itob(int(C.MIR_branch_code_p(C.MIR_insn_code_t(ic))))
}

func (ic InsnCode) IsIntegerBranch() bool {
	return itob(int(C.MIR_int_branch_code_p(C.MIR_insn_code_t(ic))))
}

func (ic InsnCode) IsFloatingPointBranch() bool {
	return itob(int(C.MIR_FP_branch_code_p(C.MIR_insn_code_t(ic))))
}

func (ic InsnCode) ReverseBranch() InsnCode {
	return InsnCode(C.MIR_reverse_branch_code(C.MIR_insn_code_t(ic)))
}
