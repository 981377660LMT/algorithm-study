/* eslint-disable func-names */
// !和python一样慢

import { matqpow } from './matqpow'

const MOD = BigInt(1e9 + 7)

/**
 * @param {number} row
 * @param {number} col
 * @param {number[][]} position
 * @return {number}
 */
function electricityExperiment(row: number, col: number, position: number[][]): number {
  const trans = Array.from({ length: row }, () => new BigUint64Array(row))
  for (let r = 0; r < row; r++) {
    trans[r][r] = 1n
    if (r !== 0) trans[r][r - 1] = 1n
    if (r !== row - 1) trans[r][r + 1] = 1n
  }

  const n = position.length
  position.sort((a, b) => a[1] - b[1])
  for (let i = 0; i < n - 1; i++) {
    const [row1, col1] = position[i]
    const [row2, col2] = position[i + 1]
    const rowDiff = Math.abs(row1 - row2)
    const colDiff = Math.abs(col1 - col2)
    if (rowDiff > colDiff) return 0
  }

  let res = 1n // !注意这里也要用BigInt
  for (let i = 0; i < n - 1; i++) {
    const [row1, col1] = position[i]
    const [row2, col2] = position[i + 1]
    const colDiff = Math.abs(col1 - col2)
    res = (res * cal(row1, row2, colDiff)) % MOD
  }
  return Number(res)

  function cal(row1: number, row2: number, k: number): bigint {
    const resTrans = matqpow(trans, BigInt(k), MOD)
    return resTrans[row1][row2]
  }
}

if (require.main === module) {
  console.time('1')
  console.timeEnd('1')
}

export {}
