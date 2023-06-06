https://kopricky.github.io/code/DataStructure_Advanced/compressed_trie.html
https://www.jianshu.com/p/4b9609222848

Patricia Trie(Compressed Trie) 相对于普通的 R-way Trie，一方面把只有单节点的细长分支压缩成了一个节点，另一方面其基于 2 进制比较，空间复杂度与字符集大小 R 无关（严格的说是和 logR 相关）
每个节点只有两个子节点（而普通的 Trie 里每个节点要开 R 个空间用来存子节点）

---

如果数据太多，内存里存不下整个 Trie 该怎么办？
分布式 Trie

解决的思路是把 Trie 分散放在多台机器上。`可以对前两个字符做一致性 hash 来路由机器`，比如以 ab 开头的词都在机器 1，以 ac 开头的词都在机器 3。
当然，假如这是一场技术面试，那么随之而来又会产生新的问题：假如有数据倾斜怎么办，有访问热点怎么办？

---

**Also known as radix tree.**
**Radix 树(Gin 框架的路由不使用反射，基于 Radix 树，内存占用少)**
与 Trie 不同的是，它对 Trie 树进行了空间优化，只有一个子节点的中间节点将被压缩
假如树中的一个节点是父节点的唯一子节点(the only child)的话，那么该子节点将会与父节点进行合并
即元素个数不是太多，但是元素之间通常有很长的相同前缀时很适合采用 radix tree 来存储

在构建 IP 路由(ip-routing)的应用方面 radix tree 也使用广泛，`因为 IP 通常具有大量相同的前缀`； 另外 radix tree 在倒排索引方面也使用广泛。

https://github.com/mjschultz/py-radix
https://blog.csdn.net/lindorx/article/details/103431530

---

https://xlinux.nist.gov/dads/HTML/patriciatree.html
