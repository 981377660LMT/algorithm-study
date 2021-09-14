/**
 * @param {number[]} arr
 * @param {number} k
 * @return {number}
 * 给你一个整数数组 arr，请你将该数组分隔为长度最多为 k 的一些（连续）子数组。
 * 分隔完成后，每个子数组的中的所有值都会变为该子数组中的最大值。
   返回将数组分隔变换后能够得到的元素最大和。
 */
const maxSumAfterPartitioning = function (arr: number[], k: number): number {
  // 记忆化递归
  const dp = (index: number, memo: Map<number, number>): number => {
    if (index >= arr.length) return 0
    if (memo.has(index)) return memo.get(index)!

    let res = 0
    let max = 0
    for (let delta = 1; delta <= k; delta++) {
      if (index + delta - 1 >= arr.length) continue
      max = Math.max(max, arr[index + delta - 1])
      res = Math.max(res, delta * max + dp(index + delta, memo))
    }

    memo.set(index, res)
    return res
  }

  return dp(0, new Map())
}

console.log(maxSumAfterPartitioning([1, 15, 7, 9, 2, 5, 10], 3))
// 因为 k=3 可以分隔成 [1,15,7] [9] [2,5,10]，结果为 [15,15,15,9,10,10,10]，和为 84，是该数组所有分隔变换后元素总和最大的。
export default 1
