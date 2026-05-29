/**
 * Veritrail — public library API.
 *
 * A portable, dependency-free (node:crypto only) toolkit for the Veritrail
 * byte-source layer: JCS canonicalization, multihash/multibase hash-strings,
 * receipt content-addressing, JWS signing + verification (Ed25519 / ES256),
 * full lineage-chain verification, and the MCP / A2A co-signing middleware.
 *
 * The MCP and A2A *transport adapters* (server.ts / demo.ts) live outside this
 * entry point because they depend on the respective SDKs (devDependencies). The
 * co-signing *middleware* exported here is SDK-free and safe for library consumers.
 */

// --- JCS canonicalization (RFC 8785) ---
export { jcs, jcsBytes, jcsString } from './commands/jcs.js';
export type { JcsInput, JcsOutput } from './commands/jcs.js';

// --- Hash-string (multihash + multibase) ---
export { hashstring, encodeHashstring } from './commands/hashstring.js';
export type { HashstringInput, HashstringOutput } from './commands/hashstring.js';

// --- Digest (sha256 of raw bytes -> hash-string) ---
export { digest, digestBytes } from './commands/digest.js';
export type { DigestInput, DigestOutput } from './commands/digest.js';

// --- Receipt content-addressing ---
export { receiptId } from './commands/receipt-id.js';
export type { ReceiptIdInput, ReceiptIdOutput } from './commands/receipt-id.js';

// --- Cost-block canonicalization ---
export { costCanon } from './commands/cost-canon.js';
export type { CostBlock, CostInput, CostOutput } from './commands/cost-canon.js';

// --- SSE decode-then-hash ---
export { sseOutputsHash } from './commands/sse-outputs-hash.js';
export type { SseInput, SseOutput } from './commands/sse-outputs-hash.js';

// --- JWS verification + JOSE hardening (§7) ---
export { verify } from './commands/verify.js';
export type { VerifyInput, VerifyOutput, RegistryJwk } from './commands/verify.js';

// --- Lineage-chain verification (§8) ---
export { verifyChain } from './commands/verify-chain.js';
export type { VerifyChainInput, VerifyChainOutput } from './commands/verify-chain.js';

// --- A2A artifact canonicalization (§9) ---
export { a2aArtifactHash, buildDescriptor } from './commands/a2a-artifact-hash.js';
export type { A2AArtifactHashInput, A2AArtifactHashOutput } from './commands/a2a-artifact-hash.js';

// --- Signing + receipt construction (SDK-free core) ---
export {
  signReceipt,
  signReceipt as sign,
  generateTestSigner,
  registryFromSigner,
  buildReceiptPayload,
  hashOf,
  stripMeta,
  computeOutputsHash,
  RECEIPT_META_KEY,
} from './mcp/veritrail-middleware.js';
export type { TestSigner, McpResultLike, BuildReceiptOptions } from './mcp/veritrail-middleware.js';

// --- MCP co-signing middleware (SDK-free) ---
export { veritrailToolResult } from './mcp/veritrail-middleware.js';

// --- A2A co-signing middleware (SDK-free) ---
export {
  veritrailArtifact,
  computeArtifactHash,
  buildA2AReceiptPayload,
} from './a2a/veritrail-a2a-middleware.js';
export type {
  A2APartLike,
  A2AArtifactLike,
  BuildA2AReceiptOptions,
} from './a2a/veritrail-a2a-middleware.js';
