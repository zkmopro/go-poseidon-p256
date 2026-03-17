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
