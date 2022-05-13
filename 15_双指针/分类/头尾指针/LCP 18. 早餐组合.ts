const MOD = 1e9 + 7

/**
 *
 * @param staple 主食
 * @param drinks
 * @param x 小扣的计划选择一份主食和一款饮料，且花费不超过 x 元。请返回小扣共有多少种购买方案
 */
function breakfastNumber(staple: number[], drinks: number[], x: number): number {
  staple.sort((a, b) => a - b)
  drinks.sort((a, b) => a - b)

  let res = 0
  let right = drinks.length - 1

  for (let left = 0; left < staple.length; left++) {
    while (right >= 0 && staple[left] + drinks[right] > x) right--
    if (right === -1) break
    res += right + 1
    res %= MOD
  }

  return res
}

console.log(breakfastNumber([10, 20, 5], [5, 5, 2], 15))

export {}
