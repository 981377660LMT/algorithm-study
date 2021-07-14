// 1.插入 insert
// 2.删除堆顶 pop
// 3.获取堆顶 peek
// 4.获取堆大小 size

class MinHeap {
  constructor(private heap: number[], private volumn?: number) {
    this.heap = heap
    this.volumn = volumn
  }

  /**
   *
   * @param val 插入的值
   * @description 将值插入数组(堆)的尾部，然后上移直至父节点不超过它
   * @description 时间复杂度为`O(log(k))`
   */
  insert(val: number) {
    if (this.volumn !== undefined && this.heap.length >= this.volumn) {
      this.pop()
    }

    this.heap.push(val)
    this.shiftUp(this.heap.length - 1)

    return this
  }

  /**
   * @description 用数组尾部元素替换堆顶(直接删除会破坏堆结构),然后下移动直至子节点都大于新堆顶
   * @description 时间复杂度为`O(log(k))`
   */
  pop() {
    this.heap[0] = this.heap.pop()!
    this.shiftDown(0)
    return this
  }

  peek() {
    return this.heap[0]
  }

  size() {
    return this.heap.length
  }

  // 上移
  private shiftUp(index: number) {
    if (index <= 0) return
    const parentIndex = this.getParentIndex(index)
    while (this.heap[parentIndex] > this.heap[index]) {
      this.swap(parentIndex, index)
      this.shiftUp(parentIndex)
    }
  }

  // 下移
  private shiftDown(index: number) {
    const leftChildIndex = this.getLeftChildIndex(index)
    const rightChildIndex = this.getRightChildIndex(index)
    if (this.heap[leftChildIndex] < this.heap[index]) {
      this.swap(leftChildIndex, index)
      this.shiftDown(leftChildIndex)
    }
    if (this.heap[rightChildIndex] < this.heap[index]) {
      this.swap(rightChildIndex, index)
      this.shiftDown(rightChildIndex)
    }
  }

  private getParentIndex(index: number) {
    // 二进制数向右移动一位，相当于Math.floor((index-1)/2)
    return (index - 1) >> 1
  }

  private getLeftChildIndex(index: number) {
    return index * 2 + 1
  }

  private getRightChildIndex(index: number) {
    return index * 2 + 2
  }

  private swap(parentIndex: number, index: number) {
    ;[this.heap[parentIndex], this.heap[index]] = [this.heap[index], this.heap[parentIndex]]
  }
}

if (require.main === module) {
  const minHeap = new MinHeap([10, 20, 30])
  console.log(minHeap)
  minHeap.insert(4)
  console.log(minHeap)
  minHeap.pop()
  console.log(minHeap)
}

export { MinHeap }
