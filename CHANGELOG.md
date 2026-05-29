# Changelog

All notable changes to veritrail are documented here. Format: [Keep a Changelog](https://keepachangelog.com/); versioning: [SemVer](https://semver.org/).

## [0.1.0] — unreleased

Initial reference release. Protocol version `veritrail/0.1`.

### Added
- **Receipt protocol** — performer-attested, content-addressed execution receipts for AI agent actions (`docs/DESIGN.md`).
- **Two reference implementations**, byte-identical across runtimes:
  - Go (`github.com/Daily-Nerd/veritrail/go`) — `Sign`, `Verify`, `VerifyChain`, `A2AArtifactHash`, `JCS`, `HashString`, `Digest`, `ReceiptID`, `SSEOutputsHash`, `CostCanon` + `veritrail-verify` CLI.
  - TypeScript (`veritrail`) — `sign`, `verify`, `verifyChain`, `a2aArtifactHash`, JCS/hash-string helpers, receipt types, and MCP + A2A co-signing middleware. Zero runtime dependencies.
- **RFC 8785 JCS** canonicalization with a pinned multibase hash-string encoding.
- **Signing & hardened verification** — Ed25519 + ES256 (JWS); rejects `alg:none`, algorithm substitution, `jwk`/`jku`/`x5*` header key-injection, non-canonical payloads, unknown/revoked keys.
- **Chain / DAG verification** with `parent_performer_id` foreign-splice defense.
- **MCP binding** (`result._meta["dev.veritrail/receipt"]`) and **A2A binding** (artifact metadata + artifact-canonicalization).
- **SSE** decode-then-hash streaming commitment.
- **Conformance suite** — 41 language-agnostic vectors + reproducible generators + cross-implementation harness; enforced in CI (Go ≡ TS, byte-identical).

### Notes
- Honest value cap: veritrail attests *what a performer did and returned* — performer non-repudiation, byte integrity, authorization-binding-by-hash, verifiable cost. It does **not** prove correctness, world side-effects, or intent integrity.
- APIs may change before `1.0`.
