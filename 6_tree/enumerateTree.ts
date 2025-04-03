// !遍历任何抽象树形结构的三个关键：根结点、isLeaf、getChildren

interface ITreeNode<T> {
  isLeaf(): boolean
  getChildren(): ITreeNode<T>[]
}

/**
 * 前序dfs遍历树.
 */
function enumerateTree<T>(root: ITreeNode<T>, f: (node: ITreeNode<T>) => void): void {
  const dfs = (node: ITreeNode<T>) => {
    f(node)
    if (node.isLeaf()) return
    const children = node.getChildren()
    for (let i = 0; i < children.length; i++) {
      const child = children[i]
      dfs(child)
    }
  }

  dfs(root)
}

export {}
