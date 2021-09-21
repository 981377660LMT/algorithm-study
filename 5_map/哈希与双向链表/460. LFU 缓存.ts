class ListNode<K, V> {
  pre!: ListNode<K, V>
  next!: ListNode<K, V>
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
  head: ListNode<K, V>
  tail: ListNode<K, V>
  size: number

  constructor() {
    this.head = new ListNode(undefined as any, undefined as any)
    this.tail = new ListNode(undefined as any, undefined as any)
    this.head.next = this.tail
    this.tail.pre = this.head
    this.size = 0
  }

  unshift(node: ListNode<K, V>) {
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

  remove(node: ListNode<K, V>) {
    const pre = node.pre
    const next = node.next
    pre.next = next
    next.pre = pre
    this.size--
  }
}

// Least Frequently Used
// key->Node->Node.freq->freqToDoubleLinkedList
class LFUCache {
  private capacity: number
  private minFreq: number
  private keyToNode: Map<number, ListNode<number, number>>
  private freqToDoubleLinkedList: Map<number, DoubleLinkedList<number, number>>

  constructor(capacity: number) {
    this.capacity = capacity
    this.minFreq = 0
    this.keyToNode = new Map()
    // 每个DoubleLinkedList都是LRU缓存
    this.freqToDoubleLinkedList = new Map()
  }

  /**
   *
   * @param key 如果键存在于缓存中，则获取键的值（总是正数），否则返回 -1
   */
  get(key: number): number {
    if (!this.keyToNode.has(key)) return -1
    const node = this.keyToNode.get(key)!
    this.upgrade(node)
    return node.value
  }

  /**
   *
   * @param key 如果键已存在，则变更其值；如果键不存在，请插入键值对
   * @param value
   * 当缓存达到其容量时，则应该在插入新项之前，使最不经常使用的项无效。
   * 在此问题中，当存在平局（即两个或更多个键具有相同使用频率）时，应该去除最久未使用的键。
   * 「项的使用次数」就是自插入该项以来对其调用 get 和 put 函数的次数之和。使用次数会在对应项被移除后置为 0 。
   * @summary
   * 1.由于频率相同时，需要删除最久未使用的键，我们每个频率都要维护head和tail 删除靠近尾部的节点(lru)
   * 2.get/put操作会使节点升级到新的DoubleLinkedList
   */
  put(key: number, value: number): void {
    if (this.capacity === 0) return
    if (this.keyToNode.has(key)) {
      const node = this.keyToNode.get(key)!
      node.value = value
      this.upgrade(node)
    } else {
      const node = new ListNode(key, value)
      this.keyToNode.set(key, node)
      // 需要删除
      if (this.keyToNode.size > this.capacity) {
        const minFreqList = this.freqToDoubleLinkedList.get(this.minFreq)!
        this.keyToNode.delete(minFreqList.pop()!.key)
      }
      this.minFreq = 1
      !this.freqToDoubleLinkedList.has(node.freq) &&
        this.freqToDoubleLinkedList.set(node.freq, new DoubleLinkedList())
      const newList = this.freqToDoubleLinkedList.get(node.freq)!
      newList.unshift(node)
    }
  }

  /**
   *
   * @param node 频率增加，升级到新的DoubleLinkedList
   */
  private upgrade(node: ListNode<number, number>) {
    const oldList = this.freqToDoubleLinkedList.get(node.freq)
    oldList?.remove(node)
    if (node.freq === this.minFreq && oldList?.size === 0) this.minFreq++
    node.freq++
    !this.freqToDoubleLinkedList.has(node.freq) &&
      this.freqToDoubleLinkedList.set(node.freq, new DoubleLinkedList())
    const newList = this.freqToDoubleLinkedList.get(node.freq)!
    newList.unshift(node)
  }

  static main() {
    const lfuCache = new LFUCache(2)
    console.log(lfuCache.put(1, 1))
    console.log(lfuCache.put(2, 2))
    console.log(lfuCache.get(1)) // 返回 1
    console.log(lfuCache.put(3, 3)) // 去除 key 2
    console.log(lfuCache.get(2)) // 返回 -1 (未找到key 2)
    console.log(lfuCache.get(3)) // 返回 3
    console.log(lfuCache.put(4, 4)) // 去除 key 1
    console.log(lfuCache.get(1)) // 返回 -1 (未找到 key 1)
    console.log(lfuCache.get(3)) // 返回 3
    console.log(lfuCache.get(4)) // 返回 4
  }
}

LFUCache.main()

export {}
