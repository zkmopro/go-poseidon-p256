package poseidon

import "math/big"

// ORDER is the base field order of secp256r1 (P-256).
var ORDER = func() *big.Int {
	n, _ := new(big.Int).SetString("FFFFFFFF00000001000000000000000000000000FFFFFFFFFFFFFFFFFFFFFFFF", 16)
	return n
}()

// BITS is the bit length of ORDER.
const BITS = 256

var bigZero = new(big.Int)
var bigOne = big.NewInt(1)

// fpAdd returns (a + b) mod ORDER.
func fpAdd(a, b *big.Int) *big.Int {
	r := new(big.Int).Add(a, b)
	return r.Mod(r, ORDER)
}

// fpMul returns (a * b) mod ORDER.
func fpMul(a, b *big.Int) *big.Int {
	r := new(big.Int).Mul(a, b)
	return r.Mod(r, ORDER)
}

// fpMulN returns a * b WITHOUT modular reduction.
func fpMulN(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

// fpSqrN returns a * a WITHOUT modular reduction.
func fpSqrN(a *big.Int) *big.Int {
	return new(big.Int).Mul(a, a)
}

// fpNeg returns (-a) mod ORDER.
func fpNeg(a *big.Int) *big.Int {
	r := new(big.Int).Neg(a)
	return r.Mod(r, ORDER)
}

// fpInv returns the modular inverse of a mod ORDER using extended GCD.
// Matches modular.ts:invert() exactly.
func fpInv(number *big.Int) *big.Int {
	if number.Sign() == 0 {
		panic("fpInv: cannot invert zero")
	}
	a := new(big.Int).Mod(number, ORDER)
	b := new(big.Int).Set(ORDER)
	x := new(big.Int)
	y := big.NewInt(1)
	u := big.NewInt(1)
	v := new(big.Int)

	for a.Sign() != 0 {
		q := new(big.Int).Div(b, a)
		r := new(big.Int).Mod(b, a)
		m := new(big.Int).Sub(x, new(big.Int).Mul(u, q))
		n := new(big.Int).Sub(y, new(big.Int).Mul(v, q))
		b.Set(a)
		a.Set(r)
		x.Set(u)
		y.Set(v)
		u.Set(m)
		v.Set(n)
	}
	return x.Mod(x, ORDER)
}

// fpCreate normalizes a value into [0, ORDER).
func fpCreate(a *big.Int) *big.Int {
	return new(big.Int).Mod(a, ORDER)
}

// fpInvertBatch inverts a slice of field elements using Montgomery's trick.
// Matches FpInvertBatch in modular.ts.
func fpInvertBatch(nums []*big.Int) []*big.Int {
	tmp := make([]*big.Int, len(nums))
	acc := new(big.Int).Set(bigOne)
	for i, num := range nums {
		if num.Sign() == 0 {
			continue
		}
		tmp[i] = new(big.Int).Set(acc)
		acc = fpMul(acc, num)
	}
	inv := fpInv(acc)
	for i := len(nums) - 1; i >= 0; i-- {
		if nums[i].Sign() == 0 {
			continue
		}
		tmp[i] = fpMul(inv, tmp[i])
		inv = fpMul(inv, nums[i])
	}
	return tmp
}

// fpSbox5 computes x^5 mod ORDER.
// Matches: Fp.mul(Fp.sqrN(Fp.sqrN(n)), n)
func fpSbox5(n *big.Int) *big.Int {
	return fpMul(fpSqrN(fpSqrN(n)), n)
}
