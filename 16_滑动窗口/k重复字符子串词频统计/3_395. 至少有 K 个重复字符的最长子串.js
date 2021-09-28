/**
 * @param {string} s
 * @param {number} k
 * @return {number}
 * 请你找出 s 中的最长子串， 要求该子串中的每一字符出现次数都不少于 k 。返回这一子串的长度
 * @summary
 * 增加限制：限制字符串包含的字符种类会方便一些
 */
var longestSubstring = function (s, k) {
  let res = 0

  for (let i = 1; i <= 26; i++) {
    let l = 0
    let r = 0
    let type = 0 // 现在有几种了
    let okCount = 0 // 字符出现次数都不少于 k的字符数
    const counter = new Map()

    while (r < s.length) {
      if ((counter.get(s[r]) || 0) === 0) type++
      if ((counter.get(s[r]) || 0) === k - 1) okCount++
      counter.set(s[r], (counter.get(s[r]) || 0) + 1)
      r++

      while (type > i) {
        if (counter.get(s[l]) === 1) type--
        if (counter.get(s[l]) === k) okCount--
        counter.set(s[l], counter.get(s[l]) - 1)
        l++
      }

      if (type === i && okCount === i) {
        res = Math.max(res, r - l) // r-l 代表移动的次数 即子串的长度
      }
    }
  }

  return res
}
console.log(longestSubstring('aaabb', 3))
