class TrieNode {
  val: string
  times: number // 大于0表示是isWord
  children: Map<string, TrieNode>

  constructor(val: string) {
    this.val = val
    this.times = 0
    this.children = new Map()
  }
}

class Trie {
  private root: TrieNode

  constructor() {
    this.root = new TrieNode('')
  }

  insert(word: string, times: number): void {
    let root = this.root
    for (const char of word) {
      if (!root.children.has(char)) root.children.set(char, new TrieNode(char))
      root = root.children.get(char)!
    }
    root.times += times
  }

  // 回溯搜索所有以word为前缀的单词
  search(word: string): [string, number][] {
    let root = this.root
    const res: [string, number][] = []
    for (const char of word) {
      if (!root.children.has(char)) return []
      root = root.children.get(char)!
    }
    this.bt(root, res, [word])
    return res
  }

  private bt(root: TrieNode, res: [string, number][], path: string[]) {
    if (root.times > 0) {
      res.push([path.join(''), root.times])
    }

    for (const child of root.children.values()) {
      path.push(child.val)
      this.bt(child, res, path)
      path.pop()
    }
  }
}

class AutocompleteSystem {
  private trie: Trie
  private curSearch: string

  constructor(sentences: string[], times: number[]) {
    this.trie = new Trie()
    for (let i = 0; i < sentences.length; i++) {
      this.trie.insert(sentences[i], times[i])
    }
    this.curSearch = ''
  }

  // 输出历史热度前三的具有相同前缀的句子。
  // 如果满足条件的句子个数少于 3，将它们全部输出
  // 如果输入了特殊字符'#'，意味着句子结束了，请返回一个空集合
  // 字符只会是小写英文字母（'a' 到 'z' ），空格（' '）和特殊字符（'#'）。
  input(c: string): string[] {
    const K = 3
    // 输入结束后插入
    if (c === '#') {
      this.trie.insert(this.curSearch, 1)
      this.curSearch = ''
      return []
    } else {
      this.curSearch += c
      const res = this.trie.search(this.curSearch)
      return res
        .sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0]))
        .map(item => item[0])
        .slice(0, K)
    }
  }

  static main() {
    const autocompleteSystem = new AutocompleteSystem(
      ['i love you', 'island', 'ironman', 'i love leetcode'],
      [5, 3, 2, 2]
    )

    // 由于 ' ' 的 ASCII 码是 32 而 'r' 的 ASCII 码是 114，
    // 所以 "i love leetcode" 在 "ironman" 前面。
    // 同时我们只输出前三的句子，所以 "ironman" 被舍弃。
    console.log(autocompleteSystem.input('i')) // ["i love you", "island","i love leetcode"]
    console.log(autocompleteSystem.input(' ')) //  ["i love you","i love leetcode"]
    console.log(autocompleteSystem.input('a')) //  [],没有句子有前缀 "i a"。
  }
}
AutocompleteSystem.main()
export {}
