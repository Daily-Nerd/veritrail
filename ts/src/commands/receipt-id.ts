/**
 * receipt-id command — derive the content address of a Receipt.
 *
 * receipt_id = "u" + base64url_nopad( multihash_sha256( JCS(receipt) ) )
 *
 * Error: if `receipt` contains a `receipt_id` key → {"error":"receipt_id_must_be_absent"}
 */
import { jcsBytes } from './jcs.js';
import { digestBytes } from './digest.js';

export interface ReceiptIdInput {
  receipt: Record<string, unknown>;
}

export interface ReceiptIdOutput {
  canonical_hex: string;
  receipt_id: string;
}

export function receiptId(input: ReceiptIdInput): ReceiptIdOutput | { error: string } {
  const receipt = input.receipt;

  // Guard: receipt_id must not be present
  if (Object.prototype.hasOwnProperty.call(receipt, 'receipt_id')) {
    return { error: 'receipt_id_must_be_absent' };
  }

  const canonical = jcsBytes(receipt);
  const receipt_id = digestBytes(canonical);

  return {
    canonical_hex: canonical.toString('hex'),
    receipt_id,
  };
}
