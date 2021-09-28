import { BinaryTree } from '../../6_tree/力扣加加/Tree'

// 设计迭代器的时候提前把所有的值遍历并且保存起来的做法并不好
// 把递归转成迭代，基本想法就是用栈:迭代时计算 next  节点
class BSTIterator {
  // 一个「单调栈」
  // 空间复杂度：O(h)，h 为树的高度
  private monoStack: BinaryTree[]

  /**
   *
   * @param root 一路到底，把根节点和它的所有左节点放到栈中；
   */
  constructor(root: BinaryTree | null) {
    this.monoStack = []
    while (root) {
      this.monoStack.push(root)
      root = root.left
    }
  }

  /**
   * 弹出栈顶的节点；
     如果它有右子树，则对右子树一路到底，把它和它的所有左节点放到栈中。
   */
  next(): number {
    const res = this.monoStack.pop()!

    let right = res.right
    while (right) {
      this.monoStack.push(right)
      right = right.left
    }

    return res.val
  }

  hasNext(): boolean {
    return this.monoStack.length > 0
  }
}

export {}
