package poseidon

import "math/big"

type poseidonConfig struct {
	t              int
	roundsFull     int
	roundsPartial  int
	roundConstants [][]*big.Int
	mds            [][]*big.Int
}

var configT3 *poseidonConfig
var configT4 *poseidonConfig

func init() {
	configT3 = loadConfigT3()
	configT4 = loadConfigT4()
}

func parseHex(s string) *big.Int {
	n, ok := new(big.Int).SetString(s[2:], 16) // strip "0x"
	if !ok {
		panic("invalid hex: " + s)
	}
	return n
}

func loadConfigT3() *poseidonConfig {
	rounds := 8 + 57
	rc := make([][]*big.Int, rounds)
	for i := 0; i < rounds; i++ {
		row := make([]*big.Int, 3)
		for j := 0; j < 3; j++ {
			row[j] = parseHex(rcHexT3[i][j])
		}
		rc[i] = row
	}
	mds := make([][]*big.Int, 3)
	for i := 0; i < 3; i++ {
		row := make([]*big.Int, 3)
		for j := 0; j < 3; j++ {
			row[j] = parseHex(mdsHexT3[i][j])
		}
		mds[i] = row
	}
	return &poseidonConfig{
		t: 3, roundsFull: 8, roundsPartial: 57,
		roundConstants: rc, mds: mds,
	}
}

func loadConfigT4() *poseidonConfig {
	rounds := 8 + 56
	rc := make([][]*big.Int, rounds)
	for i := 0; i < rounds; i++ {
		row := make([]*big.Int, 4)
		for j := 0; j < 4; j++ {
			row[j] = parseHex(rcHexT4[i][j])
		}
		rc[i] = row
	}
	mds := make([][]*big.Int, 4)
	for i := 0; i < 4; i++ {
		row := make([]*big.Int, 4)
		for j := 0; j < 4; j++ {
			row[j] = parseHex(mdsHexT4[i][j])
		}
		mds[i] = row
	}
	return &poseidonConfig{
		t: 4, roundsFull: 8, roundsPartial: 56,
		roundConstants: rc, mds: mds,
	}
}

// poseidonRound applies one round of the Poseidon permutation.
func poseidonRound(cfg *poseidonConfig, values []*big.Int, isFull bool, roundIdx int) []*big.Int {
	t := cfg.t

	// Add round constants
	for j := 0; j < t; j++ {
		values[j] = fpAdd(values[j], cfg.roundConstants[roundIdx][j])
	}

	// S-box
	if isFull {
		for j := 0; j < t; j++ {
			values[j] = fpSbox5(values[j])
		}
	} else {
		values[0] = fpSbox5(values[0])
	}

	// MDS matrix multiplication
	// Matches: mds.map(i => i.reduce((acc, i, j) => Fp.add(acc, Fp.mulN(i, values[j])), Fp.ZERO))
	newValues := make([]*big.Int, t)
	for i := 0; i < t; i++ {
		acc := new(big.Int) // Fp.ZERO
		for j := 0; j < t; j++ {
			acc = fpAdd(acc, fpMulN(cfg.mds[i][j], values[j]))
		}
		newValues[i] = acc
	}

	return newValues
}

// poseidonPermutation applies the full Poseidon permutation.
func poseidonPermutation(cfg *poseidonConfig, values []*big.Int) []*big.Int {
	halfFull := cfg.roundsFull / 2
	lastRound := 0

	for i := 0; i < halfFull; i++ {
		values = poseidonRound(cfg, values, true, lastRound)
		lastRound++
	}
	for i := 0; i < cfg.roundsPartial; i++ {
		values = poseidonRound(cfg, values, false, lastRound)
		lastRound++
	}
	for i := 0; i < halfFull; i++ {
		values = poseidonRound(cfg, values, true, lastRound)
		lastRound++
	}

	return values
}

// Hash2 computes the Poseidon hash of two field elements.
// state = [0, a, b], t=3, returns state[0].
func Hash2(a, b *big.Int) *big.Int {
	state := []*big.Int{
		new(big.Int),
		fpCreate(a),
		fpCreate(b),
	}
	result := poseidonPermutation(configT3, state)
	return result[0]
}

// Hash3 computes the Poseidon hash of three field elements.
// state = [0, a, b, c], t=4, returns state[0].
func Hash3(a, b, c *big.Int) *big.Int {
	state := []*big.Int{
		new(big.Int),
		fpCreate(a),
		fpCreate(b),
		fpCreate(c),
	}
	result := poseidonPermutation(configT4, state)
	return result[0]
}

// Hash computes the Poseidon hash of arbitrary field elements.
// t = len(inputs) + 1, state = [0, inputs...], returns state[0].
func Hash(inputs []*big.Int) *big.Int {
	switch len(inputs) {
	case 2:
		return Hash2(inputs[0], inputs[1])
	case 3:
		return Hash3(inputs[0], inputs[1], inputs[2])
	default:
		panic("unsupported number of inputs: only 2 or 3 are supported")
	}
}
