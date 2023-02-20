AC 自动机作用: 多模式匹配

```
文本串:s
模式串:patterns
```

1. 求每个 pattern 在 s 中出现的位置和次数.
   - https://leetcode.cn/problems/multi-search-lcci/
   - https://leetcode.cn/problems/stream-of-characters/
2. 求使得 s 中的子串都不在 patterns 中的最少替换次数和方案(处理敏感词).
   - https://atcoder.jp/contests/abc268/tasks/abc268_h
3. 将 s 分割成若干个单词,使得每个部分都在 patterns 中出现.
   - https://atcoder.jp/contests/jag2017autumn/tasks/jag2017autumn_h
4. AC 自动机上 dp:
   统计不含有被禁止的模式串的字符串个数.
   dp[i][state]表示当前看到第 i 位,处于 Trie 上的状态 state(pos)时,不含禁止的字符串个数.
   - https://yukicoder.me/problems/no/1269
   - https://leetcode.cn/problems/find-all-good-strings/
