// 前缀树是一个树状的数据结构，
// 用于高效地存储和检索一系列字符串的前缀。
// 前缀树有许多应用，如自动补全和拼写检查

class TrieNode<V = string> {
  // value: V // 存储结点值
  // word = '' // 当前结点的单词
  preCount = 0
  wordCount = 0
  children: Map<string, TrieNode> = new Map()
}

class Trie {
  readonly root: TrieNode = new TrieNode()

  // 将字符串 word 插入前缀树中
  insert(word: string): void {
    let root = this.root
    for (const char of word) {
      if (!root.children.has(char)) root.children.set(char, new TrieNode())
      root.children.get(char)!.preCount++
      root = root.children.get(char)!
    }

    root.wordCount++
    // root.word = word
  }

  // 返回前缀树中字符串 word 的实例个数。
  // 不过 更快的方法是直接在dfs中获取结点的信息 而不是重新遍历
  countWord(word: string): number {
    let root = this.root
    for (const char of word) {
      if (!root.children.has(char)) return 0
      root = root.children.get(char)!
    }

    return root.wordCount
  }

  // 返回前缀树中以 prefix 为前缀的字符串个数
  // 不过 更快的方法是直接在dfs中获取结点的信息 而不是重新遍历
  countPre(prefix: string): number {
    let root = this.root
    for (const char of prefix) {
      if (!root.children.has(char)) return 0
      root = root.children.get(char)!
    }

    return root.preCount
  }

  // 从前缀树中移除 `1个` 字符串 word
  // !需要保证word在前缀树中
  discard(word: string): void {
    let root = this.root
    for (const char of word) {
      if (!root.children.has(char)) return
      root.children.get(char)!.preCount--
      root = root.children.get(char)!
    }

    root.wordCount--
  }
}

export { Trie, TrieNode }
