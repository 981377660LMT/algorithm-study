import assert from 'assert'

class ListNode<K, V> {
  key: K
  val: V
  pre!: ListNode<K, V>
  next!: ListNode<K, V>

  constructor(key: K, val: V) {
    this.key = key
    this.val = val
  }
}

/**
 * @summary 本质上是 散列表+双端链表
 * @description
 * 散列表O(1) 根据key找到对应的node的val
 * 双端链表适合遍历keys
 */
class LRU<K extends PropertyKey, V> {
  private size: number
  private capacity: number
  private data: Record<K, ListNode<K, V>>
  // head与tail不用来存储值
  private head: ListNode<K, V>
  private tail: ListNode<K, V>

  constructor(capacity: number) {
    this.size = 0
    this.capacity = capacity
    this.data = Object.create(null)
    // @ts-ignore
    this.head = new ListNode<K, V>()
    // @ts-ignore
    this.tail = new ListNode<K, V>()
    this.head.next = this.tail
    this.tail.pre = this.head
  }

  /**
   *
   * @param key 如果关键字 key 存在于缓存中，则返回关键字的值，否则返回 -1 。
   */
  get(key: K) {
    if (this.data[key] !== undefined) {
      const node = this.data[key]
      // console.log(node)
      this.removeNode(node)
      this.appendHead(node)
      return node.val
    } else {
      return -1
    }
  }

  /**
   * @description 如果关键字已经存在，则变更其数据值；如果关键字不存在，
   * 则插入该组「关键字-值」。当缓存容量达到上限时，
   * 它应该在写入新数据之前删除最久未使用的数据值，从而为新的数据值留出空间。
   * 最终将这个node 插入到链表头部
   */
  put(key: K, value: V) {
    let node: ListNode<K, V>

    if (this.data[key] !== undefined) {
      node = this.data[key]
      this.removeNode(node)
      node.val = value
    } else {
      node = new ListNode(key, value)
      this.data[key] = node
      if (this.size < this.capacity) {
        this.size++
      } else {
        const key = this.removeTail()
        delete this.data[key]
      }
    }

    this.appendHead(node)
  }

  private removeNode(node: ListNode<K, V>) {
    const preNode = node.pre
    const nextNode = node.next
    preNode.next = nextNode
    nextNode.pre = preNode
  }

  private appendHead(node: ListNode<K, V>) {
    const firstNode = this.head.next
    this.head.next = node
    node.next = firstNode
    node.pre = this.head
    firstNode.pre = node
  }

  private removeTail() {
    const key = this.tail.pre.key
    this.removeNode(this.tail.pre)
    return key
  }
}

if (require.main === module) {
  const cache = new LRU(2)
  cache.put(1, 1)
  assert.strictEqual(cache.get(1), 1)
  cache.put(2, 2)
  cache.put(3, 3)
  cache.put(4, 4)
  assert.strictEqual(cache.get(1), -1)
  assert.strictEqual(cache.get(3), 3)
  cache.put(2, 2)
  assert.strictEqual(cache.get(3), 3)
}

export {}
