import { MinHeap } from './minheap'

// 暂不支持重复元素 见381. 有重复元素O(1) 时间插入、删除和获取随机元素.ts
class HashHeap<Item = number> extends MinHeap<Item> {
  private hash: Item extends object ? WeakMap<Item, number> : Map<Item, number>

  constructor(
    compareFunction: (a: Item, b: Item) => number = MinHeap.defaultCompareFunction,
    useWeakMap: boolean = false,
    volumn: number = Infinity,
    heap: Item[] = []
  ) {
    super(compareFunction, volumn, heap)
    this.hash = useWeakMap ? (new WeakMap() as any) : new Map()
  }

  override push(val: Item) {
    if (this.heap.length >= this.volumn) {
      this.shift()
    }

    this.heap.push(val)
    this.hash.set(val, this.size - 1)
    this.shiftUp(this.size - 1)

    return this.size
  }

  override shift() {
    const top = this.peek()
    this.remove(top)
    return top
  }

  override swap(parentIndex: number, index: number) {
    if (parentIndex === index) return
    this.hash.set(this.heap[parentIndex], index)
    this.hash.set(this.heap[index], parentIndex)
    ;[this.heap[parentIndex], this.heap[index]] = [this.heap[index], this.heap[parentIndex]]
  }

  remove(val: Item) {
    if (!this.size || this.hash.get(val)) return
    const removeIndex = this.hash.get(val)!
    this.swap(this.size - 1, removeIndex)
    this.heap.pop()
    this.hash.delete(val)
    if (removeIndex < this.size) {
      this.shiftUp(removeIndex)
      this.shiftDown(removeIndex)
    }
  }

  static main() {
    const hashHeap = new HashHeap((a, b) => a - b)
    hashHeap.push(19)
    hashHeap.push(2)
    hashHeap.push(31)
    console.log(hashHeap)
    console.log(hashHeap.shift())
    console.log(hashHeap)
  }
}

// HashHeap.main()
export { HashHeap }
