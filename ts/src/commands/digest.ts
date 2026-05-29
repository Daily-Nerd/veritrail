/**
 * digest command — sha256 hash of raw bytes, returned as pinned hash-string.
 */
import { createHash } from 'node:crypto';
import { encodeHashstring } from './hashstring.js';

export interface DigestInput {
  bytes_b64: string;
}

export interface DigestOutput {
  hashstr: string;
}

export function digestBytes(rawBytes: Buffer): string {
  const sha = createHash('sha256').update(rawBytes).digest();
  return encodeHashstring(sha);
}

export function digest(input: DigestInput): DigestOutput {
  const raw = Buffer.from(input.bytes_b64, 'base64');
  return { hashstr: digestBytes(raw) };
}
