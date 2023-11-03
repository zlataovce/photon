package photon

/*
#include <mir.h>
*/
import "C"

type DataType uint

const (
	DataTypeI8 DataType = iota
	DataTypeU8
	DataTypeI16
	DataTypeU16
	DataTypeI32
	DataTypeU32
	DataTypeI64
	DataTypeU64
	DataTypeF
	DataTypeD
	DataTypeLD
	DataTypeP
	DataTypeBLK
	DataTypeRBLK DataType = 4 + iota
	DataTypeUndef
	DataTypeBound
)

const DataTypeUnknown DataType = 32767

var (
	enumToDataType = map[DataType]string{
		DataTypeI8:    "i8",
		DataTypeU8:    "u8",
		DataTypeI16:   "i16",
		DataTypeU16:   "u16",
		DataTypeI32:   "i32",
		DataTypeU32:   "u32",
		DataTypeI64:   "i64",
		DataTypeU64:   "u64",
		DataTypeF:     "f",
		DataTypeD:     "d",
		DataTypeLD:    "ld",
		DataTypeP:     "p",
		DataTypeBLK:   "blk",
		DataTypeRBLK:  "rblk",
		DataTypeUndef: "undef",
		DataTypeBound: "bound",
	}

	// reverse enumToDataType
	dataTypeToEnum = make(map[string]DataType, len(enumToDataType))
)

func init() {
	for k, v := range enumToDataType {
		dataTypeToEnum[v] = k
	}
}

func NewDataType(name string) DataType {
	if dataType, ok := dataTypeToEnum[name]; ok {
		return dataType
	}

	return DataTypeUnknown
}

func (dt DataType) Name() string {
	if name, ok := enumToDataType[dt]; ok {
		return name
	}

	return ""
}

func (dt DataType) IsInteger() bool {
	return itob(int(C.MIR_int_type_p(C.MIR_type_t(dt))))
}

func (dt DataType) IsFloatingPoint() bool {
	return itob(int(C.MIR_fp_type_p(C.MIR_type_t(dt))))
}

func (dt DataType) IsBlk() bool {
	return itob(int(C.MIR_blk_type_p(C.MIR_type_t(dt))))
}
