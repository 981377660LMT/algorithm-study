class LinkedListNode<V = unknown> {
  value: V
  pre: LinkedListNode<V> | undefined
  next: LinkedListNode<V> | undefined

  constructor(value: V, left?: LinkedListNode<V>, right?: LinkedListNode<V>) {
    this.value = value
    this.pre = left
    this.next = right
  }

  /**
   * @param node 在当前node之后插入新节点 并返回新节点
   */
  insertRight(node: LinkedListNode<V>): LinkedListNode<V> {
    node.pre = this
    node.next = this.next
    node.pre.next = node
    if (node.next) node.next.pre = node
    return node
  }

  /**
   * @param node 在当前node之前插入新节点 并返回新节点
   */
  insertLeft(node: LinkedListNode<V>): LinkedListNode<V> {
    node.next = this
    node.pre = this.pre
    node.next.pre = node
    if (node.pre) node.pre.next = node
    return node
  }

  /**
   * @description 从链表里移除自身
   */
  remove(): LinkedListNode<V> {
    if (this.pre) this.pre.next = this.next
    if (this.next) this.next.pre = this.pre
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
