import assert from 'assert'

const INF = 2e15

/**
 * 朴素的辗转相除法求最大公约数.
 */
function gcd(a: number, b: number): number {
  if (a < 0) a = -a
  if (b < 0) b = -b
  while (b) {
    const mod = a % b
    a = b
    b = mod
  }
  return a
}

/**
 * 求两个数的最小公倍数.
 * @param clamp 答案超过 {@link clamp} 时返回 {@link clamp}.
 */
function lcm(a: number, b: number, clamp = INF): number {
  if (!a || !b) return 0
  if (a < 0) a = -a
  if (b < 0) b = -b
  if (a >= clamp || b >= clamp) return clamp
  const gcd_ = gcd(a, b)
  a = Math.floor(a / gcd_)
  if (a >= Math.ceil(clamp / b)) return clamp
  return a * b
}

/**
 * 扩展gcd的迭代实现.
 * 扩展欧几里得求 `a*x + b*y = gcd(a,b)` 的一组解.
 * !注意返回的gcd可能为负数.
 * @see https://oi-wiki.org/math/number-theory/gcd
 */
function exgcd(a: number, b: number): [gcd: number, x: number, y: number] {
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
  return [a1, x, y]
}

/**
 * 模逆元.
 * 求出逆元 `inv` 满足 `a*inv ≡ target (mod mod)`.
 * 如果不存在逆元则返回 `undefined`.
 * @param [target=1] 目标值.默认为1.
 */
function modInv(a: number, mod: number, target = 1): number | undefined {
  if (target === 1) {
    const [gcd_, x] = exgcd(a, mod)
    if (gcd_ !== 1) return undefined
    const res = x % mod
    return res < 0 ? res + mod : res
  }
  let [gcd_, x] = exgcd(a, mod)
  if (target % gcd_) return undefined
  x *= target / gcd_
  mod /= gcd_
  const res = x % mod
  return res < 0 ? res + mod : res
}

/**
 * 快速求两个int32的gcd.
 * @alias binaryGcd(二进制gcd).
 * @see https://nyaannyaan.github.io/library/trial/fast-gcd.hpp
 */
function gcdInt32(a: number, b: number) {
  if (!a || !b) return a + b
  if (a < 0) a = -a
  if (b < 0) b = -b

  const ctz1 = 31 - Math.clz32(a) // __builtin_ctz(a)
  const ctz2 = 31 - Math.clz32(b)
  a >>>= ctz1
  b >>>= ctz2
  while (a ^ b) {
    const ctz = 31 - Math.clz32(a - b)
    const f = a > b
    const max = f ? a : b
    b = f ? b : a
    a = (max - b) >>> ctz
  }

  return ctz1 < ctz2 ? a << ctz1 : b << ctz2
}

export { gcd, lcm, exgcd, modInv, gcdInt32, gcdInt32 as binaryGcd, gcdInt32 as fastGcd }

if (require.main === module) {
  console.log(exgcd(3, 6))

  assert.strictEqual(modInv(3, 10), 7)
  const MOD = 998244353
  const INV2 = (MOD + 1) / 2 // 结论:2的逆元(模998244353)为499122177
  assert.strictEqual(modInv(2, MOD), INV2)

  console.time('gcd')
  for (let i = 0; i < 1e7; ++i) {
    gcd(123456789, 987654321)
  }
  console.timeEnd('gcd')

  console.time('gcdInt32')
  for (let i = 0; i < 1e7; ++i) {
    gcdInt32(123456789, 987654321)
  }
  console.timeEnd('gcdInt32')
}
