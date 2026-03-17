package poseidon

import "math/big"

// grainLFSR implements the Grain LFSR used for Poseidon constant generation.
// Matches poseidon.ts:grainLFSR exactly.
type grainLFSR struct {
	state [80]int
	pos   int
}

func newGrainLFSR(state [80]int) *grainLFSR {
	g := &grainLFSR{state: state}
	// Warm up: 160 getBit calls
	for i := 0; i < 160; i++ {
		g.getBit()
	}
	return g
}

func (g *grainLFSR) getBit() int {
	bit := g.state[(g.pos+62)%80] ^
		g.state[(g.pos+51)%80] ^
		g.state[(g.pos+38)%80] ^
		g.state[(g.pos+23)%80] ^
		g.state[(g.pos+13)%80] ^
		g.state[(g.pos+0)%80]
	g.state[g.pos] = bit
	g.pos = (g.pos + 1) % 80
	return bit
}

// shrink implements the shrinking generator.
func (g *grainLFSR) shrink() int {
	for {
		b1 := g.getBit()
		b2 := g.getBit()
		if b1 == 1 {
			return b2
		}
	}
}

// poseidonGrain creates a Grain LFSR sampler for Poseidon constant generation.
// Returns a function sample(count, reject) that produces field elements.
func poseidonGrain(t, roundsFull, roundsPartial int) func(count int, reject bool) []*big.Int {
	var state [80]int
	for i := range state {
		state[i] = 1
	}

	pos := 0
	writeBits := func(value *big.Int, bitCount int) {
		for i := bitCount - 1; i >= 0; i-- {
			state[pos] = int(value.Bit(i))
			pos++
		}
	}

	writeBits(big.NewInt(1), 2)                    // prime field: [1,0]
	writeBits(big.NewInt(0), 4)                    // sbox not inverse: [0,0,0,0]
	writeBits(big.NewInt(int64(BITS)), 12)         // field bits
	writeBits(big.NewInt(int64(t)), 12)            // state size
	writeBits(big.NewInt(int64(roundsFull)), 10)   // full rounds
	writeBits(big.NewInt(int64(roundsPartial)), 10) // partial rounds

	lfsr := newGrainLFSR(state)

	return func(count int, reject bool) []*big.Int {
		res := make([]*big.Int, 0, count)
		for i := 0; i < count; i++ {
			for {
				num := new(big.Int)
				for j := 0; j < BITS; j++ {
					num.Lsh(num, 1)
					if lfsr.shrink() == 1 {
						num.Or(num, bigOne)
					}
				}
				if reject && num.Cmp(ORDER) >= 0 {
					continue
				}
				res = append(res, fpCreate(num))
				break
			}
		}
		return res
	}
}

// GenConstants generates Poseidon round constants and MDS matrix
// for the given parameters using the Grain LFSR.
func GenConstants(t, roundsFull, roundsPartial int) (roundConstants [][]*big.Int, mds [][]*big.Int) {
	rounds := roundsFull + roundsPartial
	sample := poseidonGrain(t, roundsFull, roundsPartial)

	// Generate round constants with rejection sampling
	roundConstants = make([][]*big.Int, rounds)
	for r := 0; r < rounds; r++ {
		roundConstants[r] = sample(t, true)
	}

	// Generate MDS matrix elements without rejection sampling
	xs := sample(t, false)
	ys := sample(t, false)

	// Construct MDS matrix: M[i][j] = 1 / (xs[i] + ys[j])
	mds = make([][]*big.Int, t)
	for i := 0; i < t; i++ {
		row := make([]*big.Int, t)
		for j := 0; j < t; j++ {
			row[j] = fpAdd(xs[i], ys[j])
		}
		inverted := fpInvertBatch(row)
		mds[i] = inverted
	}

	return roundConstants, mds
}
