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

  override heappush(val: Item) {
    if (this.heap.length >= this.capacity) {
      this.heappop()
    }

    this.heap.push(val)
    this.hash.set(val, this.size - 1)
    this.pushUp(this.size - 1)

    return this.size
  }

  override heappop() {
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
    if (!this.size || !this.hash.has(val)) return
    const removeIndex = this.hash.get(val)!
    this.swap(this.size - 1, removeIndex)
    this.heap.pop()
    this.hash.delete(val)
    if (removeIndex < this.size) {
      this.pushUp(removeIndex)
      this.pushDown(removeIndex)
    }
  }

  static main() {
    const hashHeap = new HashHeap<number[]>((a, b) => a[1] - a[0] - (b[1] - b[0]))
    const intervalA = [1, 2]
    const intervalB = [3, 8]
    const intervalC = [7, 9]
    hashHeap.heappush(intervalA)
    hashHeap.heappush(intervalB)
    hashHeap.heappush(intervalC)
    hashHeap.remove(intervalB)
    hashHeap.remove(intervalB)
    hashHeap.remove(intervalB)
    console.log(hashHeap.heappop())
    console.log(hashHeap)
  }
}

if (require.main === module) {
  HashHeap.main()
}
export { HashHeap }
