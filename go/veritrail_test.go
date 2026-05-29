package veritrail_test

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Daily-Nerd/veritrail/go"
)

// vectorFile represents a test vector loaded from disk.
type vectorFile struct {
	Name    string          `json:"name"`
	Command string          `json:"command"`
	Input   json.RawMessage `json:"input"`
	Anchor  json.RawMessage `json:"anchor"`
}

// loadVectors loads all *.json files from the vectors directory relative to the module root.
func loadVectors(t *testing.T) []vectorFile {
	t.Helper()
	// vectors live under conformance/vectors at the repo root;
	// this package dir is the module root (go/), so go up one then into conformance/vectors.
	dir := filepath.Join("..", "conformance", "vectors")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("cannot read vectors dir %s: %v", dir, err)
	}
	var vecs []vectorFile
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		var v vectorFile
		if err := json.Unmarshal(data, &v); err != nil {
			t.Fatalf("parse %s: %v", e.Name(), err)
		}
		vecs = append(vecs, v)
	}
	return vecs
}

// -------------------------------------------------------------------
// Unit tests for JCS
// -------------------------------------------------------------------

func TestJCS_SortedKeys(t *testing.T) {
	input := json.RawMessage(`{"b":1,"a":2}`)
	got, err := veritrail.JCS(input)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"a":2,"b":1}`
	if string(got) != want {
		t.Errorf("JCS = %q, want %q", got, want)
	}
}

func TestJCS_NestedAndUnicode(t *testing.T) {
	// "é" is U+00E9, encoded in UTF-8 as 0xC3 0xA9
	input := json.RawMessage(`{"z":"é","a":[3,1,2]}`)
	got, err := veritrail.JCS(input)
	if err != nil {
		t.Fatal(err)
	}
	wantHex := "7b2261223a5b332c312c325d2c227a223a22c3a9227d"
	if hex.EncodeToString(got) != wantHex {
		t.Errorf("JCS hex = %s, want %s", hex.EncodeToString(got), wantHex)
	}
}

// -------------------------------------------------------------------
// Unit tests for HashString
// -------------------------------------------------------------------

func TestHashString_KnownDigest(t *testing.T) {
	// sha256 of "hello" is well-known
	digestHex := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	got, err := veritrail.HashString(digestHex)
	if err != nil {
		t.Fatal(err)
	}
	// Verify it starts with "u"
	if len(got) == 0 || got[0] != 'u' {
		t.Errorf("HashString = %q, want leading 'u'", got)
	}
	// Verify decodable base64url
	b, err := base64.RawURLEncoding.DecodeString(got[1:])
	if err != nil {
		t.Fatalf("base64url decode failed: %v", err)
	}
	// First two bytes must be 0x12, 0x20 (sha2-256 varint code, varint len=32)
	if len(b) < 2 || b[0] != 0x12 || b[1] != 0x20 {
		t.Errorf("multihash prefix = %x, want 1220", b[:2])
	}
}

// -------------------------------------------------------------------
// Unit tests for Digest
// -------------------------------------------------------------------

func TestDigest_Empty(t *testing.T) {
	got, err := veritrail.Digest([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	want := "uEiDjsMRCmPwcFJr79MiZb7kkJ65B5GSbk0yklZkbeFK4VQ"
	if got != want {
		t.Errorf("Digest(empty) = %q, want %q", got, want)
	}
}

func TestDigest_Hello(t *testing.T) {
	got, err := veritrail.Digest([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	want := "uEiAs8k26X7CjDiboOyrFueKeGxYeXB-nQl5zBDNik4uYJA"
	if got != want {
		t.Errorf("Digest(hello) = %q, want %q", got, want)
	}
}

// -------------------------------------------------------------------
// Unit tests for SSE decode
// -------------------------------------------------------------------

func TestSSEDecode_TwoDataOneEvent(t *testing.T) {
	// data: hello\ndata: world\n\n  -> "hello\nworld"
	raw := []byte("data: hello\ndata: world\n\n")
	msgs, err := veritrail.SSEDecode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 {
		t.Fatalf("SSEDecode: got %d messages, want 1", len(msgs))
	}
	if msgs[0] != "hello\nworld" {
		t.Errorf("SSEDecode: msg = %q, want %q", msgs[0], "hello\nworld")
	}
}

func TestSSEDecode_TwoEvents(t *testing.T) {
	raw := []byte("data: a\n\ndata: b\n\n")
	msgs, err := veritrail.SSEDecode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 {
		t.Fatalf("SSEDecode: got %d messages, want 2", len(msgs))
	}
	if msgs[0] != "a" || msgs[1] != "b" {
		t.Errorf("SSEDecode: msgs = %v, want [a b]", msgs)
	}
}

func TestSSEDecode_BOMNoSpace(t *testing.T) {
	// BOM + data:x (no space) + \n\n
	raw := []byte("\xef\xbb\xbfdata:x\n\n")
	msgs, err := veritrail.SSEDecode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 || msgs[0] != "x" {
		t.Errorf("SSEDecode BOM: msgs = %v, want [x]", msgs)
	}
}

func TestSSEDecode_CRLFCommentEvent(t *testing.T) {
	// event: msg\r\n: comment\r\nid: 7\r\ndata: hello\r\ndata: world\r\n\r\n
	raw := []byte("event: msg\r\n: a comment line\r\nid: 7\r\ndata: hello\r\ndata: world\r\n\r\n")
	msgs, err := veritrail.SSEDecode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 || msgs[0] != "hello\nworld" {
		t.Errorf("SSEDecode CRLF: msgs = %v, want [hello\\nworld]", msgs)
	}
}

func TestSSEDecode_EmptyDataNotDispatched(t *testing.T) {
	// event with no data lines -> empty data -> not dispatched
	raw := []byte("event: ping\n\n")
	msgs, err := veritrail.SSEDecode(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 0 {
		t.Errorf("SSEDecode empty: got %d messages, want 0", len(msgs))
	}
}

// -------------------------------------------------------------------
// Unit tests for CostCanon
// -------------------------------------------------------------------

func TestCostCanon_ValidStringInts(t *testing.T) {
	input := json.RawMessage(`{"tokens":"1500","usd_micros":"10000000000","wall_ms":"845","rail_ref":null}`)
	got, err := veritrail.CostCanon(input)
	if err != nil {
		t.Fatal(err)
	}
	wantHex := "7b227261696c5f726566223a6e756c6c2c22746f6b656e73223a2231353030222c227573645f6d6963726f73223a223130303030303030303030222c2277616c6c5f6d73223a22383435227d"
	if hex.EncodeToString(got) != wantHex {
		t.Errorf("CostCanon hex = %s, want %s", hex.EncodeToString(got), wantHex)
	}
}

func TestCostCanon_NumberError(t *testing.T) {
	input := json.RawMessage(`{"tokens":1500,"usd_micros":"10","wall_ms":"5","rail_ref":null}`)
	_, err := veritrail.CostCanon(input)
	if err == nil {
		t.Fatal("expected error for numeric tokens, got nil")
	}
}

// -------------------------------------------------------------------
// Vector-driven tests (load all vectors, run against implementations)
// -------------------------------------------------------------------

func TestVectors_JCS(t *testing.T) {
	for _, v := range loadVectors(t) {
		if v.Command != "jcs" {
			continue
		}
		v := v
		t.Run(v.Name, func(t *testing.T) {
			var inp struct {
				Value json.RawMessage `json:"value"`
			}
			if err := json.Unmarshal(v.Input, &inp); err != nil {
				t.Fatal(err)
			}
			got, err := veritrail.JCS(inp.Value)
			if err != nil {
				t.Fatal(err)
			}
			// Anchor is optional: when absent, the harness only checks
			// cross-impl agreement. Here we just assert JCS did not error.
			if len(v.Anchor) == 0 || string(v.Anchor) == "null" {
				return
			}
			var anchor struct {
				CanonicalHex string `json:"canonical_hex"`
				ByteLen      int    `json:"byte_len"`
			}
			if err := json.Unmarshal(v.Anchor, &anchor); err != nil {
				t.Fatal(err)
			}
			if hex.EncodeToString(got) != anchor.CanonicalHex {
				t.Errorf("canonical_hex = %s, want %s", hex.EncodeToString(got), anchor.CanonicalHex)
			}
			if len(got) != anchor.ByteLen {
				t.Errorf("byte_len = %d, want %d", len(got), anchor.ByteLen)
			}
		})
	}
}

func TestVectors_Digest(t *testing.T) {
	for _, v := range loadVectors(t) {
		if v.Command != "digest" {
			continue
		}
		v := v
		t.Run(v.Name, func(t *testing.T) {
			var inp struct {
				BytesB64 string `json:"bytes_b64"`
			}
			if err := json.Unmarshal(v.Input, &inp); err != nil {
				t.Fatal(err)
			}
			raw, err := base64.StdEncoding.DecodeString(inp.BytesB64)
			if err != nil {
				// Try raw (no padding)
				raw, err = base64.RawStdEncoding.DecodeString(inp.BytesB64)
				if err != nil {
					t.Fatalf("base64 decode: %v", err)
				}
			}
			got, err := veritrail.Digest(raw)
			if err != nil {
				t.Fatal(err)
			}
			var anchor struct {
				Hashstr string `json:"hashstr"`
			}
			if err := json.Unmarshal(v.Anchor, &anchor); err != nil {
				t.Fatal(err)
			}
			if got != anchor.Hashstr {
				t.Errorf("hashstr = %q, want %q", got, anchor.Hashstr)
			}
		})
	}
}

func TestVectors_ReceiptID(t *testing.T) {
	for _, v := range loadVectors(t) {
		if v.Command != "receipt-id" {
			continue
		}
		v := v
		t.Run(v.Name, func(t *testing.T) {
			var inp struct {
				Receipt json.RawMessage `json:"receipt"`
			}
			if err := json.Unmarshal(v.Input, &inp); err != nil {
				t.Fatal(err)
			}
			canonHex, receiptID, err := veritrail.ReceiptID(inp.Receipt)
			if err != nil {
				t.Fatalf("ReceiptID error: %v", err)
			}
			var anchor struct {
				CanonicalHex string `json:"canonical_hex"`
				ReceiptID    string `json:"receipt_id"`
			}
			if err := json.Unmarshal(v.Anchor, &anchor); err != nil {
				t.Fatal(err)
			}
			if canonHex != anchor.CanonicalHex {
				t.Errorf("canonical_hex = %s, want %s", canonHex, anchor.CanonicalHex)
			}
			if receiptID != anchor.ReceiptID {
				t.Errorf("receipt_id = %s, want %s", receiptID, anchor.ReceiptID)
			}
		})
	}
}

func TestVectors_SSEOutputsHash(t *testing.T) {
	for _, v := range loadVectors(t) {
		if v.Command != "sse-outputs-hash" {
			continue
		}
		v := v
		t.Run(v.Name, func(t *testing.T) {
			var inp struct {
				RawB64 string `json:"raw_b64"`
				Mode   string `json:"mode"`
			}
			if err := json.Unmarshal(v.Input, &inp); err != nil {
				t.Fatal(err)
			}
			rawBytes, err := base64.StdEncoding.DecodeString(inp.RawB64)
			if err != nil {
				rawBytes, err = base64.RawStdEncoding.DecodeString(inp.RawB64)
				if err != nil {
					t.Fatalf("base64 decode: %v", err)
				}
			}
			decodedHex, outputsHash, err := veritrail.SSEOutputsHash(rawBytes, inp.Mode)
			if err != nil {
				t.Fatalf("SSEOutputsHash error: %v", err)
			}
			var anchor struct {
				DecodedHex  string `json:"decoded_hex"`
				OutputsHash string `json:"outputs_hash"`
			}
			if err := json.Unmarshal(v.Anchor, &anchor); err != nil {
				t.Fatal(err)
			}
			if decodedHex != anchor.DecodedHex {
				t.Errorf("decoded_hex = %s, want %s", decodedHex, anchor.DecodedHex)
			}
			if outputsHash != anchor.OutputsHash {
				t.Errorf("outputs_hash = %s, want %s", outputsHash, anchor.OutputsHash)
			}
		})
	}
}

func TestVectors_CostCanon(t *testing.T) {
	for _, v := range loadVectors(t) {
		if v.Command != "cost-canon" {
			continue
		}
		v := v
		t.Run(v.Name, func(t *testing.T) {
			var inp struct {
				Cost json.RawMessage `json:"cost"`
			}
			if err := json.Unmarshal(v.Input, &inp); err != nil {
				t.Fatal(err)
			}

			var anchor map[string]json.RawMessage
			if err := json.Unmarshal(v.Anchor, &anchor); err != nil {
				t.Fatal(err)
			}

			if _, hasErr := anchor["error"]; hasErr {
				// Expect error
				_, err := veritrail.CostCanon(inp.Cost)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				got, err := veritrail.CostCanon(inp.Cost)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				var expectedAnchor struct {
					CanonicalHex string `json:"canonical_hex"`
				}
				if err := json.Unmarshal(v.Anchor, &expectedAnchor); err != nil {
					t.Fatal(err)
				}
				if hex.EncodeToString(got) != expectedAnchor.CanonicalHex {
					t.Errorf("canonical_hex = %s, want %s", hex.EncodeToString(got), expectedAnchor.CanonicalHex)
				}
			}
		})
	}
}
