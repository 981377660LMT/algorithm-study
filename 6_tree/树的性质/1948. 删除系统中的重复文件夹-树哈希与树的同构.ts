// 题目要求：如果子树的哈希(结构)相同，则删除

// 0. 用一个多叉树来维护字符串信息
// 1. dfs，子树哈希值记录到 Node 结点中
// 2. dfs，如果子树哈希值重复，就不再继续向下走

class FileTreeNode {
  value: string
  subtreeHash: string
  children: Map<string, FileTreeNode>
  constructor(value: string) {
    this.value = value
    this.subtreeHash = ''
    this.children = new Map()
  }
}

class FileTree {
  root = new FileTreeNode('/')

  insert(path: string[]): void {
    let root = this.root
    for (const char of path) {
      if (!root.children.has(char)) root.children.set(char, new FileTreeNode(char))
      root = root.children.get(char)!
    }
  }
}

function deleteDuplicateFolder(paths: string[][]): string[][] {
  const tree = new FileTree()
  const hashCounter = new Map<string, number>()
  paths.forEach(path => tree.insert(path))

  // 1. 生成子树哈希，并将信息保存在每个节点上
  dfs(tree.root)
  // console.dir(trie, { depth: null })

  // 2. 看子树哈希是否重复，回溯返回结果
  const res: string[][] = []
  bt(tree.root, [])
  return res

  function dfs(root: FileTreeNode): string {
    const subTree: string[] = []
    for (const child of root.children.values()) {
      subTree.push(dfs(child))
    }

    subTree.sort()
    root.subtreeHash = subTree.join('')

    // 叶子结点子树哈希值不计入counter
    if (root.children.size !== 0) {
      hashCounter.set(root.subtreeHash, (hashCounter.get(root.subtreeHash) ?? 0) + 1)
    }

    const res = `${root.value}(${root.subtreeHash})`
    return res
  }

  function bt(root: FileTreeNode, path: string[]): void {
    if (hashCounter.get(root.subtreeHash)! >= 2) return
    if (path.length > 0) res.push(path.slice())

    for (const [childName, child] of root.children.entries()) {
      path.push(childName)
      bt(child, path)
      path.pop()
    }
  }
}

console.log(deleteDuplicateFolder([['a'], ['c'], ['d'], ['a', 'b'], ['c', 'b'], ['d', 'a']]))
export {}
