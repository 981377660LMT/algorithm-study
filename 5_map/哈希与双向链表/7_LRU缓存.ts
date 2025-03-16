// O(1) 时间复杂度内完成get与put (js的map是有序字典)
// map插入顺序 默认第一个即最早插入的值
// LRU:获取值时将值重新放到最后，设置值时删除原来值(有的话)再放到最后，超出容量则删除集合的头部元素
// this.cache.delete(this.cache.keys().next().value)
class LRU {
  private readonly capacity: number
  private readonly cache: Map<number, number> = new Map()

  constructor(capacity: number) {
    this.capacity = capacity
  }

  /**
   * @param key 如果关键字 key 存在于缓存中，则返回关键字的值，否则返回 -1 。
   */
  get(key: number): number {
    const res = this.cache.get(key)!
    if (res === undefined) return -1
    // 删除再插入保证顺序
    this.cache.delete(key)
    this.cache.set(key, res)
    return res
  }

  /**
   * @description 如果关键字已经存在，则变更其数据值；如果关键字不存在，
   * 则插入该组「关键字-值」。当缓存容量达到上限时，
   * 它应该在写入新数据之前删除最久未使用的数据值，从而为新的数据值留出空间。
   */
  put(key: number, value: number): void {
    if (this.cache.has(key)) {
      this.cache.delete(key)
    }
    this.cache.set(key, value)
    if (this.cache.size > this.capacity) {
      // 删除最久没有使用的元素
      this.cache.delete(this.cache.keys().next().value!)
    }
  }
}

if (require.main === module) {
  const cache = new LRU(2)

  cache.put(1, 1)
  console.log(cache.get(1))
  cache.put(2, 2)
  cache.put(3, 3)
  cache.put(4, 4)
  console.log(cache.get(1))
  console.log(cache.get(3))
  cache.put(2, 2)
  console.log(cache.get(3))
}

export {}
