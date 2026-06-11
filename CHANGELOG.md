# Changelog

All notable changes to veritrail are documented here. Format: [Keep a Changelog](https://keepachangelog.com/); versioning: [SemVer](https://semver.org/).

> **Historical record — frozen at 0.1.1.** From the next release onward, changelogs are generated per implementation by release-please: **`ts/CHANGELOG.md`** (npm `@daily-nerd/veritrail`) and **`go/CHANGELOG.md`** (Go module). Versions stay linked, so both bump together.

## [0.1.1] — 2026-05-29

### Changed
- **npm package name:** the TypeScript implementation publishes as **`@daily-nerd/veritrail`**. The bare `veritrail` name is taken on npm by an unrelated package, so the scoped name is final. The `veritrail-verify` CLI binary name is unchanged.
- **Go module layout:** the library package moved to the module root. The import path is now `github.com/Daily-Nerd/veritrail/go` (was `…/go/veritrail`) — use `veritrail.Sign(...)`, `veritrail.Verify(...)`. Install the CLI with `go install github.com/Daily-Nerd/veritrail/go/cmd/veritrail-verify@v0.1.1`. (v0.1.0's nested import path is superseded; behavior and conformance unchanged — still 41/41 Go ≡ TS.)

## [0.1.0] — 2026-05-29

Initial reference release. Protocol version `veritrail/0.1`.

**Release status:** Go module **live** (`go get github.com/Daily-Nerd/veritrail/go@v0.1.0`, tag `go/v0.1.0`). npm publish **pending** (awaiting the npm org / final package-name decision: `veritrail` vs `@daily-nerd/veritrail`).

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
