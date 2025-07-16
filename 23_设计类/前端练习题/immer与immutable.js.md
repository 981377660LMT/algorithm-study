深入探究 Immutable.js 的实现机制
https://juejin.cn/post/6844903679644958728
https://juejin.cn/post/6844903682891333640

- 可见对于一个 key 全是数字的 map，我们完全可以通过一颗 Vector Trie 来实现它，同时实现持久化数据结构。如果 key 不是数字怎么办呢？用一套映射机制把它转成数字就行了。 Immutable.js 实现了一个 hash 函数，可以把一个值转换成相应数字。
- 压缩的 Trie
- 解决哈希冲突

---

immer.js - 实现不可变数据的新思路
https://blog.rexskz.info/immerjs-a-new-way-to-implement-immutable-data.html

- immer.js 介绍
- 内部实现
  - Proxy 与优化：懒的 proxy
    https://github.com/immerjs/immer/blob/9064d26aaaa4e6d5cc447b1b140f4c891286e813/src/proxy.js#L134
  - Scope 与 Patch 的概念
    https://github.com/immerjs/immer/blob/9064d26aaaa4e6d5cc447b1b140f4c891286e813/src/patches.js#L106-L138
  - 对 Async-await 的处理
  - 暴力的 ES5 版本
- 性能测试
- 注意事项
  - 不要对 draftState 本身重新赋值
  - State 中不要出现相同的引用或循环引用
  - 不要一次性读写大量数据
  - 尽量将更多操作包在一个 produce 中
  - 不要直接返回 undefined
- 值得一提的事情

---

- [函数式编程所倡导使用的「不可变数据结构」如何保证性能？](https://www.zhihu.com/question/53804334)

  - structural sharing
  - lazy evaluation (COW)
  - transient

- [immer 是如何实现 immutable 的](https://zhuanlan.zhihu.com/p/602961293)

- https://bytetech.info/articles/7451882553474514954?from=lark_all_search#LxHJdQBo9owNnbx2587cjIBnngg

---

函数式数据结构只是把问题提前了而已，实现复杂度的确非常的高，而且单线程性能确实会有一点牺牲。
但是换来的是`无畏并发`、无锁多线程、多核调度、无限扩容、分布式计算，Time Travel 回滚，全部都赚回来了.
采用可变数据结构的遇到适合不可变数据结构的算法会需要付出很多额外代价；但采用不可变数据结构的遇到适合可变数据结构的算法基本是直接放弃
保障不了性能，不然这个技术术语在宣传上的关键词就不是「不可变」而是「零开销」了

---

如果我们的应用并不是那种 内存复杂度 型的应用，使用 immutable data 的开销可能就不是瓶颈，真正的瓶颈会在 IO 上，那么我们就可以忽略这些开销，而享受 immutable 带给我们的红利：`线程安全（在并发处理时候，这个的价值是巨大的，不仅保证安全，而且可以无锁化，实现更好的性能。）、GC 友好、清晰的数据处理流程、更清晰的代码和更少的 BUG。`
