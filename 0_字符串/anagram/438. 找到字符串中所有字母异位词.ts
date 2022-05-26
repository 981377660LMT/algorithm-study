/**
 *
 * @param s
 * @param p
 * 找到 s 中所有 p 的 异位词 的子串，返回这些子串的起始索引。
 * 不考虑答案输出的顺序
 */
function findAnagrams(s: string, p: string): number[] {
  const [n, m] = [s.length, p.length]
  if (n < m) return []
  const res: number[] = []
  let l = 0
  let r = 0
  let lack = m
  const counter = new Map() // 缺少的单词统计
  for (const char of p) {
    counter.set(char, (counter.get(char) || 0) + 1)
  }

  while (r < n) {
    if ((counter.get(s[r]) || 0) > 0) lack--
    counter.set(s[r], (counter.get(s[r]) || 0) - 1)
    r++

    if (lack === 0) res.push(l)

    // 免去了初始化
    if (r - l === p.length) {
      if ((counter.get(s[l]) || 0) >= 0) lack++
      counter.set(s[l], (counter.get(s[l]) || 0) + 1)
      l++
    }
  }

  return res
}

console.log(findAnagrams('cbaebabacd', 'abc'))

export {}
