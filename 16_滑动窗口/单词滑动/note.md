变位词 anagram 问题

```JS
  const [n, m] = [s1.length, s2.length]
  if (m < n) return ...
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

    if (lack === 0) ...

    if (r - l === s1.length) {
      if ((counter.get(s2[l]) || 0) >= 0) lack++
      counter.set(s2[l], (counter.get(s2[l]) || 0) + 1)
      l++
    }
  }

  return ...
```
