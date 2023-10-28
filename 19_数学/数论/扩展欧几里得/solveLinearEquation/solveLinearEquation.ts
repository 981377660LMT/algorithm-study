/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */

/**
 * a*x + b*y = c 的通解为
 * `x = (c/g)*x0 + (b/g)*k, y = (c/g)*y0 - (a/g)*k`.
 * 其中 g = gcd(a,b) 且需要满足 g|c.
 * x0 和 y0 是 ax+by=g 的一组特解（即 exgcd(a,b) 的返回值).
 *
 * 为方便讨论，这里要求输入的 a b c 必须为正整数.
 *
 * @returns
 * - 正整数解的个数（无解时为 -1,无正整数解时为 0)
 * - x 取最小正整数时的解 x1 y1,此时 y1 是最大正整数解
 * - y 取最小正整数时的解 x2 y2,此时 x2 是最大正整数解
 *
 * @param allowZero 是否允许解为 0.默认为 false.
 */
function solveLinearEquation(
  a: number,
  b: number,
  c: number,
  allowZero = false
): {
  n: number
  x1: number
  y1: number
  x2: number
  y2: number
} {
  let { gcd, x: x0, y: y0 } = exgcd(a, b)

  // 无解
  if (c % gcd) return { n: -1, x1: 0, y1: 0, x2: 0, y2: 0 }

  a /= gcd
  b /= gcd
  c /= gcd
  x0 *= c
  y0 *= c

  let x1 = x0 % b
  if (allowZero) {
    if (x1 < 0) x1 += b
  } else if (x1 <= 0) x1 += b
  const k1 = (x1 - x0) / b
  const y1 = y0 - k1 * a

  let y2 = y0 % a
  if (allowZero) {
    if (y2 < 0) y2 += a
  } else if (y2 <= 0) y2 += a
  const k2 = (y0 - y2) / a
  const x2 = x0 + k2 * b

  // 无正整数解
  if (y1 <= 0) return { n: 0, x1, y1, x2, y2 }

  // k 越大 x 越大
  return { n: k2 - k1 + 1, x1, y1, x2, y2 }
}

function exgcd(a: number, b: number): { gcd: number; x: number; y: number } {
  let x = 1
  let y = 0
  let x1 = 0
  let y1 = 1
  let a1 = a
  let b1 = b
  while (b1) {
    const q = Math.floor(a1 / b1)
    let tmp = x - q * x1
    x = x1
    x1 = tmp
    tmp = y - q * y1
    y = y1
    y1 = tmp
    tmp = a1 - q * b1
    a1 = b1
    b1 = tmp
  }
  return { gcd: a1, x, y }
}

export { solveLinearEquation }

if (require.main === module) {
  function solve(a: number, b: number, c: number): number[] {
    const { n, x1, y1, x2, y2 } = solveLinearEquation(a, b, c)
    if (n === -1) return [-1]
    if (n === 0) return [x1, y2]
    return [n, x1, y2, x2, y1]
  }

  const test = [
    [2, 11, 100],
    [3, 18, 6],
    [192, 608, 17],
    [19, 2, 60817],
    [11, 45, 14],
    [19, 19, 810],
    [98, 76, 5432]
  ]

  const expected = [[4, 6, 2, 39, 8], [2, 1], [-1], [1600, 1, 18, 3199, 30399], [34, 3], [-1], [2, 12, 7, 50, 56]]

  test.forEach((t, i) => {
    const res = solve(...t)
    if (JSON.stringify(res) !== JSON.stringify(expected[i])) {
      console.log(`Test ${i} fail:`, t, res)
      process.exit(1)
    }
  })
  console.log('Test done.')
}
