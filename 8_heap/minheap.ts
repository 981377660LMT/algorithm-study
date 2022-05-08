// 注意堆的索引从0开始，而线段树的索引从1开始
// 堆:root (root<<1)+1 (root<<1)+2
// 线段树: root root<<1 root<<1|1

import assert from 'assert'

type Comparator<T> = (a: T, b: T) => number

class MinHeap<HeapValue = number> {
  protected readonly heap: HeapValue[]
  protected readonly capacity: number
  protected readonly comparator: Comparator<HeapValue>
  static readonly defaultComparator = (a: any, b: any) => a - b

  constructor(
    comparator: Comparator<HeapValue> = MinHeap.defaultComparator,
    capacity: number = Infinity,
    heap: HeapValue[] = []
  ) {
    this.comparator = comparator
    this.capacity = capacity
    this.heap = heap
  }

  get size(): number {
    return this.heap.length
  }

  /**
   *
   * @param value 插入的值
   * @description 将值插入数组(堆)的尾部，然后上移直至父节点不超过它
   *
   * 时间复杂度为`O(log(h))`
   */
  heappush(value: HeapValue): void {
    this.heap.push(value)
    this.pushUp(this.heap.length - 1)
    if (this.heap.length > this.capacity) {
      this.heappop()
    }
  }

  /**
   * @description 用数组尾部元素替换堆顶(直接删除会破坏堆结构),然后下移动直至子节点都大于新堆顶
   *
   * 时间复杂度为`O(log(h))`
   */
  heappop(): HeapValue | undefined {
    if (this.heap.length === 0) {
      return undefined
    } else if (this.heap.length === 1) {
      return this.heap.pop()!
    }

    const returnItem = this.heap[0]
    const last = this.heap.pop()!
    this.heap[0] = last
    this.pushDown(0)
    return returnItem
  }

  /**
   *
   * @description 将非叶子节点(2^(h-1)-1个，约n/2) 倒序shiftdown
   *
   * 堆化的复杂度是O(n)
   */
  heapify(): void {
    if (this.heap.length <= 1) return
    const last = this.heap.length - 1
    const lastParent = (last - 1) >> 1
    for (let i = lastParent; ~i; i--) {
      this.pushDown(i)
    }
  }

  /**
   * @description `入堆+出堆`的更快的版本
   */
  heappushpop(value: HeapValue): HeapValue | undefined {
    if (this.heap.length > 0 && this.comparator(this.heap[0], value) < 0) {
      ;[value, this.heap[0]] = [this.heap[0], value]
      this.pushDown(0)
    }

    return value
  }

  /**
   * @description `出堆+入堆`的更快的版本
   */
  heapreplace(value: HeapValue): HeapValue | undefined {
    const returnItem = this.heap[0]
    this.heap[0] = value
    this.pushDown(0)
    return returnItem
  }

  peek(): HeapValue | undefined {
    return this.heap[0]
  }

  protected pushUp(root: number): void {
    let parent = (root - 1) >> 1
    while (parent >= 0 && this.comparator(this.heap[parent], this.heap[root]) > 0) {
      this.swap(parent, root)
      root = parent
      parent = (parent - 1) >> 1
    }
  }

  protected pushDown(root: number): void {
    // 还有孩子，即不是叶子节点
    while ((root << 1) + 1 < this.heap.length) {
      const left = (root << 1) + 1
      const right = (root << 1) + 2

      let minIndex = root

      if (left < this.heap.length && this.comparator(this.heap[left], this.heap[minIndex]) < 0) {
        minIndex = left
      }

      if (right < this.heap.length && this.comparator(this.heap[right], this.heap[minIndex]) < 0) {
        minIndex = right
      }

      if (minIndex === root) return

      this.swap(root, minIndex)
      root = minIndex
    }
  }

  protected swap(index1: number, index2: number): void {
    ;[this.heap[index1], this.heap[index2]] = [this.heap[index2], this.heap[index1]]
  }

  // private getParentIndex(index: number) {
  //   // 减一后二进制数向右移动一位，相当于Math.floor((index-1)/2)
  //   return (index - 1) >> 1
  // }

  // private getLeftChildIndex(index: number) {
  //   return index * 2 + 1
  // }

  // private getRightChildIndex(index: number) {
  //   return index * 2 + 2
  // }
}

export { MinHeap }

if (require.main === module) {
  const heap = new MinHeap()
  heap.heappush(1)
  heap.heappush(8)
  heap.heappush(3)
  heap.heappush(5)
  assert.strictEqual(heap.heappop(), 1)
  assert.strictEqual(heap.heappop(), 3)
  assert.strictEqual(heap.heappop(), 5)
  assert.strictEqual(heap.heappop(), 8)
  assert.strictEqual(heap.heappop(), undefined)
  heap.heappush(2)
  assert.strictEqual(heap.heapreplace(3), 2)
  assert.strictEqual(heap.peek(), 3)
  assert.strictEqual(heap.heappushpop(1), 1)
  assert.strictEqual(heap.peek(), 3)
}
