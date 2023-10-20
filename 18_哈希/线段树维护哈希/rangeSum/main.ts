/* eslint-disable max-len */

function rangeSum(start: number, end: number): number {
  if (start >= end) return 0
  return ((end - start) * (start + end - 1)) / 2
}

/**
 * 区间平方和.
 */
function rangeSquareSum(start: number, end: number): number {
  if (start >= end) return 0
  const tmp1 = (end * (end - 1) * (2 * end - 1)) / 6
  const tmp2 = (start * (start - 1) * (2 * start - 1)) / 6
  return tmp1 - tmp2
}

/**
 * 区间立方和.
 */
function rangeCubeSum(start: number, end: number): number {
  if (start >= end) return 0
  const tmp1 = (end * (end - 1)) / 2
  const tmp2 = (start * (start - 1)) / 2
  return tmp1 * tmp1 - tmp2 * tmp2
}

/**
 * 区间异或和.
 */
function rangeXorSum(start: number, end: number): number {
  if (start >= end) return 0
  const preXor = (upper: number): number => {
    const mod = upper % 4
    if (mod === 0) return upper
    if (mod === 1) return 1
    if (mod === 2) return upper + 1
    return 0
  }
  return preXor(end - 1) ^ preXor(start - 1)
}

/**
 * 区间k次幂和.
 */
function rangePowKSum(start: number, end: number, k: number, mod = 1e9 + 7): number {
  if (start >= end || mod === 1) return 0
  const cal = (n: number): number => {
    let sum = 1
    let p = k
    const s = Math.floor(Math.log2(n)) - 1
    for (let d = s; d >= 0; d--) {
      sum *= p + 1
      p *= p
      if ((n >>> d) & 1) {
        sum += p
        p *= k
      }
      sum %= mod
      p %= mod
    }
    return sum
  }
  return cal(end) - cal(start)
}

/**
 * 从区间和与区间长度还原区间.
 */
function recoverIntervalFrom(sum_: number, length: number): [first: number, last: number, ok: boolean] {
  if (!length) return [0, 0, false]
  if ((2 * sum_) % length) return [0, 0, false]
  const div = (2 * sum_) / length
  if ((div + 1 - length) & 1) return [0, 0, false]
  const x = (div + 1 - length) / 2
  return [x, x + length - 1, true]
}

export { rangeSum, rangeSquareSum, rangeCubeSum, rangeXorSum, rangePowKSum, recoverIntervalFrom }

if (require.main === module) {
  // check
  for (let i = 0; i < 100; i++) {
    for (let j = i; j < 100; j++) {
      if (rangeSum(i, j) !== Array.from({ length: j - i }, (_, k) => k + i).reduce((a, b) => a + b, 0)) {
        throw new Error(`rangeSum(${i},${j})`)
      }
      if (rangeSquareSum(i, j) !== Array.from({ length: j - i }, (_, k) => k + i).reduce((a, b) => a + b * b, 0)) {
        throw new Error(`rangeSquareSum(${i},${j})`)
      }
      if (rangeCubeSum(i, j) !== Array.from({ length: j - i }, (_, k) => k + i).reduce((a, b) => a + b * b * b, 0)) {
        throw new Error(`rangeCubeSum(${i},${j})`)
      }
      if (rangeXorSum(i, j) !== Array.from({ length: j - i }, (_, k) => k + i).reduce((a, b) => a ^ b, 0)) {
        throw new Error(`rangeXorSum(${i},${j})`)
      }
    }
  }
}
