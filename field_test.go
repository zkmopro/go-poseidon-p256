package poseidon

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"
)

type fieldVectors struct {
	Order string `json:"order"`
	Add   []struct {
		A, B, Result string
	} `json:"add"`
	Mul []struct {
		A, B, Result string
	} `json:"mul"`
	Neg []struct {
		A, Result string
	} `json:"neg"`
	Inv []struct {
		A, Result string
	} `json:"inv"`
	Sbox5 []struct {
		Input, Result string
	} `json:"sbox5"`
	Create []struct {
		Input, Result string
	} `json:"create"`
	MulN []struct {
		A, B, Result string
	} `json:"mulN"`
	SqrN []struct {
		A, Result string
	} `json:"sqrN"`
	InvertBatch struct {
		Inputs  []string `json:"inputs"`
		Results []string `json:"results"`
	} `json:"invertBatch"`
}

func loadFieldVectors(t *testing.T) *fieldVectors {
	t.Helper()
	data, err := os.ReadFile("testdata/field_vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var v fieldVectors
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	return &v
}

func hexToBig(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 0)
	if !ok {
		panic("invalid hex: " + s)
	}
	return n
}

func TestFpAdd(t *testing.T) {
	v := loadFieldVectors(t)
	for _, tc := range v.Add {
		a, b, want := hexToBig(tc.A), hexToBig(tc.B), hexToBig(tc.Result)
		got := fpAdd(a, b)
		if got.Cmp(want) != 0 {
			t.Errorf("fpAdd(%s, %s) = %s, want %s", tc.A, tc.B, got.Text(16), want.Text(16))
		}
	}
}

func TestFpMul(t *testing.T) {
	v := loadFieldVectors(t)
	for _, tc := range v.Mul {
		a, b, want := hexToBig(tc.A), hexToBig(tc.B), hexToBig(tc.Result)
		got := fpMul(a, b)
		if got.Cmp(want) != 0 {
			t.Errorf("fpMul(%s, %s) = %s, want %s", tc.A, tc.B, got.Text(16), want.Text(16))
		}
	}
}

func TestFpNeg(t *testing.T) {
	v := loadFieldVectors(t)
	for _, tc := range v.Neg {
		a, want := hexToBig(tc.A), hexToBig(tc.Result)
		got := fpNeg(a)
		if got.Cmp(want) != 0 {
			t.Errorf("fpNeg(%s) = %s, want %s", tc.A, got.Text(16), want.Text(16))
		}
	}
}

func TestFpInv(t *testing.T) {
	v := loadFieldVectors(t)
	for _, tc := range v.Inv {
		a, want := hexToBig(tc.A), hexToBig(tc.Result)
		got := fpInv(a)
		if got.Cmp(want) != 0 {
			t.Errorf("fpInv(%s) = %s, want %s", tc.A, got.Text(16), want.Text(16))
		}
	}
}

func TestFpSbox5(t *testing.T) {
	v := loadFieldVectors(t)
	for _, tc := range v.Sbox5 {
		n, want := hexToBig(tc.Input), hexToBig(tc.Result)
		got := fpSbox5(n)
		if got.Cmp(want) != 0 {
			t.Errorf("fpSbox5(%s) = %s, want %s", tc.Input, got.Text(16), want.Text(16))
		}
	}
}

func TestFpCreate(t *testing.T) {
	v := loadFieldVectors(t)
	for _, tc := range v.Create {
		input, want := hexToBig(tc.Input), hexToBig(tc.Result)
		got := fpCreate(input)
		if got.Cmp(want) != 0 {
			t.Errorf("fpCreate(%s) = %s, want %s", tc.Input, got.Text(16), want.Text(16))
		}
	}
}

func TestFpMulN(t *testing.T) {
	v := loadFieldVectors(t)
	for _, tc := range v.MulN {
		a, b, want := hexToBig(tc.A), hexToBig(tc.B), hexToBig(tc.Result)
		got := fpMulN(a, b)
		if got.Cmp(want) != 0 {
			t.Errorf("fpMulN(%s, %s) = %s, want %s", tc.A, tc.B, got.Text(16), want.Text(16))
		}
	}
}

func TestFpSqrN(t *testing.T) {
	v := loadFieldVectors(t)
	for _, tc := range v.SqrN {
		a, want := hexToBig(tc.A), hexToBig(tc.Result)
		got := fpSqrN(a)
		if got.Cmp(want) != 0 {
			t.Errorf("fpSqrN(%s) = %s, want %s", tc.A, got.Text(16), want.Text(16))
		}
	}
}

func TestFpInvertBatch(t *testing.T) {
	v := loadFieldVectors(t)
	inputs := make([]*big.Int, len(v.InvertBatch.Inputs))
	for i, s := range v.InvertBatch.Inputs {
		inputs[i] = hexToBig(s)
	}
	results := fpInvertBatch(inputs)
	for i, s := range v.InvertBatch.Results {
		want := hexToBig(s)
		if results[i].Cmp(want) != 0 {
			t.Errorf("fpInvertBatch[%d] = %s, want %s", i, results[i].Text(16), want.Text(16))
		}
	}
}

func TestFpIdentities(t *testing.T) {
	a := hexToBig("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	// a + 0 = a
	if fpAdd(a, big.NewInt(0)).Cmp(a) != 0 {
		t.Error("a + 0 != a")
	}
	// a * 1 = a
	if fpMul(a, big.NewInt(1)).Cmp(a) != 0 {
		t.Error("a * 1 != a")
	}
	// a * inv(a) = 1
	if fpMul(a, fpInv(a)).Cmp(big.NewInt(1)) != 0 {
		t.Error("a * inv(a) != 1")
	}
	// a + neg(a) = 0
	if fpAdd(a, fpNeg(a)).Sign() != 0 {
		t.Error("a + neg(a) != 0")
	}
}

func TestFpInvertBatchConsistency(t *testing.T) {
	vals := []*big.Int{big.NewInt(7), big.NewInt(13), big.NewInt(42)}
	batch := fpInvertBatch(vals)
	for i, v := range vals {
		single := fpInv(v)
		if batch[i].Cmp(single) != 0 {
			t.Errorf("InvertBatch[%d] = %s, Inv = %s", i, batch[i].Text(16), single.Text(16))
		}
	}
}

func TestFpAlgebraicProperties(t *testing.T) {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	small := big.NewInt(42)
	large := hexToBig("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	orderMinus1 := new(big.Int).Sub(ORDER, one)
	elems := []*big.Int{zero, one, small, large, orderMinus1}

	t.Run("Commutativity", func(t *testing.T) {
		for _, a := range elems {
			for _, b := range elems {
				if fpAdd(a, b).Cmp(fpAdd(b, a)) != 0 {
					t.Errorf("fpAdd not commutative for %s, %s", a.Text(16), b.Text(16))
				}
				if fpMul(a, b).Cmp(fpMul(b, a)) != 0 {
					t.Errorf("fpMul not commutative for %s, %s", a.Text(16), b.Text(16))
				}
			}
		}
	})

	t.Run("Associativity", func(t *testing.T) {
		triples := [][3]*big.Int{
			{one, small, large},
			{small, large, orderMinus1},
			{one, orderMinus1, small},
		}
		for _, tr := range triples {
			a, b, c := tr[0], tr[1], tr[2]
			if fpAdd(fpAdd(a, b), c).Cmp(fpAdd(a, fpAdd(b, c))) != 0 {
				t.Errorf("fpAdd not associative for %s, %s, %s", a.Text(16), b.Text(16), c.Text(16))
			}
			if fpMul(fpMul(a, b), c).Cmp(fpMul(a, fpMul(b, c))) != 0 {
				t.Errorf("fpMul not associative for %s, %s, %s", a.Text(16), b.Text(16), c.Text(16))
			}
		}
	})

	t.Run("Distributivity", func(t *testing.T) {
		triples := [][3]*big.Int{
			{one, small, large},
			{small, large, orderMinus1},
		}
		for _, tr := range triples {
			a, b, c := tr[0], tr[1], tr[2]
			lhs := fpMul(a, fpAdd(b, c))
			rhs := fpAdd(fpMul(a, b), fpMul(a, c))
			if lhs.Cmp(rhs) != 0 {
				t.Errorf("distributivity failed for %s, %s, %s", a.Text(16), b.Text(16), c.Text(16))
			}
		}
	})

	t.Run("DoubleNegation", func(t *testing.T) {
		for _, a := range elems {
			want := fpCreate(a)
			if fpNeg(fpNeg(a)).Cmp(want) != 0 {
				t.Errorf("fpNeg(fpNeg(%s)) != %s", a.Text(16), want.Text(16))
			}
		}
	})

	t.Run("DoubleInversion", func(t *testing.T) {
		for _, a := range elems {
			if a.Sign() == 0 {
				continue
			}
			want := fpCreate(a)
			if fpInv(fpInv(a)).Cmp(want) != 0 {
				t.Errorf("fpInv(fpInv(%s)) != %s", a.Text(16), want.Text(16))
			}
		}
	})

	t.Run("AdditiveInverse", func(t *testing.T) {
		for _, a := range elems {
			if fpAdd(a, fpNeg(a)).Sign() != 0 {
				t.Errorf("fpAdd(%s, fpNeg(%s)) != 0", a.Text(16), a.Text(16))
			}
		}
	})

	t.Run("MultiplicativeInverse", func(t *testing.T) {
		for _, a := range elems {
			if a.Sign() == 0 {
				continue
			}
			if fpMul(a, fpInv(a)).Cmp(big.NewInt(1)) != 0 {
				t.Errorf("fpMul(%s, fpInv(%s)) != 1", a.Text(16), a.Text(16))
			}
		}
	})

	t.Run("Sbox5Consistency", func(t *testing.T) {
		for _, a := range elems {
			// x^5 via repeated multiplication
			x2 := fpMul(a, a)
			x4 := fpMul(x2, x2)
			x5 := fpMul(x4, a)
			if fpSbox5(a).Cmp(x5) != 0 {
				t.Errorf("fpSbox5(%s) != x^5 via repeated mul", a.Text(16))
			}
		}
	})
}

func TestFpInvPanicsOnZero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for fpInv(0)")
		}
	}()
	fpInv(big.NewInt(0))
}

func TestFpCreateNegativeInput(t *testing.T) {
	neg := big.NewInt(-1)
	got := fpCreate(neg)
	if got.Sign() < 0 || got.Cmp(ORDER) >= 0 {
		t.Errorf("fpCreate(-1) = %s, not in [0, ORDER)", got.Text(16))
	}
	// -1 mod ORDER should equal ORDER-1
	want := new(big.Int).Sub(ORDER, big.NewInt(1))
	if got.Cmp(want) != 0 {
		t.Errorf("fpCreate(-1) = %s, want %s", got.Text(16), want.Text(16))
	}

	neg2 := big.NewInt(-42)
	got2 := fpCreate(neg2)
	if got2.Sign() < 0 || got2.Cmp(ORDER) >= 0 {
		t.Errorf("fpCreate(-42) = %s, not in [0, ORDER)", got2.Text(16))
	}
}

func TestFpInvertBatchEdgeCases(t *testing.T) {
	t.Run("EmptySlice", func(t *testing.T) {
		result := fpInvertBatch([]*big.Int{})
		if len(result) != 0 {
			t.Errorf("expected empty result, got length %d", len(result))
		}
	})

	t.Run("SingleElement", func(t *testing.T) {
		val := big.NewInt(7)
		batch := fpInvertBatch([]*big.Int{val})
		single := fpInv(val)
		if batch[0].Cmp(single) != 0 {
			t.Errorf("single-element batch: got %s, want %s", batch[0].Text(16), single.Text(16))
		}
	})

	t.Run("LargerBatch", func(t *testing.T) {
		vals := make([]*big.Int, 12)
		for i := range vals {
			vals[i] = big.NewInt(int64(i + 1))
		}
		batch := fpInvertBatch(vals)
		for i, v := range vals {
			prod := fpMul(v, batch[i])
			if prod.Cmp(big.NewInt(1)) != 0 {
				t.Errorf("batch[%d] * original != 1", i)
			}
		}
	})

	t.Run("ContainsZero", func(t *testing.T) {
		// Zero inputs are skipped by fpInvertBatch; result[i] is nil for zero inputs.
		vals := []*big.Int{big.NewInt(7), big.NewInt(0), big.NewInt(13)}
		batch := fpInvertBatch(vals)
		if batch[1] != nil {
			t.Error("expected nil for zero input")
		}
		if batch[0].Cmp(fpInv(big.NewInt(7))) != 0 {
			t.Error("batch[0] wrong")
		}
		if batch[2].Cmp(fpInv(big.NewInt(13))) != 0 {
			t.Error("batch[2] wrong")
		}
	})
}
