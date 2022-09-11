/* eslint-disable no-shadow */
/* eslint-disable no-param-reassign */
// !注意用Uint64还是int64
// !尽量不要用js的快速幂 大数运算非常慢(和python一样的速度)

import assert from 'assert'

type BigIntMatrix = BigUint64Array[]

function mul(m1: BigIntMatrix, m2: BigIntMatrix, mod: bigint): BigIntMatrix {
  const [ROW, COL] = [m1.length, m2[0].length]
  const res = Array.from({ length: ROW }, () => new BigUint64Array(COL))
  for (let i = 0; i < ROW; i++) {
    for (let j = 0; j < COL; j++) {
      for (let k = 0; k < ROW; k++) {
        res[i][j] += m1[i][k] * m2[k][j]
        res[i][j] %= mod
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
  const pow2: BigIntMatrix[] = [base]
  return inner(base, exp, mod)

  function inner(base: BigIntMatrix, exp: bigint, mod: bigint): BigIntMatrix {
    const n = base.length
    let res = Array.from({ length: n }, () => new BigUint64Array(n))
    for (let r = 0; r < n; r++) {
      res[r][r] = 1n
    }

    let bit = 0
    while (exp) {
      if (exp & 1n) {
        res = mul(res, pow2[bit], mod)
      }

      exp >>= 1n
      bit++

      if (bit === pow2.length) {
        const last = pow2[bit - 1]
        pow2.push(mul(last, last, mod))
      }
    }

    return res
  }
}

export { matqpow }

if (require.main === module) {
  const n = 876543210987654321n
  const MOD = BigInt(1e9 + 7)
  let res = [new BigUint64Array([2n]), new BigUint64Array([1n]), new BigUint64Array([1n])]
  const trans = [
    new BigUint64Array([1n, 1n, 1n]),
    new BigUint64Array([1n, 0n, 0n]),
    new BigUint64Array([0n, 1n, 0n])
  ]

  const tmp = matqpow(trans, n - 3n, MOD)
  res = mul(tmp, res, MOD)
  assert.strictEqual(Number(res[0][0]), 639479200)
}
