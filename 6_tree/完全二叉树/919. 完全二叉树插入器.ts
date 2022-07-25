import { BinaryTree as TreeNode } from '../分类/Tree'

// !数组保存树节点 方便获取父结点
class CBTInserter {
  private readonly root: TreeNode | null
  private readonly nodes: TreeNode[] = [new TreeNode(-1)] // !数组保存树节点

  constructor(root: TreeNode | null) {
    this.root = root

    const queue = [root]
    while (queue.length) {
      const node = queue.shift()
      if (node) {
        this.nodes.push(node)
        node.left && queue.push(node.left)
        node.right && queue.push(node.right)
      }
    }
  }

  insert(val: number): number {
    const newNode = new TreeNode(val)
    this.nodes.push(newNode)
    const newIndex = this.nodes.length - 1
    const parent = this.nodes[Math.floor(newIndex / 2)] // index>>1 为父节点索引

    if (newIndex & 1) {
      parent.right = newNode
    } else {
      parent.left = newNode
    }

    return parent.val
  }

  get_root(): TreeNode | null {
    return this.root
  }
}

export {}

if (require.main === module) {
  // ["CBTInserter","insert","get_root"]
  // [[[1]],[2],[]]
  const cbt = new CBTInserter(new TreeNode(1))
  console.log(cbt.insert(2))
  console.log(cbt.get_root())
}
