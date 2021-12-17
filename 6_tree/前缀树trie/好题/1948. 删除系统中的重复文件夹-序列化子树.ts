// 题目要求：如果子树的哈希(结构)相同，则删除

// 1. 为了得到某一个节点的子文件夹的「结构」，
// 我们应当首先遍历完成该节点的所有子节点，
// 再回溯遍历该节点本身。这就对应着多叉树的后序遍历。
// 子树使用序列化来记录在每个结点

// 2. 在回溯到某节点时，我们需要将该节点的「结构」存储下来，
// 记录在某一「数据结构」中，以便于与其它节点的「结构」进行比较。

class TrieNode {
  subtreeHash: string
  children: Map<string, TrieNode>
  constructor() {
    this.subtreeHash = ''
    this.children = new Map()
  }
}

class Trie {
  root: TrieNode

  constructor() {
    this.root = new TrieNode()
  }

  insert(path: string[]): void {
    let root = this.root
    for (const char of path) {
      if (!root.children.has(char)) {
        root.children.set(char, new TrieNode())
      }

      root = root.children.get(char)!
    }
  }
}

function deleteDuplicateFolder(paths: string[][]): string[][] {
  const trie = new Trie()
  paths.forEach(path => trie.insert(path))

  const hashCounter = new Map<string, number>()
  genHash(trie.root, hashCounter)
  // console.dir(trie, { depth: null })

  const res: string[][] = []
  bt(trie.root, [], hashCounter)
  return res

  function genHash(root: TrieNode, counter: Map<string, number>): void {
    // 这句话很关键
    if (root.children.size === 0) return

    const sb: string[] = []
    for (const [childName, child] of root.children.entries()) {
      genHash(child, counter)
      sb.push(`${childName}(${child.subtreeHash})`)
    }

    sb.sort()
    root.subtreeHash = sb.join('')
    counter.set(root.subtreeHash, (counter.get(root.subtreeHash) || 0) + 1)
  }

  function bt(root: TrieNode, path: string[], counter: Map<string, number>): void {
    if (counter.get(root.subtreeHash)! >= 2) return
    if (path.length > 0) res.push(path.slice())

    for (const [childName, child] of root.children.entries()) {
      path.push(childName)
      bt(child, path, counter)
      path.pop()
    }
  }
}

console.log(
  deleteDuplicateFolder([
    ['a'],
    ['a', 'x'],
    ['a', 'x', 'y'],
    ['a', 'z'],
    ['b'],
    ['b', 'x'],
    ['b', 'x', 'y'],
    ['b', 'z'],
    ['b', 'w'],
  ])
)
export {}
