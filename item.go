package photon

/*
#include <mir.h>
*/
import "C"
import "unsafe"

type Item struct {
	c C.MIR_item_t
}

func (i *Item) Module() *Module {
	return &Module{c: i.c.module}
}

func (i *Item) Type() ItemType {
	return ItemType(i.c.item_type)
}

func (i *Item) Item() interface{} {
	addr := unsafe.Pointer(&(i.c.u[i.c.item_type]))
	switch i.Type() {
	case ItemTypeFunc:
		return &Func{c: *(*C.MIR_func_t)(addr)}
	case ItemTypeProto:
		return &Proto{c: *(*C.MIR_proto_t)(addr)}
	case ItemTypeImport, ItemTypeExport, ItemTypeForward:
		return C.GoString(*(**C.char)(addr))
	case ItemTypeData:
		return &Data{c: *(*C.MIR_data_t)(addr)}
	case ItemTypeRefData:
		return &RefData{c: *(*C.MIR_ref_data_t)(addr)}
	case ItemTypeExprData:
		return &ExprData{c: *(*C.MIR_expr_data_t)(addr)}
	case ItemTypeBss:
		return &BSS{c: *(*C.MIR_bss_t)(addr)}
	}

	return nil
}

func (i *Item) Func() *Func {
	if func_, ok := i.Item().(*Func); ok {
		return func_
	}

	return nil
}

func (i *Item) Proto() *Proto {
	if proto, ok := i.Item().(*Proto); ok {
		return proto
	}

	return nil
}

func (i *Item) ImportExportForward() string {
	if ief, ok := i.Item().(string); ok {
		return ief
	}

	return ""
}

func (i *Item) Data() *Data {
	if data, ok := i.Item().(*Data); ok {
		return data
	}

	return nil
}

func (i *Item) RefData() *RefData {
	if refData, ok := i.Item().(*RefData); ok {
		return refData
	}

	return nil
}

func (i *Item) ExprData() *ExprData {
	if exprData, ok := i.Item().(*ExprData); ok {
		return exprData
	}

	return nil
}

func (i *Item) BSS() *BSS {
	if bss, ok := i.Item().(*BSS); ok {
		return bss
	}

	return nil
}
