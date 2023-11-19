/* eslint-disable @typescript-eslint/no-non-null-assertion */
// 前缀树是一个树状的数据结构，
// 用于高效地存储和检索一系列字符串的前缀。
// 前缀树有许多应用，如自动补全和拼写检查

class TrieNode<K> {
  preCount = 0
  wordCount = 0
  children = new Map<K, TrieNode<K>>()
}

class Trie<K extends PropertyKey> {
  readonly root: TrieNode<K> = new TrieNode()

  insert(s: ArrayLike<K>): TrieNode<K> {
    if (!s.length) return this.root
    let { root } = this
    for (let i = 0; i < s.length; i++) {
      const char = s[i]
      if (!root.children.has(char)) root.children.set(char, new TrieNode())
      root.children.get(char)!.preCount++
      root = root.children.get(char)!
    }
    root.wordCount++
    return root
  }

  /**
   * @param s 从前缀树中移除 `1个` s
   * !需要保证s在前缀树中
   */
  remove(s: ArrayLike<K>): TrieNode<K> {
    if (!s.length) return this.root
    let { root } = this
    for (let i = 0; i < s.length; i++) {
      const char = s[i]
      if (!root.children.has(char)) throw new Error(`word ${s} not in trie`)
      root.children.get(char)!.preCount--
      root = root.children.get(char)!
    }
    root.wordCount--
    return root
  }

  find(s: ArrayLike<K>): TrieNode<K> | undefined {
    if (!s.length) return undefined
    let { root } = this
    for (let i = 0; i < s.length; i++) {
      const char = s[i]
      if (!root.children.has(char)) return undefined
      root = root.children.get(char)!
    }
    return root
  }

  enumerate(s: ArrayLike<K>, f: (index: number, node: TrieNode<K>) => boolean | void): void {
    if (!s.length) return
    let { root } = this
    for (let i = 0; i < s.length; i++) {
      const char = s[i]
      if (!root.children.has(char)) return
      root = root.children.get(char)!
      if (f(i, root)) return
    }
  }
}

export { Trie }

if (require.main === module) {
  // 1804. 实现 Trie （前缀树） II
  // https://leetcode.cn/problems/implement-trie-ii-prefix-tree/description/
  class Trie2 {
    private readonly _trie = new Trie<string>()

    insert(word: string): void {
      this._trie.insert(word)
    }

    erase(word: string): void {
      this._trie.remove(word)
    }

    countWordsEqualTo(word: string): number {
      let res = 0
      this._trie.enumerate(word, (i, node) => {
        if (i === word.length - 1) res = node.wordCount
      })
      return res
    }

    countWordsStartingWith(prefix: string): number {
      let res = 0
      this._trie.enumerate(prefix, (i, node) => {
        if (i === prefix.length - 1) res = node.preCount
      })
      return res
    }
  }

  /**
   * Your Trie object will be instantiated and called as such:
   * var obj = new Trie()
   * obj.insert(word)
   * var param_2 = obj.countWordsEqualTo(word)
   * var param_3 = obj.countWordsStartingWith(prefix)
   * obj.erase(word)
   */
}
