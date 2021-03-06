/**
 * @param {number[]} arr
 * @return {number}
 * @description
 * 请你删除一个子数组（可以为空），使得 arr 中剩下的元素是 非递减 的。
 *
 * 注意是删除最短的子数组而不是子序列
 * 返回满足题目要求删除的最短子数组的长度。
 * 必须连续，那就考虑滑动窗口。
 */
function findLengthOfShortestSubarray(arr: number[]): number {
  // 删左边/删右边/删中间
  const n = arr.length
  let left = 0
  let right = n - 1
  let res = Infinity

  while (left + 1 < n && arr[left] <= arr[left + 1]) left++
  if (left === n - 1) return 0

  while (right > 0 && arr[right] >= arr[right - 1]) right--

  // 删除0到r-1  或者删除l+1到n-1
  res = Math.min(right, n - left - 1)

  // 删中间：画平行线，看交点
  let i = 0
  while (i <= left && right <= n - 1) {
    if (arr[i] <= arr[right]) {
      // 删除i+1 到 r-1
      res = Math.min(res, right - 1 - (i + 1) + 1)
      i++
    } else {
      right++
    }
  }

  return res
}

console.log(findLengthOfShortestSubarray([1, 2, 3, 10, 4, 2, 3, 5]))

export {}
