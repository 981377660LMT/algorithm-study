1. 技巧： `AC自动机 linkWord dp`

   下面这段代码的复杂度是多少？
   `O(nsqrt(∑))`
   自然根号：fail 指针每次命中，都至少有一个比指针深度更长的单词出现，因此每个位置最坏情况下不超过 O(sqrt(∑))次命中

   ```go
   pos := int32(0)
   for i, char := range target {
     pos = trie.Move(pos, char)
     // 对当前文本串后缀，找到每个匹配的模式串.
     for cur := pos; cur != 0; cur = trie.LinkWord(cur) {
       dp[i+1] = min(dp[i+1], dp[int32(i)+1-nodeDepth[cur]]+nodeCosts[cur])
     }
   }
   ```

   https://leetcode.cn/problems/construct-string-with-minimum-cost/solutions/2833826/aczi-dong-ji-you-hua-dp-onljie-fa-by-vcl-t6fx/
