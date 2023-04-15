/* eslint-disable no-param-reassign */

class LinkedListNode<E> {
  value: E | undefined
  pre: LinkedListNode<E> | undefined
  next: LinkedListNode<E> | undefined

  constructor(value?: E, left?: LinkedListNode<E>, right?: LinkedListNode<E>) {
    this.value = value
    this.pre = left
    this.next = right
  }

  /**
   * @param node 在当前node之后插入新节点 并返回新节点
   */
  insertAfter(node: LinkedListNode<E>): LinkedListNode<E> {
    node.pre = this
    node.next = this.next
    node.pre.next = node
    if (node.next) node.next.pre = node
    return node
  }

  /**
   * @param node 在当前node之前插入新节点 并返回新节点
   */
  insertBefore(node: LinkedListNode<E>): LinkedListNode<E> {
    node.next = this
    node.pre = this.pre
    node.next.pre = node
    if (node.pre) node.pre.next = node
    return node
  }

  /**
   * @description 从链表里移除自身
   */
  remove(): LinkedListNode<E> {
    if (this.pre) this.pre.next = this.next
    if (this.next) this.next.pre = this.pre
    this.pre = void 0
    this.next = void 0
    return this
  }

  toString(): string {
    return `${this.value}->${this.next}`
  }
}

if (require.main === module) {
  const node = new LinkedListNode(
    1,
    new LinkedListNode(2, new LinkedListNode(3)),
    new LinkedListNode(4, undefined, new LinkedListNode(5))
  )
  console.log(node.toString())
}

export { LinkedListNode }
