class Trie {
  constructor() {
    this.root = {}
  }

  insert(word) {
    let node = this.root
    for (let c of word) {
      if (node[c] == null) node[c] = {}
      node = node[c]
    }
    node.isWord = true
  }

  traverse(word) {
    let node = this.root
    for (let c of word) {
      node = node[c]
      if (node == null) return null
    }
    return node
  }

  search(word) {
    const node = this.traverse(word)
    return node != null && node.isWord === true
  }

  startsWith(prefix) {
    return this.traverse(prefix) != null
  }
}

const trie = new Trie()

trie.insert('google')
trie.insert('react')
console.log(trie.search('a'))
console.log(trie.startsWith('goo'))
console.dir(trie.root, { depth: null })
