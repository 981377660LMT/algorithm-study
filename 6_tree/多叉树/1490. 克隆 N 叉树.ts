// 1490. 克隆 N 叉树
// https://leetcode.cn/problems/clone-n-ary-tree/description/

// data
interface ITreeNode {
  children: this[]
}

function cloneTree<T extends ITreeNode>(
  root: T,
  operations: { getChildren: (node: T) => T[]; cloneNode: (node: T) => T }
): T {
  const { getChildren, cloneNode } = operations
  const dfs = (node: T): T => {
    const newNode = cloneNode(node)
    const children = getChildren(node)
    for (const child of children) {
      newNode.children.push(dfs(child))
    }
    return newNode
  }
  return dfs(root)
}

export { cloneTree }

if (require.main === module) {
  class Node {
    val: number
    children: Node[]

    constructor(val = 0, children: Node[] = []) {
      this.val = val
      this.children = children
    }
  }

  // eslint-disable-next-line no-inner-declarations
  function cloneTree_(root: Node | null): Node | null {
    if (!root) return null
    return cloneTree<Node>(root, {
      getChildren: node => node.children,
      cloneNode: node => new Node(node.val)
    })
  }
}
