// 1.插入 push
// 2.删除堆顶 shift
// 3.获取堆顶 peek
// 4.获取堆大小 size

class MinHeap<Item> {
  private heap: Item[]
  private volumn: number
  private compareFunction: (a: Item, b: Item) => number
  private static defaultCompareFunction = (a: any, b: any) => a - b

  constructor(compareFunction: (a: Item, b: Item) => number = MinHeap.defaultCompareFunction) {
    this.compareFunction = compareFunction
    this.heap = []
    this.volumn = Infinity
  }

  /**
   *
   * @param val 插入的值
   * @description 将值插入数组(堆)的尾部，然后上移直至父节点不超过它
   * @description 时间复杂度为`O(log(h))`
   */
  push(val: Item) {
    if (this.heap.length >= this.volumn) {
      this.pop()
    }

    this.heap.push(val)
    this.shiftUp(this.heap.length - 1)

    return this.size
  }

  /**
   * @description 用数组尾部元素替换堆顶(直接删除会破坏堆结构),然后下移动直至子节点都大于新堆顶
   * @description 时间复杂度为`O(log(h))`
   */
  pop() {
    const top = this.peek()
    const last = this.heap.pop()!
    if (this.size !== 0) {
      this.heap[0] = last
      this.shiftDown(0)
    }

    return top
  }

  peek() {
    return this.heap[0]
  }

  get size() {
    return this.heap.length
  }

  /**
   * 取出堆顶元素，替换成val;
   * 一次O(log(h)的操作)
   */
  replace(val: Item) {
    this.heap[0] = val
    this.shiftDown(0)
    return this
  }

  /**
   *
   * @description 将非叶子节点(2^(h-1)-1个，约n/2) 倒序shiftdown
   * @description 堆化的复杂度是O(n)
   */
  heapify() {
    const start = this.getParentIndex(this.size - 1)
    for (let i = start; i >= 0; i--) {
      this.shiftDown(i)
    }
  }

  /**
   *
   * @param index 数组中的index
   * @returns
   */
  private shiftUp(index: number) {
    if (index <= 0) return
    const parentIndex = this.getParentIndex(index)
    while (this.compareFunction(this.heap[parentIndex], this.heap[index]) > 0) {
      this.swap(parentIndex, index)
      this.shiftUp(parentIndex)
    }
  }

  // 下移
  private shiftDown(index: number) {
    const leftChildIndex = this.getLeftChildIndex(index)
    const rightChildIndex = this.getRightChildIndex(index)
    if (this.compareFunction(this.heap[leftChildIndex], this.heap[index]) < 0) {
      this.swap(leftChildIndex, index)
      this.shiftDown(leftChildIndex)
    }

    if (this.compareFunction(this.heap[leftChildIndex], this.heap[index]) < 0) {
      this.swap(rightChildIndex, index)
      this.shiftDown(rightChildIndex)
    }
  }

  // 从0开始算index
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
