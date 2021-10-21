/**
 *
 * @param nums 1 <= arr.length <= 10^5  1 <= arr[i] <= 10000
 * @param m 将数组分为最多 m 子数组，每个子数组可以删除最多一个数，求 m 个子数组和的最小值
 * @description 我们应该贪心地删除子数组中最大的数。
 */
function minTime(nums: number[], m: number): number {
  let l = 0
  let r = nums.reduce((pre, cur) => pre + cur, 0)
  while (l <= r) {
    const mid = (l + r) >> 1
    if (calMinSegment(mid) > m) l = mid + 1
    else r = mid - 1
  }

  return l

  function calMinSegment(threshold: number): number {
    let res = 1
    let sum = 0 // 一开始默认删除
    let max = nums[0]
    for (const num of nums.slice(1)) {
      if (sum + Math.min(max, num) > threshold) {
        res++
        sum = 0
        max = num
      } else {
        sum += Math.min(max, num)
        max = Math.max(max, num)
      }
    }

    return res
  }
}

console.log(minTime([1, 2, 3, 3], 2))
