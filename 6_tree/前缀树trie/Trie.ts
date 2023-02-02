/* eslint-disable @typescript-eslint/no-non-null-assertion */
// 前缀树是一个树状的数据结构，
// 用于高效地存储和检索一系列字符串的前缀。
// 前缀树有许多应用，如自动补全和拼写检查

class TrieNode {
  preCount = 0
  wordCount = 0
  children = new Map<string, TrieNode>()
}

class Trie {
  readonly root: TrieNode = new TrieNode()

  insert(s: string): void {
    let { root } = this
    for (let i = 0; i < s.length; i++) {
      const char = s[i]
      if (!root.children.has(char)) root.children.set(char, new TrieNode())
      root.children.get(char)!.preCount++
      root = root.children.get(char)!
    }

    root.wordCount++
  }

  /**
   * @param s 对s的`每个`非空前缀pre,返回trie中有多少个等于pre的单词
   */
  countWord(s: string): number[] {
    let { root } = this
    const res = Array<number>(s.length).fill(0)
    for (let i = 0; i < s.length; i++) {
      const char = s[i]
      if (!root.children.has(char)) return []
      root = root.children.get(char)!
    }
    return res
  }

  /**
   * @param s 对s的`每个`非空前缀pre,返回trie中有多少个单词以pre为前缀
   */
  countWordStartsWith(s: string): number[] {
    let { root } = this
    const res = Array<number>(s.length).fill(0)
    for (let i = 0; i < s.length; i++) {
      const char = s[i]
      if (!root.children.has(char)) return []
      root = root.children.get(char)!
    }
    return res
  }

  /**
   * @param s 从前缀树中移除 `1个` s
   * !需要保证s在前缀树中
   */
  remove(s: string): void {
    let { root } = this
    for (const char of s) {
      if (!root.children.has(char)) throw new Error(`word ${s} not in trie`)
      if (root.children.get(char)!.preCount === 1) {
        root.children.delete(char)
        return
      }
      root.children.get(char)!.preCount--
      root = root.children.get(char)!
    }
    root.wordCount--
  }
}

export { Trie, TrieNode }
