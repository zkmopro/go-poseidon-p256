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
	// 0x20390b3870f5480a3fd0cc5ba71ba0c2930faeba9f5b7cb639b863f2c30ec415
}

func ExampleHash3() {
	a, _ := new(big.Int).SetString("a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90", 16)
	b, _ := new(big.Int).SetString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", 16)
	c, _ := new(big.Int).SetString("fedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321", 16)
	h := poseidon.Hash3(a, b, c)
	fmt.Printf("0x%s\n", h.Text(16))
	// Output:
	// 0xb5b1978a26f5990a05bd9ed83ce3a101ad5778ab09db7c3c0f2667b9ad2be13f
}
