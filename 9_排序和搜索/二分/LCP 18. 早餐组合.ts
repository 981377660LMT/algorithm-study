/**
 *
 * @param staple 主食
 * @param drinks
 * @param x 小扣的计划选择一份主食和一款饮料，且花费不超过 x 元。请返回小扣共有多少种购买方案
 */
function breakfastNumber(staple: number[], drinks: number[], x: number): number {
  const mod = 10 ** 9 + 7
  staple.sort((a, b) => a - b)
  drinks.sort((a, b) => a - b)
  let res = 0
  let r = drinks.length - 1
  for (let l = 0; l < staple.length; l++) {
    while (r >= 0 && staple[l] + drinks[r] > x) r--
    if (r === -1) break
    res += r + 1
    res %= mod
  }

  return res
}

console.log(breakfastNumber([10, 20, 5], [5, 5, 2], 15))
