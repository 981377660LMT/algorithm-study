class Node {
  constructor(
    public val: string,
    public isWord = false,
    public children: Map<string, Node> = new Map()
  ) {}
}

class StreamCheckerTrie {
  private root: Node

  constructor() {
    this.root = new Node('')
  }

  insert(s: string) {
    let rootP = this.root
    for (const letter of s) {
      if (!rootP.children.has(letter)) rootP.children.set(letter, new Node(letter))
      const next = rootP.children.get(letter)!
      rootP = next
    }
    rootP.isWord = true
  }

  search(s: string[] | string) {
    let rootP = this.root
    for (const letter of s) {
      const next = rootP.children.get(letter)
      if (!next) return false
      if (next.isWord) return true
      rootP = next
    }
    return false
  }
}

class StreamChecker {
  private trie: StreamCheckerTrie
  private queue: string[] = []

  constructor(words: string[]) {
    this.trie = new StreamCheckerTrie()
    // 反向插入
    words.forEach(word => this.trie.insert(word.split('').reverse().join('')))
    this.queue = []
  }

  /**
   *
   * @param letter 如果存在某些 k >= 1，
   * 可以用查询的最后 k个字符（按从旧到新顺序，包括刚刚查询的字母）拼写出给定字词表中的某一字词时，返回 true
   */
  query(letter: string): boolean {
    // 用dequeue快一些
    this.queue.unshift(letter)
    return this.trie.search(this.queue)
  }
}

const streamChecker = new StreamChecker(['cd', 'f', 'kl'])

console.log(streamChecker.query('a'))
console.log(streamChecker.query('b'))
console.log(streamChecker.query('c'))
console.log(streamChecker.query('d')) // cd在里面，true
