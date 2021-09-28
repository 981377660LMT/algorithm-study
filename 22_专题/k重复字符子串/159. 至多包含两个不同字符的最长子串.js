/**
 * @param {string} s
 * @return {number}
 * 这里要记录字符的种类
 */
var lengthOfLongestSubstringTwoDistinct = function (s) {
  let l = 0
  let r = 0
  let max = 0
  let type = 0 // 现在有几种了
  const counter = new Map() // 这里换成map 记录count 可以解决有k种字符的最长字串

  while (r < s.length) {
    if ((counter.get(s[r]) || 0) === 0) type++
    counter.set(s[r], (counter.get(s[r]) || 0) + 1)
    r++
    while (type > 2) {
      if (counter.get(s[l]) === 1) type--
      counter.set(s[l], counter.get(s[l]) - 1)
      l++
    }
    max = Math.max(max, r - l)
  }

  return max
}
console.log(lengthOfLongestSubstringTwoDistinct('eceba'))
