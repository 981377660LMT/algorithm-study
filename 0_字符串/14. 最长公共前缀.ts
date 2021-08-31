// 解法1：逐个比较
// 时间复杂度：O(s)，s 是所有字符串中字符数量的总和
var longestCommonPrefix1 = (strs: string) => {
  if (!strs) return ''
  if (strs.length === 1) return strs[0]

  let prevs = strs[0]
  for (let i = 1; i < strs.length; i++) {
    let j = 0
    for (; j < prevs.length && j < strs[i].length; j++) {
      if (prevs.charAt(j) !== strs[i].charAt(j)) break
    }
    prevs = prevs.substring(0, j)
    if (prevs === '') return ''
  }
  return prevs
}
// 解法二：仅需最大、最小字符串的最长公共前缀
// 最小 ab 与最大 ac 的最长公共前缀一定也是 abc 、 abcd 的公共前缀
var longestCommonPrefix = function (strs: string) {
  if (!strs) return ''
  if (strs.length === 1) return strs[0]
  let min = 0,
    max = 0
  for (let i = 1; i < strs.length; i++) {
    if (strs[min] > strs[i]) min = i
    if (strs[max] < strs[i]) max = i
  }
  for (let j = 0; j < strs[min].length; j++) {
    if (strs[min].charAt(j) !== strs[max].charAt(j)) {
      return strs[min].substring(0, j)
    }
  }
  return strs[min]
}
