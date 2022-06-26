type NumberMatrix = number[][]
type BigIntMatrix = BigInt64Array[]

/**
 * @description 矩阵快速幂
 * @param base 矩阵(N*N)
 * @param exp 幂次
 * @param mod 每个数最后取的模
 */
function matqpow(base: NumberMatrix, exp: number, mod = 1e9 + 7): NumberMatrix {
  const N = base.length
  const bigMod = BigInt(mod)
  let bigBase = base.map(row => BigInt64Array.from(row, BigInt))
  let bigRes = Array.from({ length: N }, () => new BigInt64Array(N))
  for (let r = 0; r < N; r++) bigRes[r][r] = 1n

  while (exp) {
    if (exp % 2 === 1) bigRes = mul(bigRes, bigBase)
    exp = Math.floor(exp / 2)
    bigBase = mul(bigBase, bigBase)
  }

  return bigRes.map(row => Array.from(row, Number))

  function mul(m1: BigIntMatrix, m2: BigIntMatrix): BigIntMatrix {
    const [ROW, COL] = [m1.length, m2[0].length]
    const res = Array.from({ length: ROW }, () => new BigInt64Array(COL))
    for (let r = 0; r < ROW; r++) {
      for (let c = 0; c < COL; c++) {
        for (let i = 0; i < ROW; i++) {
          res[r][c] = m1[r][i] * m2[i][c] + res[r][c]
          res[r][c] %= bigMod
        }
      }
    }

    return res
  }
}

export { matqpow }

if (require.main === module) {
  const base = [
    [1, 1],
    [1, 0],
  ]

  console.log(matqpow(base, 1e9))
}
