// 判断 s2 是否包含 s1 的排列。如果是，返回 true ；否则，返回 false
// 类似于438. 找到字符串中所有字母异位词
function checkInclusion(s1: string, s2: string): boolean {
  const [n, m] = [s1.length, s2.length]
  if (m < n) return false
  let l = 0
  let r = 0
  let lack = n
  const counter = new Map<string, number>() // 缺少的单词统计
  for (const char of s1) {
    counter.set(char, (counter.get(char) || 0) + 1)
  }

  while (r < m) {
    if ((counter.get(s2[r]) || 0) > 0) lack--
    counter.set(s2[r], (counter.get(s2[r]) || 0) - 1)
    r++

    if (lack === 0) return true

    // 免去了初始化
    if (r - l === s1.length) {
      if ((counter.get(s2[l]) || 0) >= 0) lack++
      counter.set(s2[l], (counter.get(s2[l]) || 0) + 1)
      l++
    }
  }

  return false
}

console.log(checkInclusion('ab', 'eidbaooo'))
console.log(checkInclusion('ab', 'eidboaoo'))
