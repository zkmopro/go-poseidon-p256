# go-poseidon-p256

[![Go Reference](https://pkg.go.dev/badge/github.com/zkmopro/go-poseidon-p256.svg)](https://pkg.go.dev/github.com/zkmopro/go-poseidon-p256)

A pure-Go implementation of the [Poseidon hash function](https://eprint.iacr.org/2019/458) over the **P-256 (secp256r1) scalar field**.

Poseidon is an arithmetic-friendly hash function designed for zero-knowledge proof systems. This implementation targets the P-256 scalar field, making it suitable for ZK applications that operate over the NIST P-256 curve (e.g., proving ECDSA signature validity).

## Installation

```bash
go get github.com/zkmopro/go-poseidon-p256
```

## Usage

```go
package main

import (
	"fmt"
	"math/big"

	poseidon "github.com/zkmopro/go-poseidon-p256"
)

func main() {
	a := big.NewInt(1)
	b := big.NewInt(2)

	// Hash two field elements
	h := poseidon.Hash2(a, b)
	fmt.Printf("Hash2: 0x%s\n", h.Text(16))

	// Hash three field elements
	c := big.NewInt(3)
	h3 := poseidon.Hash3(a, b, c)
	fmt.Printf("Hash3: 0x%s\n", h3.Text(16))
}
```

## API

| Function | Description |
|---|---|
| `Hash2(a, b *big.Int) *big.Int` | Hash two field elements (t=3 Poseidon) |
| `Hash3(a, b, c *big.Int) *big.Int` | Hash three field elements (t=4 Poseidon) |
| `Hash(inputs []*big.Int) *big.Int` | Generic hash for 2 or 3 inputs |
| `GenConstants(t, roundsFull, roundsPartial int)` | Generate Poseidon round constants and MDS matrix |

Inputs are automatically reduced modulo the P-256 scalar field order. Results are in `[0, ORDER)`.

## Benchmarks

Run benchmarks with:

```bash
go test -bench=. -benchmem ./...
```

Results on Apple M3 chip using random 256-bit field elements:

| Function | ns/op | B/op | allocs/op |
|---|---|---|---|
| Hash2 | ~212,000 | 271,697 | 3,494 |
| Hash3 | ~342,000 | 426,455 | 5,495 |

## Security

This implementation uses Go's `math/big` for field arithmetic, which is **not constant-time**. It is suitable for hashing public data in ZK circuits but should **not** be used for operations involving secret keys or other sensitive values where timing side-channels are a concern.

## References

- [Poseidon: A New Hash Function for Zero-Knowledge Proof Systems](https://eprint.iacr.org/2019/458) (Grassi et al., 2019)
- [iden3/go-iden3-crypto Poseidon](https://github.com/iden3/go-iden3-crypto/tree/master/poseidon) — Go Poseidon implementation over the BN254 scalar field
