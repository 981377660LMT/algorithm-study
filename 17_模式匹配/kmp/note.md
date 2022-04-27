kmp 的 next 数组的使用

1. fallback 后前进：匹配不成功，j 往右直接走到最长公共后缀的位置
2. 匹配：匹配成功，j 往右走一步

```JS
  const next = getNext(needle)
  let hit = 0

  for (let i = 0; i < pattern.length; i++) {
    while (hit > 0 && pattern[i] !== needle[hit]) {
      hit = next[hit - 1]  // 关键：如果不匹配，那么needle串前进到后缀匹配的地方，hit个数变为最大首尾公共前缀的长度，即不要把"搜索位置"移回已经比较过的位置，继续把它向`后移`
    }

    if (haystack[i] === needle[hit]) hit++

    // 找到头了
    if (hit === needle.length) {
      return i - needle.length + 1
    }
  }
```
