/**
 * @description 仅用于LinkedList内部的结点
 */
class _Node<T = number> {
  value: T
  prev!: _Node<T>
  next!: _Node<T>
  constructor(val: T) {
    this.value = val
  }
}

/**
 * @description 链表实现的双端队列
 */
class LinkedList<T = number> {
  private readonly head: _Node<T>
  private readonly tail: _Node<T>
  length: number

  constructor() {
    this.head = new _Node(undefined as any)
    this.tail = new _Node(undefined as any)
    this.head.next = this.tail
    this.tail.prev = this.head
    this.length = 0
  }

  unshift(val: T): number {
    const node = new _Node(val)
    const next = this.head.next
    this.head.next = node
    node.prev = this.head
    node.next = next
    next.prev = node
    this.length++
    return this.length
  }

  shift(): T | undefined {
    if (this.length > 0) {
      const first = this.head.next
      this.remove(first)
      return first.value
    }
    return undefined
  }

  push(val: T): number {
    const node = new _Node(val)
    const prev = this.tail.prev
    this.tail.prev = node
    node.next = this.tail
    node.prev = prev
    prev.next = node
    this.length++
    return this.length
  }

  pop(): T | undefined {
    if (this.length > 0) {
      const last = this.tail.prev
      this.remove(last)
      return last.value
    }
    return undefined
  }

  first(): T | undefined {
    if (this.length === 0) return undefined
    return this.head.next.value
  }

  last(): T | undefined {
    if (this.length === 0) return undefined
    return this.tail.prev.value
  }

  private remove(node: _Node<T>) {
    const prev = node.prev
    const next = node.next
    prev.next = next
    next.prev = prev
  }
}

export {}

// java查找链表元素：起点折半查找 这样最坏情况也只要找一半就可以了。
// Node<E> node(int index) {
//   // assert isElementIndex(index);

//   if (index < (size >> 1)) {
//       Node<E> x = first;
//       for (int i = 0; i < index; i++)
//           x = x.next;
//       return x;
//   } else {
//       Node<E> x = last;
//       for (int i = size - 1; i > index; i--)
//           x = x.prev;
//       return x;
//   }
// }
