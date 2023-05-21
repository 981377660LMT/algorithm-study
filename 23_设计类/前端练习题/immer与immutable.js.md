深入探究 Immutable.js 的实现机制
https://juejin.cn/post/6844903679644958728
https://juejin.cn/post/6844903682891333640

- 可见对于一个 key 全是数字的 map，我们完全可以通过一颗 Vector Trie 来实现它，同时实现持久化数据结构。如果 key 不是数字怎么办呢？用一套映射机制把它转成数字就行了。 Immutable.js 实现了一个 hash 函数，可以把一个值转换成相应数字。
- 压缩的 Trie
- 解决哈希冲突
