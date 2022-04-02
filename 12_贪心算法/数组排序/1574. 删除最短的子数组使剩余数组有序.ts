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
  let l = 0
  let r = arr.length - 1
  let res = Infinity

  while (l + 1 < arr.length && arr[l] <= arr[l + 1]) l++
  if (l === arr.length - 1) return 0

  while (r > 0 && arr[r] >= arr[r - 1]) r--

  // 删除0到r-1  或者删除l+1到arr.length-1
  res = Math.min(r, arr.length - l - 1)

  // 删中间：画平行线
  let i = 0
  while (i <= l && r <= arr.length - 1) {
    if (arr[i] <= arr[r]) {
      // 删除i+1 到 r-1
      res = Math.min(res, r - 1 - i)
      i++
    } else {
      r++
    }
  }

  return res
}

console.log(findLengthOfShortestSubarray([1, 2, 3, 10, 4, 2, 3, 5]))

export {}
