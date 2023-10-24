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
  const { gcd, x: x0, y: y0 } = exgcd(a, b)
  if (c % gcd) return { n: -1, x1: 0, y1: 0, x2: 0, y2: 0 }
  a /= gcd
  b /= gcd
  c /= gcd
  let x1 = x0 * c
  let y1 = y0 * c
  x1 %= b
  if (allowZero ? x1 < 0 : x1 <= 0) x1 += b
  const k1 = Math.floor((x1 - x0) / b)
  y1 = y0 - k1 * a
  let x2 = x0 * c
  let y2 = y0 * c
  y2 %= a
  if (allowZero ? y2 < 0 : y2 <= 0) y2 += a
  const k2 = Math.floor((y0 - y2) / a)
  x2 = x0 + k2 * b
  if (y1 <= 0) return { n: 0, x1, y1, x2, y2 }
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
