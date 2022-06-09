**Radix 树(Gin 框架的路由不使用反射，基于 Radix 树，内存占用少)**
与 Trie 不同的是，它对 Trie 树进行了空间优化，只有一个子节点的中间节点将被压缩
假如树中的一个节点是父节点的唯一子节点(the only child)的话，那么该子节点将会与父节点进行合并
即元素个数不是太多，但是元素之间通常有很长的相同前缀时很适合采用 radix tree 来存储

在构建 IP 路由(ip-routing)的应用方面 radix tree 也使用广泛，`因为 IP 通常具有大量相同的前缀`； 另外 radix tree 在倒排索引方面也使用广泛。

https://github.com/mjschultz/py-radix
https://blog.csdn.net/lindorx/article/details/103431530
