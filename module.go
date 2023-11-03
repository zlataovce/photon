package photon

/*
#include <mir.h>
#include <photon-dlist.h>

DLIST_ACCUMULATOR(module_items, MIR_item_t, MIR_module_t, item->items)
*/
import "C"
import (
	"unsafe"
)

type Module struct {
	c C.MIR_module_t
}

func (m *Module) Data() unsafe.Pointer {
	return m.c.data
}

func (m *Module) Name() string {
	return C.GoString(m.c.name)
}

func (m *Module) Items() []*Item {
	pModuleItems := C.photon_list_module_items(m.c)
	if pModuleItems == nil {
		panic("memory allocation failed")
	}

	defer C.free(unsafe.Pointer(pModuleItems.items))
	defer C.free(unsafe.Pointer(pModuleItems))

	cModuleItems := unsafe.Slice(&pModuleItems.items, int(pModuleItems.length))
	moduleItems := make([]*Item, pModuleItems.length)
	for i, moduleItem := range cModuleItems {
		moduleItems[i] = &Item{c: *moduleItem}
	}
	return moduleItems
}
