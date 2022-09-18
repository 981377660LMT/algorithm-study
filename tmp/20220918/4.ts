class TrieNode {
  // word = '' // 当前结点的单词
  preCount = 0
  wordCount = 0
  children: Map<string, TrieNode> = new Map()
}

class Trie {
  readonly root: TrieNode = new TrieNode()

  /**
   * @param word 将字符串 word 插入前缀树中
   */
  insert(word: string): void {
    let { root } = this
    for (let i = 0; i < word.length; i++) {
      const char = word[i]
      if (!root.children.has(char)) root.children.set(char, new TrieNode())
      root.children.get(char)!.preCount++
      root = root.children.get(char)!
    }

    root.wordCount++
    // root.word = word
  }

  /**
   * @param word 树中有多少个单词word
   *
   * 更快的方法是直接在dfs中获取结点的信息 而不是重新遍历
   */
  countWord(word: string): number {
    let { root } = this
    for (let i = 0; i < word.length; i++) {
      const char = word[i]
      if (!root.children.has(char)) return 0
      root = root.children.get(char)!
    }

    return root.wordCount
  }

  /**
   * @param prefix 树中有多少个以prefix为前缀的单词
   */
  countWordStartsWith(prefix: string): number {
    let { root } = this
    for (let i = 0; i < prefix.length; i++) {
      const char = prefix[i]
      if (!root.children.has(char)) return 0
      root = root.children.get(char)!
    }

    return root.preCount
  }

  /**
   * @param word 从前缀树中移除 `1个` word
   *
   * !需要保证word在前缀树中
   */
  remove(word: string): void {
    let { root } = this
    for (const char of word) {
      if (!root.children.has(char)) throw new Error(`word ${word} not in trie`)
      root.children.get(char)!.preCount--
      root = root.children.get(char)!
    }

    root.wordCount--
  }
}

function sumPrefixScores(words: string[]): number[] {
  const trie = new Trie()
  words.forEach(word => trie.insert(word))
  const res = []
  for (let i = 0; i < words.length; i++) {
    let cur = 0
    for (let j = 0; j < words[i].length; j++) {
      const pre = words[i].slice(0, j + 1)
      cur += trie.countWordStartsWith(pre)
    }
    res.push(cur)
  }

  return res
}
