// 删除子数组

/**
 * !请你删除一个子数组（可以为空），使得 arr 中剩下的元素是 非递减 的。
 *
 * 注意是删除最短的子数组而不是子序列
 * 返回满足题目要求删除的最短子数组的长度。
 * !双指针，一共三种情况 1、开头一段+末尾一段 2、开头一段 3、末尾一段
 */
function findLengthOfShortestSubarray(arr: number[]): number {
  // 删左边/删右边/删中间
  const n = arr.length
  let left = 0
  let right = arr.length - 1
  let res = 2e15

  while (left + 1 < n && arr[left] <= arr[left + 1]) left++
  while (right > 0 && arr[right] >= arr[right - 1]) right--

  // 删除0到r-1  或者删除l+1到n-1
  res = Math.min(right, n - left - 1)
  if (res === 0) return 0

  // 删中间：
  // 只需遍历左边数组的每一个位置，找到右边数组相应的非递减的端点，我们就可以得到答案
  // 可以用二分解决 注意到滑动的单调性 考虑双指针
  let j = right
  for (let i = 0; i <= left; i++) {
    while (j < n && arr[i] > arr[j]) j++
    res = Math.min(res, j - i - 1)
  }

  return res
}

if (require.main === module) {
  console.log(findLengthOfShortestSubarray([1, 2, 3, 10, 4, 2, 3, 5]))
}

export {}
