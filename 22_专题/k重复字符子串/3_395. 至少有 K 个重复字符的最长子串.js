/**
 * @param {string} s
 * @param {number} k
 * @return {number}
 */
var longestSubstring = function (s, k) {
  let l = 0
  let r = 0
  let max = 0
  let overlap = 0 // 现在的最大重复数量
  const counter = new Map() // 这里换成map 记录count 可以解决有k种字符的最长字串

  while (r < s.length) {
    if ((counter.get(s[r]) || 0) > 0) overlap++
    counter.set(s[r], (counter.get(s[r]) || 0) + 1)
    r++
    while (overlap > 0) {
      if (counter.get(s[l]) > 1) overlap--
      counter.set(s[l], counter.get(s[l]) - 1)
      l++
    }
    max = Math.max(max, r - l)
  }

  return max
}
