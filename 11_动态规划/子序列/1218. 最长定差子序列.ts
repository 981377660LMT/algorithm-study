type Value = number
type Count = number
/**
 * @param {number[]} arr
 * @param {number} difference
 * @return {number}
 * 找出并返回 arr 中最长等差子序列的长度，该子序列中相邻元素之间的差等于 difference 。
 */
var longestSubsequence = function (arr: number[], difference: number): number {
  // 结尾
  const map = new Map<Value, Count>()
  for (const num of arr) {
    const pre = num - difference
    if (map.has(pre)) map.set(num, map.get(pre)! + 1)
    else map.set(num, 1)
  }

  return Math.max(...map.values(), 1)
}

console.log(longestSubsequence([1, 2, 3, 4], 1))

export default 1
