class TrieNode {
  val: string
  count: number
  children: TrieNode[]

  constructor(val: string) {
    this.val = val
    this.count = 0
    this.children = Array<TrieNode>(26)
  }
}

class Trie {
  private root: TrieNode

  constructor() {
    this.root = new TrieNode('')
  }

  insert(str: string) {
    let root = this.root
    for (const char of str) {
      const key = char.codePointAt(0)! - 97
      !root.children[key] && (root.children[key] = new TrieNode(char))
      root = root.children[key]
    }
    root.count++
  }

  search(str: string) {
    let root = this.root
    for (const char of str) {
      const key = char.codePointAt(0)! - 97
      if (!root.children[key]) return 0
      root = root.children[key]
    }
    return root.count
  }
}

class WordsFrequency {
  private trie: Trie

  constructor(book: string[]) {
    this.trie = new Trie()
    book.forEach(word => this.trie.insert(word))
  }

  get(word: string): number {
    return this.trie.search(word)
  }
}

const wordsFrequency = new WordsFrequency(['i', 'have', 'an', 'apple', 'he', 'have', 'a', 'pen'])

console.log(wordsFrequency.get('you'))
console.log(wordsFrequency.get('have'))
console.log(wordsFrequency.get('an'))
console.log(wordsFrequency.get('apple'))
console.log(wordsFrequency.get('pen'))
