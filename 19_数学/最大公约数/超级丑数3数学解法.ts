/**
 * @param {number} n
 * @param {number[]} primes
 * @return {number}
 * @description
   给你三个数字 a，b，c，你需要找到第 n 个（n 从 0 开始）有序序列的值，这个有序序列是由 a，b，c 的整数倍构成的。
   @summary 采用二分法 先确定解空间 再最左能力二分
 */
const nthSuperUglyNumber = function (a: number, b: number, c: number, n: number): number {
  const _GCD = (a: number, b: number): number => (b === 0 ? a : GCD(b, a % b))
  const GCD = (...arr: number[]) => arr.reduce(_GCD)
  const _LCM = (a: number, b: number): number => (a * b) / GCD(a, b)
  const LCM = (...arr: number[]) => arr.reduce(_LCM)

  const possible = (mid: number) =>
    ~~mid / a +
      ~~mid / b +
      ~~mid / c -
      ~~mid / LCM(a, b) -
      ~~mid / LCM(a, c) -
      ~~mid / LCM(b, c) +
      ~~mid / LCM(a, b, c) >=
    n

  let l = 0
  let r = n * Math.max(a, b, c)
  while (l <= r) {
    const mid = ~~((l + r) / 2)
    if (possible(mid)) r = mid - 1
    else l = mid + 1
  }

  return l
}

console.log(nthSuperUglyNumber(2, 5, 7, 8))
