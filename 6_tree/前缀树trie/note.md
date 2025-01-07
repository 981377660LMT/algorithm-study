使用 Trie 的提示:
**所有字符串总长度<=2e5**
`trie的时间空间复杂度等于所有字符总长度`

**trie 中存储了所有 words 的信息**

算法的复杂度瓶颈在字符串查找，并且字符串有很多公共前缀，就可以用前缀树优化。
例如: 异或判断是否大于 m
树的每一个节点存储的是：n 个数中，从根节点到当前节点形成的前缀有多少个是一样的，

处理单词的前缀时，应该在 **dfs 中向下处理节点的信息 而不是每次调 api 对前缀重新遍历**

---

可以使用数组, 而不是对象来表示 **TrieNode**
可以节省空间.
**functools.lru_cache 中链表的实现就是用的四个元素的 list.**

---

利用 Trie 树的变种优化带参数路由的匹配
https://blog.rexskz.info/use-variant-trie-to-optimize-router-match-with-params.html
Gin 用了一个叫 Radix Tree 的数据结构，基本可以理解为“压缩节点后的 Trie”

https://engineering.linecorp.com/ja/blog/simple-tries
