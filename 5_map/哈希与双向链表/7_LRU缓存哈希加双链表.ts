import assert from 'assert'
import { LinkedListNode } from '../../3_linkedList/LinkedListNode'

interface NodeValue<K extends PropertyKey, V = unknown> {
  key: K
  value: V
}

/**
 * @summary 本质上是 散列表+双端链表
 * @description
 * 根据题目要求,存储的数据需要保证顺序关系(逻辑层面) ===> 使用数组,链表等保证顺序关系
   同时需要对数据进行频繁的增删, 时间复杂度 O(1) ==> 使用链表等
   对数据进行读取时, 时间复杂度 O(1) ===> 使用哈希表 最终采取双向链表 + 哈希表
   双向链表按最后一次访问的时间的顺序进行排列, 链表`尾部`为最近访问的节点
   哈希表,以关键字为键,以链表节点的地址为值
 */
class LRUCache<K extends PropertyKey, V = unknown> {
  size = 0
  readonly capacity: number
  private readonly linkedMap: Record<K, LinkedListNode<NodeValue<K, V>>>
  private readonly root: LinkedListNode<NodeValue<K, V>>

  constructor(capacity: number) {
    this.capacity = capacity
    this.linkedMap = Object.create(null)
    // @ts-ignore
    this.root = new LinkedListNode<NodeValue<K, V>>()
    this.root.next = this.root
    this.root.pre = this.root
  }

  /**
   * @param key 如果关键字 key 存在于缓存中，则返回关键字的值，否则返回 -1 。
   */
  get(key: K): V | -1 {
    if (this.linkedMap[key] != void 0) {
      const node = this.linkedMap[key]
      node.remove()
      this.root.insertBefore(node)
      return node.value.value
    }

    return -1
  }

  /**
   * @description 如果关键字已经存在，则变更其数据值；如果关键字不存在，
   * 则插入该组「关键字-值」。当缓存容量达到上限时，
   *
   * !它应该在`写入新数据之前`删除最久未使用的数据值，从而为新的数据值留出空间。
   */
  put(key: K, value: V): void {
    let newNode: LinkedListNode<NodeValue<K, V>>
    if (this.linkedMap[key] != void 0) {
      newNode = this.linkedMap[key].remove()
      newNode.value.value = value
    } else {
      if (this.size === this.capacity) {
        const deleted = this.root.next?.remove()
        delete this.linkedMap[deleted!.value.key]
        this.size--
      }
      newNode = new LinkedListNode<NodeValue<K, V>>({ key, value })
      this.linkedMap[key] = newNode
      this.size++
    }

    this.root.insertBefore(newNode)
  }
}

if (require.main === module) {
  const cache = new LRUCache(2)
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
