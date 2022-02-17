// 注意堆的索引从0开始，而线段树的索引从1开始
// 堆:root (root<<1)+1 (root<<1)+2
// 线段树: root root<<1 root<<1|1

import assert from 'assert'

type CompareFunction<T> = (a: T, b: T) => number

class MinHeap<HeapValue = number> {
  protected readonly heap: HeapValue[]
  protected readonly capacity: number
  protected readonly compare: CompareFunction<HeapValue>
  static readonly defaultCompareFunction = (a: any, b: any) => a - b

  constructor(
    compareFunction: CompareFunction<HeapValue> = MinHeap.defaultCompareFunction,
    capacity: number = Infinity,
    heap: HeapValue[] = []
  ) {
    this.compare = compareFunction
    this.capacity = capacity
    this.heap = heap
  }

  /**
   *
   * @param value 插入的值
   * @description 将值插入数组(堆)的尾部，然后上移直至父节点不超过它
   * @description 时间复杂度为`O(log(h))`
   */
  heappush(value: HeapValue): void {
    if (this.heap.length >= this.capacity) {
      this.heappop()
    }

    this.heap.push(value)
    this.pushUp(this.heap.length - 1)
  }

  /**
   * @description 用数组尾部元素替换堆顶(直接删除会破坏堆结构),然后下移动直至子节点都大于新堆顶
   * @description 时间复杂度为`O(log(h))`
   */
  heappop(): HeapValue | undefined {
    if (this.heap.length === 0) {
      return undefined
    } else if (this.heap.length === 1) {
      return this.heap.pop()!
    }

    const top = this.peek()
    const last = this.heap.pop()!
    this.heap[0] = last
    this.pushDown(0)
    return top
  }

  peek(): HeapValue | undefined {
    return this.heap[0]
  }

  get size(): number {
    return this.heap.length
  }

  /**
   * 取出堆顶元素，替换成val;
   * 一次O(log(h)的操作)
   */
  // heapreplace(val: Item) {
  //   this.heap[0] = val
  //   this.shiftDown(0)
  // }

  /**
   *
   * @description 将非叶子节点(2^(h-1)-1个，约n/2) 倒序shiftdown
   * @description 堆化的复杂度是O(n)
   */
  heapify(): void {
    if (this.heap.length <= 1) return
    const last = this.heap.length - 1
    const lastParent = (last - 1) >> 1
    for (let i = lastParent; ~i; i--) {
      this.pushDown(i)
    }
  }

  protected pushUp(root: number): void {
    let parent = (root - 1) >> 1
    while (parent >= 0 && this.compare(this.heap[parent], this.heap[root]) > 0) {
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

      if (left < this.heap.length && this.compare(this.heap[left], this.heap[minIndex]) < 0) {
        minIndex = left
      }

      if (right < this.heap.length && this.compare(this.heap[right], this.heap[minIndex]) < 0) {
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
}
