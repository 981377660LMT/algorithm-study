import { BinaryTree } from '../力扣加加/Tree'

// 此题类似于树上倍增那题
class CBTInserter {
  private nodes: BinaryTree[] // 层序遍历结果 数组下标即节点位置
  private root: BinaryTree | null

  constructor(root: BinaryTree | null) {
    this.nodes = []
    this.root = root

    const queue = [root]
    while (queue.length) {
      const head = queue.shift()
      if (head) {
        this.nodes.push(head)
        head.left && queue.push(head.left)
        head.right && queue.push(head.right)
      }
    }
  }

  // 辅助列表降低插入复杂度
  insert(val: number): number {
    const newNode = new BinaryTree(val)
    this.nodes.push(newNode)
    const index = this.nodes.length - 1
    const parent = this.nodes[(index - 1) >> 1] // (index-1)>>1 为父节点索引

    if (index & 1) {
      parent.left = newNode
    } else {
      parent.right = newNode
    }

    return parent.val
  }

  get_root(): BinaryTree | null {
    return this.root
  }
}

export {}
