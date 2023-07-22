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
