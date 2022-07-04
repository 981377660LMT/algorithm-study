import assert from 'assert'

type BigIntMatrix = BigInt64Array[]

function mul(m1: BigIntMatrix, m2: BigIntMatrix, mod: bigint): BigIntMatrix {
  const [ROW, COL] = [m1.length, m2[0].length]
  const res = Array.from({ length: ROW }, () => new BigInt64Array(COL))
  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      for (let i = 0; i < ROW; i++) {
        res[r][c] = m1[r][i] * m2[i][c] + res[r][c]
        res[r][c] %= mod
      }
    }
  }

  return res
}

/**
 * @description 矩阵快速幂
 * @param base 矩阵(N*N)
 * @param exp 幂次
 * @param mod 每个数最后取的模
 */
function matqpow(base: BigIntMatrix, exp: bigint, mod: bigint): BigIntMatrix {
  const N = base.length
  let res = Array.from({ length: N }, () => new BigInt64Array(N))
  for (let r = 0; r < N; r++) res[r][r] = 1n

  while (exp) {
    if (exp & 1n) res = mul(res, base, mod)
    exp = exp >> 1n
    base = mul(base, base, mod)
  }

  return res
}

export { matqpow }

if (require.main === module) {
  const n = 876543210987654321n
  const MOD = BigInt(1e9 + 7)
  let res = [new BigInt64Array([2n]), new BigInt64Array([1n]), new BigInt64Array([1n])]
  const trans = [
    new BigInt64Array([1n, 1n, 1n]),
    new BigInt64Array([1n, 0n, 0n]),
    new BigInt64Array([0n, 1n, 0n]),
  ]

  const tmp = matqpow(trans, n - 3n, MOD)
  res = mul(tmp, res, MOD)
  assert.strictEqual(Number(res[0][0]), 639479200)
}
