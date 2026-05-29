/**
 * cost-canon command — JCS of a cost block with string-encoded integers.
 *
 * Validates that tokens, usd_micros, and wall_ms are JSON strings (not numbers).
 * If any magnitude is a JSON number, error: {"error":"cost_must_be_string_int"}.
 *
 * CRITICAL: After JSON.parse, typeof number === "number". We must check
 * the parsed type — not the string representation — since JSON.parse
 * distinguishes numbers from strings.
 */
import { jcsBytes } from './jcs.js';

export interface CostBlock {
  tokens: unknown;
  usd_micros: unknown;
  wall_ms: unknown;
  rail_ref: unknown;
}

export interface CostInput {
  cost: CostBlock;
}

export interface CostOutput {
  canonical_hex: string;
}

export function costCanon(input: CostInput): CostOutput | { error: string } {
  const { cost } = input;

  // Validate: tokens, usd_micros, wall_ms must be strings, not numbers
  const magnitudeFields: Array<keyof CostBlock> = ['tokens', 'usd_micros', 'wall_ms'];
  for (const field of magnitudeFields) {
    if (typeof cost[field] === 'number') {
      return { error: 'cost_must_be_string_int' };
    }
  }

  const canonical = jcsBytes(cost);
  return { canonical_hex: canonical.toString('hex') };
}
