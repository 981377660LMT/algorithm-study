class LinkedListNode<V = unknown> {
  value: V
  left: LinkedListNode<V> | undefined
  right: LinkedListNode<V> | undefined

  constructor(value: V, left?: LinkedListNode<V>, right?: LinkedListNode<V>) {
    this.value = value
    this.left = left
    this.right = right
  }

  /**
   * @param node 在当前node之后插入新节点 并返回新节点
   */
  insertRight(node: LinkedListNode<V>): LinkedListNode<V> {
    node.left = this
    node.right = this.right
    node.left.right = node
    if (node.right) node.right.left = node
    return node
  }

  /**
   * @param node 在当前node之前插入新节点 并返回新节点
   */
  insertLeft(node: LinkedListNode<V>): LinkedListNode<V> {
    node.right = this
    node.left = this.left
    node.right.left = node
    if (node.left) node.left.right = node
    return node
  }

  /**
   * @description 从链表里移除自身
   */
  remove(): void {
    if (this.left) this.left.right = this.right
    if (this.right) this.right.left = this.left
  }
}

export { LinkedListNode }
