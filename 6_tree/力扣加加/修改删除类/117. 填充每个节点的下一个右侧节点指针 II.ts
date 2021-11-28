import { NodeWithNext } from './116. 填充每个节点的下一个右侧节点指针常量空间'

class Node {
  val: number
  left: Node | null
  right: Node | null
  next: Node | null
  constructor(val = 0) {
    this.val = val
    this.left = null
    this.right = null
    this.next = null
  }
}
/**
 *
 * @param root 二叉树 并不一定完全二叉树
 * 利用next指针，将每层节点链接成链表进行遍历
 */
function connect(root: Node | null): Node | null {
  if (!root) return null
  let rootP: Node | null = root

  while (rootP) {
    const dummy = new Node(-1) // 为下一行的之前的节点，相当于下一行所有节点链表的头结点；
    let pre = dummy

    while (rootP) {
      if (rootP.left) {
        // 链接下一行的节点
        pre.next = rootP.left
        pre = pre.next
      }
      if (rootP.right) {
        pre.next = rootP.right
        pre = pre.next
      }
      rootP = rootP.next
    }

    rootP = dummy.next // 此处为换行操作，更新到下一行
  }

  return root
}

export {}
