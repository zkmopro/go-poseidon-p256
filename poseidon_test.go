package poseidon

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"
)

type hashVectors struct {
	Hash2 []struct {
		Inputs []string `json:"inputs"`
		Result string   `json:"result"`
	} `json:"hash2"`
	Hash3 []struct {
		Inputs []string `json:"inputs"`
		Result string   `json:"result"`
	} `json:"hash3"`
	RoundTrace struct {
		Inputs []string   `json:"inputs"`
		Trace  [][]string `json:"trace"`
	} `json:"roundTrace"`
	SMTTree struct {
		Leaves        []string `json:"leaves"`
		InternalNodes []string `json:"internalNodes"`
		Root          string   `json:"root"`
	} `json:"smtTree"`
	Normalization struct {
		NormalResult   string `json:"normalResult"`
		OverflowResult string `json:"overflowResult"`
	} `json:"normalization"`
}

func loadHashVectors(t *testing.T) *hashVectors {
	t.Helper()
	data, err := os.ReadFile("testdata/hash_vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var v hashVectors
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	return &v
}

func TestHash2(t *testing.T) {
	v := loadHashVectors(t)
	for i, tc := range v.Hash2 {
		a := hexToBig(tc.Inputs[0])
		b := hexToBig(tc.Inputs[1])
		want := hexToBig(tc.Result)
		got := Hash2(a, b)
		if got.Cmp(want) != 0 {
			t.Errorf("Hash2 case %d: got %s, want %s", i, got.Text(16), want.Text(16))
		}
	}
}

func TestHash3(t *testing.T) {
	v := loadHashVectors(t)
	for i, tc := range v.Hash3 {
		a := hexToBig(tc.Inputs[0])
		b := hexToBig(tc.Inputs[1])
		c := hexToBig(tc.Inputs[2])
		want := hexToBig(tc.Result)
		got := Hash3(a, b, c)
		if got.Cmp(want) != 0 {
			t.Errorf("Hash3 case %d: got %s, want %s", i, got.Text(16), want.Text(16))
		}
	}
}

func TestRoundByRound(t *testing.T) {
	v := loadHashVectors(t)
	trace := v.RoundTrace

	// Set up initial state
	values := make([]*big.Int, 3)
	for i, s := range trace.Inputs {
		values[i] = fpCreate(hexToBig(s))
	}

	cfg := configT3
	halfFull := cfg.roundsFull / 2
	rounds := cfg.roundsFull + cfg.roundsPartial

	for round := 0; round < rounds; round++ {
		isFull := round < halfFull || round >= halfFull+cfg.roundsPartial
		values = poseidonRound(cfg, values, isFull, round)

		for j := 0; j < 3; j++ {
			want := hexToBig(trace.Trace[round][j])
			if values[j].Cmp(want) != 0 {
				t.Fatalf("round %d, element %d: got %s, want %s",
					round, j, values[j].Text(16), want.Text(16))
			}
		}
	}
}

func TestSMTTree(t *testing.T) {
	v := loadHashVectors(t)
	smt := v.SMTTree

	// Verify leaf hashes
	for i := 0; i < 4; i++ {
		leaf := Hash3(big.NewInt(int64(i)), big.NewInt(int64(i+1)), big.NewInt(1))
		want := hexToBig(smt.Leaves[i])
		if leaf.Cmp(want) != 0 {
			t.Fatalf("leaf[%d]: got %s, want %s", i, leaf.Text(16), want.Text(16))
		}
	}

	// Verify internal nodes
	l0 := hexToBig(smt.Leaves[0])
	l1 := hexToBig(smt.Leaves[1])
	l2 := hexToBig(smt.Leaves[2])
	l3 := hexToBig(smt.Leaves[3])

	n01 := Hash2(l0, l1)
	wantN01 := hexToBig(smt.InternalNodes[0])
	if n01.Cmp(wantN01) != 0 {
		t.Fatalf("n01: got %s, want %s", n01.Text(16), wantN01.Text(16))
	}

	n23 := Hash2(l2, l3)
	wantN23 := hexToBig(smt.InternalNodes[1])
	if n23.Cmp(wantN23) != 0 {
		t.Fatalf("n23: got %s, want %s", n23.Text(16), wantN23.Text(16))
	}

	root := Hash2(n01, n23)
	wantRoot := hexToBig(smt.Root)
	if root.Cmp(wantRoot) != 0 {
		t.Fatalf("root: got %s, want %s", root.Text(16), wantRoot.Text(16))
	}
}

func TestInputNormalization(t *testing.T) {
	v := loadHashVectors(t)
	// Hash2(ORDER+1, ORDER) should equal Hash2(1, 0)
	orderPlus1 := new(big.Int).Add(ORDER, big.NewInt(1))
	got := Hash2(orderPlus1, new(big.Int).Set(ORDER))
	want := hexToBig(v.Normalization.NormalResult)
	if got.Cmp(want) != 0 {
		t.Errorf("normalization: got %s, want %s", got.Text(16), want.Text(16))
	}
}

func TestDeterminism(t *testing.T) {
	a := big.NewInt(42)
	b := big.NewInt(99)
	first := Hash2(a, b)
	for i := 0; i < 100; i++ {
		got := Hash2(a, b)
		if got.Cmp(first) != 0 {
			t.Fatalf("non-deterministic at iteration %d", i)
		}
	}
}

func TestNoMutation(t *testing.T) {
	a := big.NewInt(42)
	b := big.NewInt(99)
	aCopy := new(big.Int).Set(a)
	bCopy := new(big.Int).Set(b)
	Hash2(a, b)
	if a.Cmp(aCopy) != 0 {
		t.Error("Hash2 mutated input a")
	}
	if b.Cmp(bCopy) != 0 {
		t.Error("Hash2 mutated input b")
	}
}

func TestHashPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid input length")
		}
	}()
	Hash([]*big.Int{big.NewInt(1)}) // only 1 input - should panic
}

func TestHashGeneric(t *testing.T) {
	a := big.NewInt(1)
	b := big.NewInt(2)
	c := big.NewInt(3)

	// Hash with 2 inputs should match Hash2
	h2 := Hash([]*big.Int{a, b})
	if h2.Cmp(Hash2(a, b)) != 0 {
		t.Error("Hash(2 inputs) != Hash2")
	}

	// Hash with 3 inputs should match Hash3
	h3 := Hash([]*big.Int{a, b, c})
	if h3.Cmp(Hash3(a, b, c)) != 0 {
		t.Error("Hash(3 inputs) != Hash3")
	}
}

func TestHash2NonCommutative(t *testing.T) {
	pairs := [][2]*big.Int{
		{big.NewInt(1), big.NewInt(2)},
		{big.NewInt(42), big.NewInt(99)},
		{big.NewInt(0), big.NewInt(1)},
	}
	for _, p := range pairs {
		a, b := p[0], p[1]
		if Hash2(a, b).Cmp(Hash2(b, a)) == 0 {
			t.Errorf("Hash2(%s, %s) == Hash2(%s, %s), expected different", a, b, b, a)
		}
	}
}

func TestHash3NonCommutative(t *testing.T) {
	a, b, c := big.NewInt(1), big.NewInt(2), big.NewInt(3)
	h := Hash3(a, b, c)
	// Swapping any pair should produce a different hash
	if h.Cmp(Hash3(b, a, c)) == 0 {
		t.Error("Hash3(a,b,c) == Hash3(b,a,c)")
	}
	if h.Cmp(Hash3(a, c, b)) == 0 {
		t.Error("Hash3(a,b,c) == Hash3(a,c,b)")
	}
	if h.Cmp(Hash3(c, b, a)) == 0 {
		t.Error("Hash3(a,b,c) == Hash3(c,b,a)")
	}
}

func TestHashAvalanche(t *testing.T) {
	a := big.NewInt(42)
	b := big.NewInt(99)
	h1 := Hash2(a, b)

	// Flip one bit in input a
	aFlipped := new(big.Int).Xor(a, big.NewInt(1))
	h2 := Hash2(aFlipped, b)

	// Count differing bits
	diff := new(big.Int).Xor(h1, h2)
	bitsChanged := 0
	for _, word := range diff.Bits() {
		for word != 0 {
			bitsChanged++
			word &= word - 1
		}
	}

	// Expect at least 25% of 256 bits to differ (weak avalanche check)
	if bitsChanged < 64 {
		t.Errorf("avalanche: only %d/256 bits changed, expected at least 64", bitsChanged)
	}
}

func TestNoMutationHash3(t *testing.T) {
	a := big.NewInt(10)
	b := big.NewInt(20)
	c := big.NewInt(30)
	aCopy := new(big.Int).Set(a)
	bCopy := new(big.Int).Set(b)
	cCopy := new(big.Int).Set(c)
	Hash3(a, b, c)
	if a.Cmp(aCopy) != 0 {
		t.Error("Hash3 mutated input a")
	}
	if b.Cmp(bCopy) != 0 {
		t.Error("Hash3 mutated input b")
	}
	if c.Cmp(cCopy) != 0 {
		t.Error("Hash3 mutated input c")
	}
}

func TestHashPanicCases(t *testing.T) {
	cases := []struct {
		name   string
		inputs []*big.Int
	}{
		{"OneInput", []*big.Int{big.NewInt(1)}},
		{"FourInputs", []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4)}},
		{"ZeroInputs", []*big.Int{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for Hash with %d inputs", len(tc.inputs))
				}
			}()
			Hash(tc.inputs)
		})
	}
}

func TestHashOutputInField(t *testing.T) {
	inputs := [][2]*big.Int{
		{big.NewInt(0), big.NewInt(0)},
		{big.NewInt(1), big.NewInt(2)},
		{new(big.Int).Sub(ORDER, big.NewInt(1)), new(big.Int).Sub(ORDER, big.NewInt(2))},
	}
	for _, pair := range inputs {
		h := Hash2(pair[0], pair[1])
		if h.Sign() < 0 || h.Cmp(ORDER) >= 0 {
			t.Errorf("Hash2 output %s not in [0, ORDER)", h.Text(16))
		}
	}
	// Also check Hash3
	h3 := Hash3(big.NewInt(1), big.NewInt(2), big.NewInt(3))
	if h3.Sign() < 0 || h3.Cmp(ORDER) >= 0 {
		t.Errorf("Hash3 output %s not in [0, ORDER)", h3.Text(16))
	}
}
