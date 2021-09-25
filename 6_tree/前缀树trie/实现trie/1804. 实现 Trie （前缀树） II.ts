class TrieNode {
  val: string
  asPrefix: number // 插入word时，每一个结点+1
  asWord: number // 插入word结束时，统计一下
  children: Map<string, TrieNode>

  constructor(val: string) {
    this.val = val
    this.asPrefix = 0
    this.asWord = 0
    this.children = new Map()
  }
}

class Trie {
  private root: TrieNode
  constructor() {
    this.root = new TrieNode('')
  }

  // 将字符串 word 插入前缀树中
  insert(word: string): void {
    let root = this.root
    for (const char of word) {
      if (!root.children.has(char)) root.children.set(char, new TrieNode(char))
      root.children.get(char)!.asPrefix++
      root = root.children.get(char)!
    }
    root.asWord++
  }

  // 返回前缀树中字符串 word 的实例个数。
  countWordsEqualTo(word: string): number {
    let root = this.root
    for (const char of word) {
      if (!root.children.has(char)) return 0
      root = root.children.get(char)!
    }
    return root.asWord
  }

  // 返回前缀树中以 prefix 为前缀的字符串个数
  countWordsStartingWith(prefix: string): number {
    let root = this.root
    for (const char of prefix) {
      if (!root.children.has(char)) return 0
      root = root.children.get(char)!
    }
    return root.asPrefix
  }

  // 从前缀树中移除字符串 word
  erase(word: string): void {
    let root = this.root
    for (const char of word) {
      if (!root.children.has(char)) return
      root.children.get(char)!.asPrefix--
      root = root.children.get(char)!
    }
    root.asWord--
  }
}

export {}
