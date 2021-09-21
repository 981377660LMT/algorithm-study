// 一般解法
// 输出smalls中的字符串在big里出现的所有位置positions
// function multiSearch(big: string, smalls: string[]): number[][] {
//   const allIndexOf = function (str: string, searchElement: string) {
//     if (searchElement === '') return []
//     const res: number[] = []
//     let idx = str.indexOf(searchElement)
//     while (idx !== -1) {
//       res.push(idx)
//       idx = str.indexOf(searchElement, idx + 1)
//     }
//     return res
//   }
//   return smalls.map(v => allIndexOf(big, v))
// }

// indexof是顺序查找 而

class TrieNode {
  val: string
  index: number
  children: Map<string, TrieNode>

  constructor(val: string, index: boolean = false) {
    this.val = val
    this.index = -1
    this.children = new Map()
  }
}

class Trie {
  private root: TrieNode

  constructor() {
    this.root = new TrieNode('')
  }

  insert(str: string, index: number) {
    let root = this.root
    for (const char of str) {
      !root.children.has(char) && root.children.set(char, new TrieNode(char))
      root = root.children.get(char)!
    }
    root.index = index
  }

  search(str: string) {
    const res: number[] = []
    let root = this.root
    for (const char of str) {
      if (!root.children.has(char)) return res
      const childNode = root.children.get(char)!
      childNode.index !== -1 && res.push(childNode.index)
      root = childNode
    }
    return res
  }
}

function multiSearch(big: string, smalls: string[]): number[][] {
  const trie = new Trie()
  const res = Array.from<any, number[]>({ length: smalls.length }, () => [])

  smalls.forEach((small, index) => trie.insert(small, index))

  for (let i = 0; i < big.length; i++) {
    const hit = trie.search(big.slice(i))
    for (const index of hit) {
      res[index].push(i)
    }
  }

  return res
}

console.log(multiSearch('mississippi', ['is', 'ppi', 'hi', 'sis', 'i', 'ssippi']))
