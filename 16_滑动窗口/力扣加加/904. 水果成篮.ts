import { resolve4 } from 'dns/promises'

/**
 * @param {number[]} fruits
 * @return {number}
 * @description 你有两个篮子，每个篮子可以携带任何数量的水果，但你希望每个篮子只携带一种类型的水果。
 * @summary 求只包含两种元素的最长连续子序列
 */
const totalFruit = function (fruits: number[]): number {
  let res = 0
  let l = 0
  const map = new Map<number, number>()

  for (let r = 0; r < fruits.length; r++) {
    const cur = fruits[r]
    map.set(cur, map.get(cur)! + 1 || 1)
    while (map.size > 2) {
      l++
      const pre = fruits[l - 1]
      const count = map.get(pre)!
      if (count === 1) map.delete(pre)
      else map.set(pre, count - 1)
    }
    res = Math.max(res, r - l + 1)
  }

  return res
}

console.log(totalFruit([3, 3, 3, 1, 2, 1, 1, 2, 3, 3, 4]))
// 5
