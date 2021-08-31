/**
 * @param {number[]} nums
 * @return {number}
 * 请你计算并返回 nums 中任意两个数之间汉明距离的总和。
 * @summary 两两统计不行 解法：按位统计
 * 其中第 i 位上的汉明距离之和是：res[i] * (nums.length - res[i])
 * @description 这里的每位计数可以用 比特位计数优化
 */
var totalHammingDistance = function (nums: number[]): number {
  const len = nums.length
  // 根据题目要求，不超过10^9，所以30位就可以了
  const res = Array<number>(30).fill(0)
  for (const num of nums) {
    let mask = 1
    for (let bit = 0; bit < 30; bit++) {
      if (num & mask) res[bit]++
      mask = mask << 1
    }
  }

  console.log(res)

  return res.reduce((pre, cur) => pre + cur * (len - cur), 0)
}

console.log(totalHammingDistance([4, 14, 2]))
console.log(0x01)

export {}
