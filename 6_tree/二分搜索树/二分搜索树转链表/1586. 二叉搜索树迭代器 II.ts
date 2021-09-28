import { BinaryTree } from '../../力扣加加/Tree'

//
class BSTIterator {
  private stack: BinaryTree[]
  private dummy: BinaryTree
  private cur: BinaryTree | null | undefined

  // 使用树中不存在的最小值来初始化内部指针
  constructor(private root: BinaryTree | null) {
    this.stack = []
    this.dummy = new BinaryTree(Infinity)
    this.cur = this.dummy
    this.pushLeft(root)
  }

  // 游标后面有（中序遍历序列的）节点，或者栈非空（游标节点有父节点）
  hasNext(): boolean {
    return this.cur?.right != null || this.stack.length > 0
  }

  // 从栈中取下一个节点 并对next的right节点执行pushLeft操作 然后删除next的right节点
  next(): number {
    if (this.cur?.right == null) {
      const next = this.stack.pop()!
      this.pushLeft(next.right)
      next.right = null // 这句很关键

      // next和cur串在一起  类似中序遍历转链表常用的pre 和 root
      next.left = this.cur!
      this.cur!.right = next
    }
    this.cur = this.cur?.right
    return this.cur!.val
  }

  hasPrev(): boolean {
    return this.cur !== this.dummy && this.cur?.left !== this.dummy
  }

  prev(): number {
    this.cur = this.cur?.left
    return this.cur!.val
  }

  private pushLeft(root: BinaryTree | null) {
    while (root) {
      this.stack.push(root)
      root = root.left
    }
  }
}

export {}
// 你可以假设 next() 和 prev() 的调用总是有效的
