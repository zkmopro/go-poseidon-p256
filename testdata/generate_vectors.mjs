/**
 * Generate test vectors for go-poseidon-p256 using @noble/curves.
 *
 * Outputs:
 *   constants_t3.json  - round constants + MDS for t=3, RF=8, RP=57
 *   constants_t4.json  - round constants + MDS for t=4, RF=8, RP=56
 *   hash_vectors.json  - Poseidon hash input/output pairs
 *   field_vectors.json - field arithmetic test cases
 */

import { writeFileSync } from "fs";
import { Field } from "@noble/curves/abstract/modular";
import { grainGenConstants, poseidon } from "@noble/curves/abstract/poseidon";

const ORDER =
  0xffffffff00000001000000000000000000000000ffffffffffffffffffffffffn;
const Fp = Field(ORDER);

const ROUNDS_FULL = 8;

function toHex(n) {
  return "0x" + n.toString(16);
}

// ─── Field vectors ───────────────────────────────────────────────────────────

function generateFieldVectors() {
  const a = 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefn;
  const b = 0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321n;

  const vectors = {
    order: toHex(ORDER),
    add: [{ a: toHex(a), b: toHex(b), result: toHex(Fp.add(a, b)) }],
    mul: [{ a: toHex(a), b: toHex(b), result: toHex(Fp.mul(a, b)) }],
    neg: [{ a: toHex(a), result: toHex(Fp.neg(a)) }],
    inv: [
      { a: toHex(a), result: toHex(Fp.inv(a)) },
      { a: "0x1", result: toHex(Fp.inv(1n)) },
      { a: "0x2", result: toHex(Fp.inv(2n)) },
    ],
    sbox5: [],
    create: [
      { input: toHex(ORDER), result: "0x0" },
      { input: toHex(ORDER + 1n), result: "0x1" },
      {
        input: toHex(ORDER - 1n),
        result: toHex(ORDER - 1n),
      },
      { input: "0x0", result: "0x0" },
    ],
    mulN: [{ a: toHex(a), b: toHex(b), result: toHex(a * b) }],
    sqrN: [{ a: toHex(a), result: toHex(a * a) }],
  };

  // sbox5: n^5 mod ORDER
  const sboxInputs = [1n, 2n, 3n, a, ORDER - 1n];
  for (const n of sboxInputs) {
    // sbox5(n) = Fp.mul(Fp.sqrN(Fp.sqrN(n)), n)
    const result = Fp.mul(Fp.sqrN(Fp.sqrN(n)), n);
    vectors.sbox5.push({ input: toHex(n), result: toHex(result) });
  }

  // InvertBatch
  const batchInputs = [a, b, 42n];
  const batchResults = Fp.invertBatch(batchInputs);
  vectors.invertBatch = {
    inputs: batchInputs.map(toHex),
    results: batchResults.map(toHex),
  };

  // More add/mul edge cases
  vectors.add.push(
    { a: "0x0", b: "0x0", result: "0x0" },
    { a: toHex(ORDER - 1n), b: "0x1", result: "0x0" },
    { a: toHex(ORDER - 1n), b: "0x2", result: "0x1" }
  );
  vectors.mul.push(
    { a: "0x0", b: toHex(a), result: "0x0" },
    { a: "0x1", b: toHex(a), result: toHex(a) },
    {
      a: toHex(ORDER - 1n),
      b: toHex(ORDER - 1n),
      result: toHex(Fp.mul(ORDER - 1n, ORDER - 1n)),
    }
  );

  writeFileSync("field_vectors.json", JSON.stringify(vectors, null, 2));
  console.log("Wrote field_vectors.json");
}

// ─── Constants ───────────────────────────────────────────────────────────────

function generateConstants(t, rp) {
  const { roundConstants, mds } = grainGenConstants({
    Fp,
    t,
    roundsFull: ROUNDS_FULL,
    roundsPartial: rp,
  });

  const data = {
    t,
    roundsFull: ROUNDS_FULL,
    roundsPartial: rp,
    roundConstants: roundConstants.map((row) => row.map(toHex)),
    mds: mds.map((row) => row.map(toHex)),
  };

  writeFileSync(`constants_t${t}.json`, JSON.stringify(data, null, 2));
  console.log(`Wrote constants_t${t}.json`);
  return { roundConstants, mds };
}

// ─── Hash vectors ────────────────────────────────────────────────────────────

function generateHashVectors() {
  const configs = [
    { t: 3, rp: 57 },
    { t: 4, rp: 56 },
  ];

  const hashFns = {};
  for (const { t, rp } of configs) {
    const { roundConstants, mds } = grainGenConstants({
      Fp,
      t,
      roundsFull: ROUNDS_FULL,
      roundsPartial: rp,
    });
    hashFns[t] = poseidon({
      Fp,
      t,
      roundsFull: ROUNDS_FULL,
      roundsPartial: rp,
      sboxPower: 5,
      roundConstants,
      mds,
    });
  }

  const hash2 = (a, b) => hashFns[3]([0n, a, b])[0];
  const hash3 = (a, b, c) => hashFns[4]([0n, a, b, c])[0];

  const vectors = {
    hash2: [],
    hash3: [],
  };

  // Hash2 test cases
  const hash2Cases = [
    [0n, 0n],
    [1n, 2n],
    [ORDER - 1n, ORDER - 2n],
    [
      0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefn,
      0xcafebabedeadcafecafebabedeadcafecafebabedeadcafecafebabedeadcafen,
    ],
    [42n, 0n],
    [0n, 42n],
  ];

  for (const [a, b] of hash2Cases) {
    vectors.hash2.push({
      inputs: [toHex(a), toHex(b)],
      result: toHex(hash2(a, b)),
    });
  }

  // Hash3 test cases
  const hash3Cases = [
    [0n, 0n, 0n],
    [1n, 2n, 3n],
    [ORDER - 1n, ORDER - 2n, ORDER - 3n],
    [42n, 43n, 44n],
  ];

  for (const [a, b, c] of hash3Cases) {
    vectors.hash3.push({
      inputs: [toHex(a), toHex(b), toHex(c)],
      result: toHex(hash3(a, b, c)),
    });
  }

  // Round-by-round trace for Hash2(1, 2) - for debugging
  {
    const t = 3;
    const rp = 57;
    const { roundConstants, mds } = grainGenConstants({
      Fp,
      t,
      roundsFull: ROUNDS_FULL,
      roundsPartial: rp,
    });

    let values = [0n, 1n, 2n].map((v) => Fp.create(v));
    const halfFull = ROUNDS_FULL / 2;
    const rounds = ROUNDS_FULL + rp;
    const trace = [];

    const sboxFn = (n) => Fp.mul(Fp.sqrN(Fp.sqrN(n)), n);

    for (let round = 0; round < rounds; round++) {
      const isFull =
        round < halfFull || round >= halfFull + rp;

      // Add round constants
      values = values.map((v, j) => Fp.add(v, roundConstants[round][j]));

      // S-box
      if (isFull) {
        values = values.map((v) => sboxFn(v));
      } else {
        values[0] = sboxFn(values[0]);
      }

      // MDS
      values = mds.map((row) =>
        row.reduce((acc, m, j) => Fp.add(acc, Fp.mulN(m, values[j])), Fp.ZERO)
      );

      trace.push(values.map(toHex));
    }

    vectors.roundTrace = {
      description: "Round-by-round state for Hash2(1, 2), t=3",
      inputs: ["0x0", "0x1", "0x2"],
      trace,
    };
  }

  // SMT integration test: small 3-level Merkle tree
  {
    // Leaves: L0=1, L1=2, L2=3, L3=4
    // Leaf hash: Hash3(key, value, 1) where key=index
    const leafHash = (key, value) => hash3(BigInt(key), value, 1n);
    const l0 = leafHash(0, 1n);
    const l1 = leafHash(1, 2n);
    const l2 = leafHash(2, 3n);
    const l3 = leafHash(3, 4n);

    // Internal nodes
    const n01 = hash2(l0, l1);
    const n23 = hash2(l2, l3);
    const root = hash2(n01, n23);

    vectors.smtTree = {
      description: "3-level SMT: leaves=[1,2,3,4], leaf=Hash3(key,val,1), node=Hash2",
      leaves: [toHex(l0), toHex(l1), toHex(l2), toHex(l3)],
      internalNodes: [toHex(n01), toHex(n23)],
      root: toHex(root),
    };
  }

  // Input normalization test
  {
    const normalResult = hash2(1n, 0n);
    const overflowResult = hash2(ORDER + 1n, ORDER); // ORDER+1 -> 1, ORDER -> 0
    vectors.normalization = {
      description: "Hash2(ORDER+1, ORDER) should equal Hash2(1, 0)",
      normalInputs: ["0x1", "0x0"],
      normalResult: toHex(normalResult),
      overflowInputs: [toHex(ORDER + 1n), toHex(ORDER)],
      overflowResult: toHex(overflowResult),
    };
  }

  writeFileSync("hash_vectors.json", JSON.stringify(vectors, null, 2));
  console.log("Wrote hash_vectors.json");
}

// ─── Main ────────────────────────────────────────────────────────────────────

generateFieldVectors();
generateConstants(3, 57);
generateConstants(4, 56);
generateHashVectors();
console.log("Done.");
