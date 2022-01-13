/**
 * @param {number[]} fruits
 * @return {number}
 * @description 你有两个篮子，每个篮子可以携带任何数量的水果，但你希望每个篮子只携带一种类型的水果。
 * @summary 求只包含两种元素的最长连续子序列
 */
const totalFruit = function (fruits: number[]): number {
  let res = 0
  let left = 0
  const counter = new Map<number, number>()

  for (let right = 0; right < fruits.length; right++) {
    const cur = fruits[right]
    counter.set(cur, (counter.get(cur) ?? 0) + 1)

    while (counter.size > 2) {
      const pre = fruits[left]
      const preCount = counter.get(pre)!
      if (preCount === 1) counter.delete(pre)
      else counter.set(pre, preCount - 1)
      left++
    }

    res = Math.max(res, right - left + 1)
  }

  return res
}

console.log(totalFruit([3, 3, 3, 1, 2, 1, 1, 2, 3, 3, 4]))
// 5

export default 1
