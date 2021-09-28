/**
 * @param {string} s
 * @param {string} t
 * @return {string}
 * 返回 s 中涵盖 t 所有字符的最小子串。
 * 如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 ""
 */
var minWindow = function (s, t) {
  if (s.length < t.length) return ''
  let l = 0
  let r = 0
  let min = Infinity
  let lack = t.length // 缺少的数量，到0就统计
  const counter = new Map() //  缺少的对应关系,负数代表多了
  let res = ''

  for (const char of t) {
    counter.set(char, (counter.get(char) || 0) + 1)
  }

  while (r < s.length) {
    if ((counter.get(s[r]) || 0) > 0) lack--
    counter.set(s[r], (counter.get(s[r]) || 0) - 1)
    r++
    while (lack === 0) {
      if (min > r - l) {
        min = r - l
        res = s.slice(l, r)
      }
      if (counter.get(s[l]) === 0) lack++
      counter.set(s[l], counter.get(s[l]) + 1)
      l++
    }
  }

  return res
}
