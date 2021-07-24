// 前缀树 是一种树形数据结构，用于高效地存储和检索字符串数据集中的键。
// 这一数据结构有相当多的应用情景，例如自动补完和拼写检查。
interface ITrie {
  insert: (word: string) => Trie
  search: (word: string) => boolean
  startsWith: (prefix: string) => boolean
}

// 根节点不包含字符，除根节点意外每个节点只包含一个字符.
// 从根节点到某一个节点，路径上经过的字符连接起来，为该节点对应的字符串。

// 使用数组的话插入和查找的复杂度为O(n)，前缀的复杂度为O(n*k),n为数组长度
// 使用前缀树的话插入查找和前缀的复杂度为O(k),k为字符串长度

class TrieNode {
  val?: string
  isWord: boolean
  children: Map<string, TrieNode>

  constructor(val?: string, isWord: boolean = false) {
    this.val = val
    this.isWord = isWord
    this.children = new Map()
  }
}

// Trie:root节点(TrieNode或者Map<string,TriNode>)
// TrieNode:孩子节点(Map<string,Trie>或者Array<TrieNode>),结束flag,对应的值
class Trie implements ITrie {
  private root: TrieNode

  constructor() {
    this.root = new TrieNode()
  }

  insert(word: string): Trie {
    let root = this.root
    for (const letter of word) {
      if (!root.children.has(letter)) root.children.set(letter, new TrieNode(letter))
      root = root.children.get(letter)!
    }
    root.isWord = true
    return this
  }

  search(word: string): boolean {
    const node = this.traverse(word)
    return node !== null && node.isWord === true
  }

  startsWith(prefix: string): boolean {
    return this.traverse(prefix) !== null
  }

  /**
   *
   * @param word 返回val为word的TrieNode
   * @returns
   */
  private traverse(word: string): TrieNode | null {
    let root = this.root
    for (const letter of word) {
      const childNode = root.children.get(letter)
      if (!childNode) return null
      root = childNode
    }
    return root
  }
}

if (require.main === module) {
  const trie = new Trie()
  trie.insert('google')
  console.dir(trie, { depth: 4 })

  console.log(trie.search('google'))
  console.log(trie.startsWith('agoo'))
}

export { Trie }
