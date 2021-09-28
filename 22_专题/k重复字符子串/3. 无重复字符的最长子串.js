/**
 * @param {string} s
 * @return {number}
 * @summary
 * 这里要记录字符的最大重复数量
 */
var lengthOfLongestSubstring = function (s) {
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

console.log(lengthOfLongestSubstring('bbbbb'))

// var lengthOfLongestSubstring2 = function (s) {
//   let left = 0
//   let right = 0
//   let max = 0
//   const set = new Set()

//   while (right <= s.length - 1) {
//     if (!set.has(s[right])) {
//       set.add(s[right])
//       max = Math.max(set.size, max)
//       right++
//     } else {
//       left++
//       set.delete(s[left - 1])
//     }
//   }

//   return max
// }
