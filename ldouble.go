package photon

/*
#include <photon-math.h>
*/
import "C"
import "math/big"

var ldblExp = float64(C.LDBL_EXP)

func newLongDouble(v *big.Float) C.photon_ldouble_t {
	mant := new(big.Float)
	exp := v.MantExp(mant)

	mant.Mul(mant, big.NewFloat(ldblExp))
	if mant.Signbit() { // move mantissa signbit to exponent
		exp *= -1
		mant.Neg(mant)
	}

	uintMant, _ := mant.Uint64()
	return &C.struct_photon_ldouble{
		exp:  C.int(exp),
		mant: C.uint64_t(uintMant),
	}
}

func newBigFloat(v C.photon_ldouble_t) *big.Float {
	mant := new(big.Float).SetUint64(uint64(v.mant))
	mant.Quo(mant, big.NewFloat(ldblExp))

	exp := int(v.exp)
	if exp < 0 { // move exponent signbit to mantissa
		exp *= -1
		mant.Neg(mant)
	}

	return new(big.Float).SetMantExp(mant, exp)
}
