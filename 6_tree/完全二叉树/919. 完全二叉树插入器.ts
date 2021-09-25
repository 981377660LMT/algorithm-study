import { BinaryTree } from '../力扣加加/Tree'

class CBTInserter {
  private parent: BinaryTree[] // 层序遍历结果 数组下标即节点位置
  private root: BinaryTree | null

  constructor(root: BinaryTree | null) {
    this.parent = []
    this.root = root

    const queue = [root]
    while (queue.length) {
      const head = queue.shift()
      if (head) {
        this.parent.push(head)
        head.left && queue.push(head.left)
        head.right && queue.push(head.right)
      }
    }
  }

  // 辅助列表降低插入复杂度   (index-1)>>1 为父节点索引
  insert(val: number): number {
    const index = this.parent.length - 1
    const newNode = new BinaryTree(val)
    this.parent.push(newNode)

    if (index & 1) {
      this.parent[index >> 1].right = newNode
    } else {
      this.parent[index >> 1].left = newNode
    }

    return this.parent[index >> 1].val
  }

  get_root(): BinaryTree | null {
    return this.root
  }
}

export {}
