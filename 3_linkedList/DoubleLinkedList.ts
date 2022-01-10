class Node<K, V> {
  pre!: Node<K, V>
  next!: Node<K, V>
  key: K
  value: V
  freq: number
  constructor(key: K, value: V) {
    this.key = key
    this.value = value
    this.freq = 1
  }
}

class DoubleLinkedList<K, V> {
  head: Node<K, V>
  tail: Node<K, V>
  size: number

  constructor() {
    this.head = new Node(undefined as any, undefined as any)
    this.tail = new Node(undefined as any, undefined as any)
    this.head.next = this.tail
    this.tail.pre = this.head
    this.size = 0
  }

  unshift(node: Node<K, V>) {
    const next = this.head.next
    this.head.next = node
    node.pre = this.head
    node.next = next
    next.pre = node
    this.size++
  }

  pop() {
    if (this.size > 0) {
      const rear = this.tail.pre
      this.remove(rear)
      return rear
    }
    return undefined
  }

  remove(node: Node<K, V>) {
    const pre = node.pre
    const next = node.next
    pre.next = next
    next.pre = pre
  }
}

export {}
