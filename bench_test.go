package poseidon

import (
	"crypto/rand"
	"math/big"
	"testing"
)

// randFieldElement returns a random element in [0, ORDER).
func randFieldElement() *big.Int {
	n, err := rand.Int(rand.Reader, ORDER)
	if err != nil {
		panic(err)
	}
	return n
}

func BenchmarkHash2(b *testing.B) {
	a := randFieldElement()
	c := randFieldElement()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash2(a, c)
	}
}

func BenchmarkHash3(b *testing.B) {
	a := randFieldElement()
	c := randFieldElement()
	d := randFieldElement()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash3(a, c, d)
	}
}
