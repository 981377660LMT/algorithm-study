const MOD = 1337

/**
 * @param {number} a
 * @param {number[]} b
 * @return {number}
 * 计算 a**b 对 1337 取模
 * @summary
 * 指数转换:a^[1,5,6,4] = a^4 * a^[1, 5, 6, 0] = a^4 * (a^[1, 5, 6])^10
 * 求模运算:(a∗b)%c=((a%c)∗(b%c))%c
 */
function superPow(a: number, b: number[]): number {
  if (b.length === 0) return 1

  // 计算 a**b%mod 的值
  const last = b.pop()!
  return (qpow(a, last, MOD) * qpow(superPow(a, b), 10, MOD)) % MOD
}

function qpow(a: number, b: number, mod: number): number {
  let res = 1

  while (b) {
    if (b & 1) {
      res *= a
      res %= mod
    }

    a **= 2
    b = b >> 1
    a %= mod
  }

  return res
}

console.log(qpow(2, 3, 100))

console.log(superPow(1, [4, 3, 3, 8, 5, 2]))
console.log(superPow(2147483647, [2, 0, 0]))

export {}
