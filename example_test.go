package poseidon_test

import (
	"fmt"
	"math/big"

	poseidon "github.com/zkmopro/go-poseidon-p256"
)

func ExampleHash2() {
	a, _ := new(big.Int).SetString("deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef", 16)
	b, _ := new(big.Int).SetString("cafebabedeadcafecafebabedeadcafecafebabedeadcafecafebabedeadcafe", 16)
	h := poseidon.Hash2(a, b)
	fmt.Printf("0x%s\n", h.Text(16))
	// Output:
	// 0x66ac887f89cc1740dab07d27b8fe70c153a2728c4a69bc73a457677bca1ee5c7
}

func ExampleHash3() {
	a, _ := new(big.Int).SetString("a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90", 16)
	b, _ := new(big.Int).SetString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", 16)
	c, _ := new(big.Int).SetString("fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321", 16)
	h := poseidon.Hash3(a, b, c)
	fmt.Printf("0x%s\n", h.Text(16))
	// Output:
	// 0x3aa328978565e3e21352ab59da9386b222c739c23436e59535af2aaaf507b417
}
