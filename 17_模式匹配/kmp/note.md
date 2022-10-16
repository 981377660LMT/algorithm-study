kmp 的 next 数组的使用

1. fallback 后前进：匹配不成功，j 往右直接走到最长公共后缀的位置
2. 匹配：匹配成功，j 往右走一步
   理解：https://www.ruanyifeng.com/blog/2013/05/Knuth%E2%80%93Morris%E2%80%93Pratt_algorithm.html

```JS
  const next = getNext(short)
  let hit = 0

  for (let i = 0; i < long.length; i++) {
    while (hit > 0 && long[i] !== short[hit]) {
      hit = next[hit - 1]  // 关键：如果不匹配，那么needle串前进到后缀匹配的地方，hit个数变为最大首尾公共前缀的长度，即不要把"搜索位置"移回已经比较过的位置，继续把它向`后移`
    }

    if (haystack[i] === short[hit]) hit++

    // 找到头了
    if (hit === short.length) {
      return i - short.length + 1
    }
  }
```

https://leetcode.cn/problems/find-the-index-of-the-first-occurrence-in-a-string/solution/shua-chuan-lc-shuang-bai-po-su-jie-fa-km-tb86/
暴力的匹配每次失败时,longer 串的指针都要回到起点后再+1,
而 `kmp 算法不需要移动原串的指针`,且利用 `shorter 串已经匹配部分的相同前后缀加速匹配`
