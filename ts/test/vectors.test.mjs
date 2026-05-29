/**
 * Vector-driven conformance tests.
 * Each test loads a vector JSON, invokes the implementation function,
 * and asserts the result matches the anchor byte-for-byte (as JSON string).
 *
 * STRICT TDD: these tests are written BEFORE the implementation exists.
 * They will fail until ts/dist/ is built with the correct logic.
 */
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const VECTORS_DIR = join(__dirname, '../../conformance/vectors');
const DIST = join(__dirname, '../dist');

// Lazy-load built modules (fails loudly if build hasn't run)
const { jcs } = await import(`${DIST}/commands/jcs.js`);
const { hashstring } = await import(`${DIST}/commands/hashstring.js`);
const { digest } = await import(`${DIST}/commands/digest.js`);
const { receiptId } = await import(`${DIST}/commands/receipt-id.js`);
const { sseOutputsHash } = await import(`${DIST}/commands/sse-outputs-hash.js`);
const { costCanon } = await import(`${DIST}/commands/cost-canon.js`);
const { verify } = await import(`${DIST}/commands/verify.js`);
const { verifyChain } = await import(`${DIST}/commands/verify-chain.js`);
const { a2aArtifactHash } = await import(`${DIST}/commands/a2a-artifact-hash.js`);

function loadVector(filename) {
  return JSON.parse(readFileSync(join(VECTORS_DIR, filename), 'utf8'));
}

function runCommand(cmd, input) {
  switch (cmd) {
    case 'jcs': return jcs(input);
    case 'hashstring': return hashstring(input);
    case 'digest': return digest(input);
    case 'receipt-id': return receiptId(input);
    case 'sse-outputs-hash': return sseOutputsHash(input);
    case 'cost-canon': return costCanon(input);
    case 'verify': return verify(input);
    case 'verify-chain': return verifyChain(input);
    case 'a2a-artifact-hash': return a2aArtifactHash(input);
    default: return { error: 'unknown_command' };
  }
}

// Helper: sort keys recursively for JCS output comparison
// We compare each field individually to get good error messages
function assertAnchorMatch(result, anchor, vectorName) {
  for (const [key, expected] of Object.entries(anchor)) {
    assert.deepStrictEqual(
      result[key],
      expected,
      `Vector "${vectorName}": field "${key}" mismatch.\n  got:      ${JSON.stringify(result[key])}\n  expected: ${JSON.stringify(expected)}`
    );
  }
}

// --- jcs/sorted-keys ---
test('jcs/sorted-keys', () => {
  const v = loadVector('jcs-sorted-keys.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- jcs/nested-and-unicode ---
test('jcs/nested-and-unicode', () => {
  const v = loadVector('jcs-nested.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- jcs/max-safe-int ---
test('jcs/max-safe-int', () => {
  const v = loadVector('jcs-max-safe-int.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- jcs/over-safe-int (no anchor: just assert it produces a result, no throw) ---
test('jcs/over-safe-int-no-anchor', () => {
  const v = loadVector('jcs-over-safe-int.json');
  const result = runCommand(v.command, v.input);
  assert.equal(typeof result.canonical_hex, 'string');
  assert.equal(typeof result.byte_len, 'number');
});

// --- digest/empty ---
test('digest/empty', () => {
  const v = loadVector('digest-empty.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- digest/hello ---
test('digest/hello', () => {
  const v = loadVector('digest-hello.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- receipt-id/basic ---
test('receipt-id/basic', () => {
  const v = loadVector('receipt-id-basic.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- receipt-id: error when receipt_id key present ---
test('receipt-id/error-when-receipt-id-present', () => {
  const result = receiptId({ receipt: { receipt_id: 'u123', v: '0.1' } });
  assert.deepStrictEqual(result, { error: 'receipt_id_must_be_absent' });
});

// --- sse/bom-no-space ---
test('sse/bom-no-space', () => {
  const v = loadVector('sse-bom-no-space.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- sse/two-data-one-event ---
test('sse/two-data-one-event', () => {
  const v = loadVector('sse-two-data-one-event.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- sse/two-events ---
test('sse/two-events', () => {
  const v = loadVector('sse-two-events.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- sse/reframed-crlf-comment-event ---
test('sse/reframed-crlf-comment-event', () => {
  const v = loadVector('sse-reframed-equivalent.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- sse/jsonrpc-inner-result ---
test('sse/jsonrpc-inner-result', () => {
  const v = loadVector('sse-jsonrpc-inner-result.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- sse/jsonrpc-null-result ---
test('sse/jsonrpc-null-result', () => {
  const v = loadVector('sse-jsonrpc-null-result.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- sse/zero-event ---
test('sse/zero-event', () => {
  const v = loadVector('sse-zero-event.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- cost/valid-string-ints ---
test('cost/valid-string-ints', () => {
  const v = loadVector('cost-valid.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- cost/number-must-error ---
test('cost/number-must-error', () => {
  const v = loadVector('cost-number-error.json');
  const result = runCommand(v.command, v.input);
  assertAnchorMatch(result, v.anchor, v.name);
});

// --- verify: all 12 vectors, table-driven ---
const VERIFY_VECTORS = [
  'verify-valid-ed25519.json',
  'verify-valid-es256.json',
  'verify-alg-none.json',
  'verify-alg-substitution.json',
  'verify-jwk-header-injection.json',
  'verify-bad-signature.json',
  'verify-malformed-payload.json',
  'verify-non-canonical-payload.json',
  'verify-unknown-kid.json',
  'verify-revoked-key.json',
  'verify-context-mismatch.json',
  'verify-wrong-key.json',
];

for (const filename of VERIFY_VECTORS) {
  const v = loadVector(filename);
  test(v.name, () => {
    const result = runCommand(v.command, v.input);
    assertAnchorMatch(result, v.anchor, v.name);
  });
}

// --- verify-chain: all 7 vectors, table-driven ---
const CHAIN_VECTORS = [
  'chain-valid-2hop.json',
  'chain-valid-3hop.json',
  'chain-link-mismatch.json',
  'chain-forged-parent.json',
  'chain-dangling-parent.json',
  'chain-parent-identity-mismatch.json',
  'chain-middle-null-parent.json',
];

for (const filename of CHAIN_VECTORS) {
  const v = loadVector(filename);
  test(v.name, () => {
    const result = runCommand(v.command, v.input);
    assertAnchorMatch(result, v.anchor, v.name);
  });
}

// --- a2a-artifact-hash: all 6 vectors, table-driven ---
const A2A_VECTORS = [
  'a2a-text.json',
  'a2a-data.json',
  'a2a-file-inline.json',
  'a2a-file-uri.json',
  'a2a-mixed.json',
  'a2a-unsupported.json',
];

for (const filename of A2A_VECTORS) {
  const v = loadVector(filename);
  test(v.name, () => {
    const result = runCommand(v.command, v.input);
    assertAnchorMatch(result, v.anchor, v.name);
  });
}
