// veritrail-verify is the Veritrail conformance verifier CLI.
// Usage: veritrail-verify <command> < input.json
// Output: a single line of JCS-canonical JSON to stdout, exit 0.
// On error: {"error":"<code>"} to stdout, exit 0.
package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Daily-Nerd/veritrail/go"
	"github.com/gowebpki/jcs"
)

func main() {
	if len(os.Args) < 2 {
		writeError("missing_command")
		return
	}
	command := os.Args[1]

	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		writeError("stdin_read_error")
		return
	}

	var result json.RawMessage
	var cmdErr error

	switch command {
	case "jcs":
		result, cmdErr = runJCS(stdin)
	case "hashstring":
		result, cmdErr = runHashString(stdin)
	case "digest":
		result, cmdErr = runDigest(stdin)
	case "receipt-id":
		result, cmdErr = runReceiptID(stdin)
	case "sse-outputs-hash":
		result, cmdErr = runSSEOutputsHash(stdin)
	case "cost-canon":
		result, cmdErr = runCostCanon(stdin)
	case "verify":
		result, cmdErr = runVerify(stdin)
	case "verify-chain":
		result, cmdErr = runVerifyChain(stdin)
	case "a2a-artifact-hash":
		result, cmdErr = runA2AArtifactHash(stdin)
	default:
		writeError("unsupported_command")
		return
	}

	if cmdErr != nil {
		writeError(errorCode(cmdErr))
		return
	}

	// The result must itself be JCS-canonical. Since we build it via
	// json.Marshal on structs with explicit field names, and Go's
	// encoding/json emits object keys in declaration order (not sorted),
	// we must pass it through jcs.Transform before printing.
	canonical, err := jcs.Transform(result)
	if err != nil {
		writeError("jcs_output_error")
		return
	}
	fmt.Fprintf(os.Stdout, "%s\n", canonical)
}

func writeError(code string) {
	// {"error":"<code>"} — keys already sorted (only one key), so safe to
	// build manually and put through JCS for byte-exactness.
	out := map[string]string{"error": code}
	b, _ := json.Marshal(out)
	canonical, _ := jcs.Transform(b)
	fmt.Fprintf(os.Stdout, "%s\n", canonical)
}

func errorCode(err error) string {
	if errors.Is(err, veritrail.ErrCostMustBeStringInt) {
		return "cost_must_be_string_int"
	}
	if errors.Is(err, veritrail.ErrReceiptIDMustBeAbsent) {
		return "receipt_id_must_be_absent"
	}
	if errors.Is(err, veritrail.ErrUnsupportedPart) {
		return "unsupported_part"
	}
	return "invalid_input"
}

// -------------------------------------------------------------------
// Command implementations
// -------------------------------------------------------------------

func runJCS(stdin []byte) (json.RawMessage, error) {
	var inp struct {
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	canon, err := veritrail.JCS(inp.Value)
	if err != nil {
		return nil, err
	}
	out := struct {
		CanonicalHex string `json:"canonical_hex"`
		ByteLen      int    `json:"byte_len"`
	}{
		CanonicalHex: hex.EncodeToString(canon),
		ByteLen:      len(canon),
	}
	return json.Marshal(out)
}

func runHashString(stdin []byte) (json.RawMessage, error) {
	var inp struct {
		Algo      string `json:"algo"`
		DigestHex string `json:"digest_hex"`
	}
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	if inp.Algo != "sha2-256" {
		return nil, errors.New("unsupported algo: " + inp.Algo)
	}
	hashStr, err := veritrail.HashString(inp.DigestHex)
	if err != nil {
		return nil, err
	}
	out := struct {
		Hashstr string `json:"hashstr"`
	}{Hashstr: hashStr}
	return json.Marshal(out)
}

func runDigest(stdin []byte) (json.RawMessage, error) {
	var inp struct {
		BytesB64 string `json:"bytes_b64"`
	}
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	raw, err := base64.StdEncoding.DecodeString(inp.BytesB64)
	if err != nil {
		// Try without padding
		raw, err = base64.RawStdEncoding.DecodeString(inp.BytesB64)
		if err != nil {
			return nil, err
		}
	}
	hashStr, err := veritrail.Digest(raw)
	if err != nil {
		return nil, err
	}
	out := struct {
		Hashstr string `json:"hashstr"`
	}{Hashstr: hashStr}
	return json.Marshal(out)
}

func runReceiptID(stdin []byte) (json.RawMessage, error) {
	var inp struct {
		Receipt json.RawMessage `json:"receipt"`
	}
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	canonHex, receiptID, err := veritrail.ReceiptID(inp.Receipt)
	if err != nil {
		return nil, err
	}
	out := struct {
		CanonicalHex string `json:"canonical_hex"`
		ReceiptID    string `json:"receipt_id"`
	}{
		CanonicalHex: canonHex,
		ReceiptID:    receiptID,
	}
	return json.Marshal(out)
}

func runSSEOutputsHash(stdin []byte) (json.RawMessage, error) {
	var inp struct {
		RawB64 string `json:"raw_b64"`
		Mode   string `json:"mode"`
	}
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	rawBytes, err := base64.StdEncoding.DecodeString(inp.RawB64)
	if err != nil {
		rawBytes, err = base64.RawStdEncoding.DecodeString(inp.RawB64)
		if err != nil {
			return nil, err
		}
	}
	decodedHex, outputsHash, err := veritrail.SSEOutputsHash(rawBytes, inp.Mode)
	if err != nil {
		return nil, err
	}
	out := struct {
		DecodedHex  string `json:"decoded_hex"`
		OutputsHash string `json:"outputs_hash"`
	}{
		DecodedHex:  decodedHex,
		OutputsHash: outputsHash,
	}
	return json.Marshal(out)
}

func runCostCanon(stdin []byte) (json.RawMessage, error) {
	var inp struct {
		Cost json.RawMessage `json:"cost"`
	}
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	canon, err := veritrail.CostCanon(inp.Cost)
	if err != nil {
		return nil, err
	}
	out := struct {
		CanonicalHex string `json:"canonical_hex"`
	}{CanonicalHex: hex.EncodeToString(canon)}
	return json.Marshal(out)
}

func runVerify(stdin []byte) (json.RawMessage, error) {
	var inp veritrail.VerifyInput
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	valid, reason := veritrail.Verify(inp)
	out := struct {
		Valid  bool   `json:"valid"`
		Reason string `json:"reason"`
	}{
		Valid:  valid,
		Reason: reason,
	}
	return json.Marshal(out)
}

func runVerifyChain(stdin []byte) (json.RawMessage, error) {
	var inp veritrail.VerifyChainInput
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	valid, reason, chainLen := veritrail.VerifyChain(inp)
	out := struct {
		Valid    bool   `json:"valid"`
		Reason   string `json:"reason"`
		ChainLen int    `json:"chain_len"`
	}{
		Valid:    valid,
		Reason:   reason,
		ChainLen: chainLen,
	}
	return json.Marshal(out)
}

func runA2AArtifactHash(stdin []byte) (json.RawMessage, error) {
	var inp struct {
		Artifact json.RawMessage `json:"artifact"`
	}
	if err := json.Unmarshal(stdin, &inp); err != nil {
		return nil, err
	}
	descriptorBytes, outputsHash, err := veritrail.A2AArtifactHash(inp.Artifact)
	if err != nil {
		return nil, err
	}
	out := struct {
		OutputsHash   string `json:"outputs_hash"`
		DescriptorHex string `json:"descriptor_hex"`
	}{
		OutputsHash:   outputsHash,
		DescriptorHex: hex.EncodeToString(descriptorBytes),
	}
	return json.Marshal(out)
}
