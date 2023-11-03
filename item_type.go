package photon

type ItemType uint

const (
	ItemTypeFunc ItemType = iota
	ItemTypeProto
	ItemTypeImport
	ItemTypeExport
	ItemTypeForward
	ItemTypeData
	ItemTypeRefData
	ItemTypeExprData
	ItemTypeBss
)

const ItemTypeUnknown ItemType = 32767

var (
	enumToItemType = map[ItemType]string{
		ItemTypeFunc:     "func",
		ItemTypeProto:    "proto",
		ItemTypeImport:   "import",
		ItemTypeExport:   "export",
		ItemTypeForward:  "forward",
		ItemTypeData:     "data",
		ItemTypeRefData:  "ref_data",
		ItemTypeExprData: "expr_data",
		ItemTypeBss:      "bss",
	}

	// reverse enumToItemType
	itemTypeToEnum = make(map[string]ItemType, len(enumToItemType))
)

func init() {
	for k, v := range enumToItemType {
		itemTypeToEnum[v] = k
	}
}

func NewItemType(name string) ItemType {
	if itemType, ok := itemTypeToEnum[name]; ok {
		return itemType
	}

	return ItemTypeUnknown
}

func (it ItemType) Name() string {
	if name, ok := enumToItemType[it]; ok {
		return name
	}

	return ""
}
