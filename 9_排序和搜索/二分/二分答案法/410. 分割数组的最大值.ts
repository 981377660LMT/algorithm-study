/**
 * @param {number[]} nums 1 <= nums.length <= 1000
 * @param {number} m
 * @return {number}
 * 给定一个非负整数数组 nums 和一个整数 m ，你需要将这个数组分成 m 个非空的连续子数组。
 * 设计一个算法使得这 m 个子数组各自和的最大值最小
 * 返回这个和
 * @summary
 * 非负整数数组:和的单调性=>二分查找
 */
const splitArray = function (nums: number[], m: number): number {
  let l = Math.max.apply(null, nums)
  let r = nums.reduce((pre, cur) => pre + cur, 0)
  while (l <= r) {
    const mid = (l + r) >> 1
    const count = countNGT(nums, mid)
    if (count > m) {
      l = mid + 1
    } else {
      r = mid - 1
    }
  }

  return l

  /**
   *
   * @param nums
   * @param mid
   * @returns 每个子数组和不超过mid 至少需要分成多少个子数组
   */
  function countNGT(nums: number[], mid: number): number {
    let res = 1
    let curSum = 0

    for (const num of nums) {
      if (num + curSum > mid) {
        curSum = num
        res++
      } else {
        curSum += num
      }
    }

    return res
  }
}

console.log(splitArray([7, 2, 5, 10, 8], 2))

export default 1
