/**
 * @param {number} a
 * @param {number[]} b
 * @return {number}
 * 计算 a**b 对 1337 取模
 * @summary
 * 指数转换:a^[1,5,6,4] = a^4 * a^[1, 5, 6, 0] = a^4 * (a^[1, 5, 6])^10
 * 求模运算:(a∗b)%c=((a%c)∗(b%c))%c
 */
var superPow = function (a: number, b: number[]): number {
  if (!b.length) return 1
  const BASE = 1337
  // 计算 a**b%mod 的值
  const myPow = (a: number, b: number, mod = BASE): number => {
    if (b === 0) return 1
    else if (b % 2 === 1) return (a * myPow(a, b - 1)) % mod
    else return myPow(a, b / 2) ** 2 % mod
  }

  const last = b.pop()!
  return (myPow(a, last) * myPow(superPow(a, b), 10)) % BASE
}

console.log(superPow(1, [4, 3, 3, 8, 5, 2]))
