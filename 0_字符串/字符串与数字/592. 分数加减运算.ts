/**
 * @param {string} expression
 * @return {string}
 * "+"或者"-"可以看成分子的一部分，这样每个分数都是+a/b或者-a/b的形式,更容易计算
 */
var fractionAddition = function (expression: string): string {
  const _GCD = (a: number, b: number): number => (b === 0 ? a : GCD(b, a % b))
  const GCD = (...arr: number[]) => arr.reduce(_GCD)
  const _LCM = (a: number, b: number): number => (a * b) / GCD(a, b)
  const LCM = (...arr: number[]) => arr.reduce(_LCM)

  const divisor = [...expression.matchAll(/\/([-+]?\d+)/g)].map(m => parseInt(m[1]))
  const lcm = LCM(...divisor)
  const dividend = [...expression.matchAll(/([-+]?\d+)\//g)]
    .map(m => parseInt(m[1]))
    .map((m, i) => (m * lcm) / divisor[i])

  console.log(dividend)
  console.log(divisor)
  const d1 = dividend.reduce((pre, cur) => pre + cur)
  const gcd = GCD(d1, lcm)

  const res = `${d1 / gcd}/${lcm / gcd}`
  return res.includes('-') ? '-' + res.replace('-', '') : res
}

console.log(fractionAddition('-1/2+1/2'))
// 输出: "0/1"
export default 1
