// 链表是实现的结构不是抽象的结构，
// 由于缓存极不友好，实际表现比从算法复杂度上得出的感觉差很多，
// 实际应用里面但凡能用连续内存（通过将旧内存复制到新内存来扩容）
// 做的都不会用链表这种实现（除非复制的成本非常非常非常高，而这种情况很不常见）。

class ListNode<T = number> {
  value: T
  prev!: ListNode<T>
  next!: ListNode<T>
  constructor(val: T) {
    this.value = val
  }
}

/**
 * @description 双向链表
 */
class LinkedList<T = number> {
  private head: ListNode<T>
  private tail: ListNode<T>
  length: number

  constructor() {
    this.head = new ListNode(undefined as any)
    this.tail = new ListNode(undefined as any)
    this.head.next = this.tail
    this.tail.prev = this.head
    this.length = 0
  }

  unshift(val: T): number {
    const node = new ListNode(val)
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
    const node = new ListNode(val)
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

  private remove(node: ListNode<T>): void {
    const prev = node.prev
    const next = node.next
    prev.next = next
    next.prev = prev
  }
}

export { LinkedList }

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
