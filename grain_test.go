package poseidon

import (
	"encoding/json"
	"os"
	"testing"
)

type constantsJSON struct {
	T              int        `json:"t"`
	RoundsFull     int        `json:"roundsFull"`
	RoundsPartial  int        `json:"roundsPartial"`
	RoundConstants [][]string `json:"roundConstants"`
	MDS            [][]string `json:"mds"`
}

func loadConstants(t *testing.T, filename string) *constantsJSON {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var c constantsJSON
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	return &c
}

func TestGrainConstantsT3(t *testing.T) {
	ref := loadConstants(t, "testdata/constants_t3.json")
	rc, mds := GenConstants(3, 8, 57)

	if len(rc) != len(ref.RoundConstants) {
		t.Fatalf("round constants length: got %d, want %d", len(rc), len(ref.RoundConstants))
	}
	for i, row := range ref.RoundConstants {
		for j, hexStr := range row {
			want := hexToBig(hexStr)
			if rc[i][j].Cmp(want) != 0 {
				t.Fatalf("rc[%d][%d]: got %s, want %s", i, j, rc[i][j].Text(16), want.Text(16))
			}
		}
	}

	for i, row := range ref.MDS {
		for j, hexStr := range row {
			want := hexToBig(hexStr)
			if mds[i][j].Cmp(want) != 0 {
				t.Fatalf("mds[%d][%d]: got %s, want %s", i, j, mds[i][j].Text(16), want.Text(16))
			}
		}
	}
}

func TestGrainConstantsT4(t *testing.T) {
	ref := loadConstants(t, "testdata/constants_t4.json")
	rc, mds := GenConstants(4, 8, 56)

	if len(rc) != len(ref.RoundConstants) {
		t.Fatalf("round constants length: got %d, want %d", len(rc), len(ref.RoundConstants))
	}
	for i, row := range ref.RoundConstants {
		for j, hexStr := range row {
			want := hexToBig(hexStr)
			if rc[i][j].Cmp(want) != 0 {
				t.Fatalf("rc[%d][%d]: got %s, want %s", i, j, rc[i][j].Text(16), want.Text(16))
			}
		}
	}

	for i, row := range ref.MDS {
		for j, hexStr := range row {
			want := hexToBig(hexStr)
			if mds[i][j].Cmp(want) != 0 {
				t.Fatalf("mds[%d][%d]: got %s, want %s", i, j, mds[i][j].Text(16), want.Text(16))
			}
		}
	}
}

func TestHardcodedMatchesGrain(t *testing.T) {
	// Verify that constants.go hardcoded values match Grain LFSR output.
	rc, mds := GenConstants(3, 8, 57)
	for i, row := range rc {
		for j, val := range row {
			hardcoded := parseHex(rcHexT3[i][j])
			if val.Cmp(hardcoded) != 0 {
				t.Fatalf("T3 rc[%d][%d] mismatch: grain=%s, hardcoded=%s",
					i, j, val.Text(16), hardcoded.Text(16))
			}
		}
	}
	for i, row := range mds {
		for j, val := range row {
			hardcoded := parseHex(mdsHexT3[i][j])
			if val.Cmp(hardcoded) != 0 {
				t.Fatalf("T3 mds[%d][%d] mismatch: grain=%s, hardcoded=%s",
					i, j, val.Text(16), hardcoded.Text(16))
			}
		}
	}

	rc4, mds4 := GenConstants(4, 8, 56)
	for i, row := range rc4 {
		for j, val := range row {
			hardcoded := parseHex(rcHexT4[i][j])
			if val.Cmp(hardcoded) != 0 {
				t.Fatalf("T4 rc[%d][%d] mismatch: grain=%s, hardcoded=%s",
					i, j, val.Text(16), hardcoded.Text(16))
			}
		}
	}
	for i, row := range mds4 {
		for j, val := range row {
			hardcoded := parseHex(mdsHexT4[i][j])
			if val.Cmp(hardcoded) != 0 {
				t.Fatalf("T4 mds[%d][%d] mismatch: grain=%s, hardcoded=%s",
					i, j, val.Text(16), hardcoded.Text(16))
			}
		}
	}
}
