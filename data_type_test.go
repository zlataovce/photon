package photon

import "testing"

func TestDataType_IsInteger(t *testing.T) {
	if !DataTypeI8.IsInteger() {
		t.Fatalf("expected DataTypeI8 to be an integer type")
	}

	if !DataTypeU64.IsInteger() {
		t.Fatalf("expected DataTypeU64 to be an integer type")
	}
}

func TestDataType_IsFloatingPoint(t *testing.T) {
	if !DataTypeF.IsFloatingPoint() {
		t.Fatalf("expected DataTypeF to be a floating-point type")
	}

	if !DataTypeD.IsFloatingPoint() {
		t.Fatalf("expected DataTypeD to be a floating-point type")
	}
}

func TestDataType_IsBlk(t *testing.T) {
	if !DataTypeBLK.IsBlk() {
		t.Fatalf("expected DataTypeBLK to be a block type")
	}

	if DataTypeRBLK.IsBlk() {
		t.Fatalf("expected DataTypeRBLK to not be a block type")
	}
}
