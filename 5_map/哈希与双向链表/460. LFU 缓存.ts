import { LinkedListNode } from '../../3_linkedList/LinkedListNode'

interface NodeValue<K extends PropertyKey, V = unknown> {
  key: K
  value: V
  freq: number
}

// !启示:副作用很多且影响范围很多的逻辑 需要把类似的操作封装成单独的函数

// Least Frequently Used
// key->Node->Node.freq->freqToDoubleLinkedList
class LFUCache<K extends PropertyKey, V = unknown> {
  readonly capacity: number
  private _minFreq = 1
  private readonly _nodeMap: Map<K, LinkedListNode<NodeValue<K, V>>> = new Map()

  /**
   * @description 每个LinkedList都是LRU缓存
   */
  private readonly _freqMap: Map<number, LinkedListNode<NodeValue<K, V>>> = new Map()
  private static _createList(): LinkedListNode<any> {
    const node = new LinkedListNode(void 0)
    node.next = node
    node.pre = node
    return node
  }

  constructor(capacity: number) {
    this.capacity = capacity
  }

  /**
   * @param key 如果键存在于缓存中，则获取键的值（总是正数），否则返回 -1
   */
  get(key: K): V | -1 {
    if (!this._nodeMap.has(key)) return -1
    const node = this._nodeMap.get(key)!
    this._increase(node)
    return node.value.value
  }

  /**
   * @param key 如果键已存在，则变更其值；如果键不存在，请插入键值对
   * @summary
   * 1.由于频率相同时，需要删除最久未使用的键，我们每个频率都要维护一个LRU
   * 2.get/put操作会使节点升级到新的LinkedList
   *
   * !3.应该在`写入新数据之前`删除最久未使用的数据值，从而为新的数据值留出空间。
   */
  put(key: K, value: V): void {
    if (this.capacity === 0) return
    let newNode: LinkedListNode<NodeValue<K, V>>

    if (this._nodeMap.has(key)) {
      newNode = this._nodeMap.get(key)!
      newNode.value.value = value
    } else {
      if (this._nodeMap.size === this.capacity) this._removeLFU()
      newNode = new LinkedListNode({ key, value, freq: 0 })
      this._nodeMap.set(key, newNode)
    }

    this._increase(newNode)
  }

  /**
   * @param node 频率增加
   * @description
   * 1. 频率增加
   * 2. freqMap移除
   * 3. freqMap加入
   * 4. 维护minFreq
   */
  private _increase(node: LinkedListNode<NodeValue<K, V>>) {
    const preFreq = node.value.freq
    const newFreq = preFreq + 1
    node.value.freq++

    node.remove()

    !this._freqMap.has(newFreq) && this._freqMap.set(newFreq, LFUCache._createList())
    const newRoot = this._freqMap.get(newFreq)!
    newRoot.insertPre(node)

    if (newFreq === 1) {
      this._minFreq = 1
    } else if (this._minFreq === preFreq) {
      const preRoot = this._freqMap.get(preFreq)!
      if (preRoot.next === preRoot) this._minFreq = newFreq
    }
  }

  /**
   * 从LFU中移除 {@link LinkedListNode}
   * 1. nodeMap移除
   * 2. freqMap移除
   * 3. 维护minFreq(移除LFU只发生在加入新节点时,此时minFreq为1)
   */
  private _removeLFU() {
    const toRemove = this._freqMap.get(this._minFreq)!.next!

    this._nodeMap.delete(toRemove.value.key)

    toRemove.remove()

    this._minFreq = 1
  }
}

if (require.main === module) {
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

export {}
