// 745. 前缀和后缀搜索后缀修饰的Trie
// 745. 前缀和后缀搜索

class TrieNode {
  index: number // 对应单词下标
  children: Map<string, TrieNode> = new Map()

  constructor(index = 0) {
    this.index = index
  }
}

// Trie:root节点(TrieNode或者Map<string,TriNode>)
// TrieNode:孩子节点(Map<string,Trie>或者Array<TrieNode>),结束flag,对应的值
class Trie {
  readonly root: TrieNode = new TrieNode()

  insert(word: string, index: number) {
    if (!word) return
    let { root } = this
    for (const letter of word) {
      if (!root.children.has(letter)) root.children.set(letter, new TrieNode())
      root = root.children.get(letter)!
      root.index = index // 更新weight
    }
  }

  search(prefix: string, suffix: string): number {
    const target = `${suffix}#${prefix}`
    let { root } = this
    for (const char of target) {
      if (!root.children.has(char)) return -1
      root = root.children.get(char)!
    }
    return root.index
  }
}

// 我们将在单词查找树中插入
// '#apple', 'e#apple', 'le#apple', 'ple#apple', 'pple#apple', 'apple#apple'。
// 然后对于 prefix = "ap", suffix = "le" 这样的查询，
// 我们可以通过查询单词查找树找到 le#ap。
class WordFilter {
  private trie: Trie

  constructor(words: string[]) {
    this.trie = new Trie()
    for (let i = 0; i < words.length; i++) {
      const word = words[i]
      let suffix = ''
      for (let j = word.length; j >= 0; j--) {
        suffix = word.slice(j, word.length)
        this.trie.insert(`${suffix}#${word}`, j)
      }
    }
  }

  /**
   *
   * @param prefix
   * @param suffix
   * 返回词典中具有前缀 prefix 和后缀suffix 的单词的下标。
   * 如果存在不止一个满足要求的下标，返回其中 最大的下标 。
   * 如果不存在这样的单词，返回 -1 。
   */
  f(prefix: string, suffix: string): number {
    return this.trie.search(prefix, suffix)
  }

  static main() {
    const wordFilter = new WordFilter(['apple'])
    console.log(wordFilter.f('a', 'e')) // 0
  }
}

export {}
