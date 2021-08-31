<!-- 专门用于处理子串查询 -->

<!-- 字符串问题 -->
<!-- 模式匹配 -->
<!-- 编译原理 -->

算法的复杂度瓶颈在字符串查找，并且字符串有很多公共前缀，就可以用前缀树优化。
例如: 异或判断是否大于 m
树的每一个节点存储的是：n 个数中，从根节点到当前节点形成的前缀有多少个是一样的，

解决的问题:**这个字符串的前缀是否存在于某个数据结构中?**

遍历 trie 子节点技巧:

```JS
  // 求子节点所有的weight(trie词频统计)
  private _sum(node: Node): number {
    let res = node.weight
    for (const next of node.children.values()) {
      res += this._sum(next)
    }
    return res
  }
```

```JS
  // 字符串模糊匹配
  // 用index比每次slice要好
  match(node: TrieNode, word: string, index: number): boolean {
    // 递归终点
    if (index === word.length) return node.isWord

    if (word[index] !== '.') {
      const next = node.children.get(word[index])
      if (!next) return false
      return this.match(node.children.get(word[index])!, word, index + 1)
    } else {
      // '.'时遍历所有孩子 如果一个true则为true
      for (const next of node.children.values()) {
        if (this.match(next, word, index + 1)) return true
      }
      return false
    }
  }
```

trie 的局限性：空间消耗大
改进：压缩字典树/**三分搜索字典树(>/===/<),每个节点只有三个孩子，但是牺牲时间换空间**

<!-- 一个场景，在一个输入框输入内容，怎么更加高效的去提示用户你输入的信息，
// 举个例子，你输入天猫，那么对应的提示信息是天猫商城，天猫集团，
// 这个信息如何最快的获取，有没有不需要发请求的方式来实现？

// 提示：

// 数据请求：防抖、节流
// 数据存储处理：Trie树 -->
