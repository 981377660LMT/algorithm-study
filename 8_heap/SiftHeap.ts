// 二叉堆有个缺点：当你想维护固定若干个数据的最小值时，你可能需要`修改`某个数据的值。
// 然而，你并不能精准地找到这个值在堆中的位置，然后根据它的值的增大或者减小，将它在堆中动态“上浮”或者“下沉”。
// 经典案例就是在 `Dijkstra` 最短路算法借助 `std::priority_queue` 的常规实现中，
// 当你将到某个点的距离缩短时，你并未将该点在堆中的旧版本“上浮”，而是在堆中加入一个该点的新版本表示。
// 所以在同一时刻，堆中可能有同一元素的不同版本，除了最新版本外，其它旧版本都是没有作用的。
// 当有元素出栈时，还要判断这个元素是旧版本还是新版本，如果是旧版本就置之不理。这里实际上相当于“惰性删除”。
// 本数据结构解决了这一问题：在对某一点的值进行修改时，会在堆中该元素的旧版本基础上进行修改。
// （对于默认的大顶堆来说），如果将该点的值增大，可以调用“上浮”方法，在原来的位置基础上上浮；如果将该点的值减小，可以调用“下沉”方法，
// 在原来的位置基础上下沉。堆中，任何时刻，对于同一个点只存在一个版本的表示。
// 可以将 Dijkstra 算法的时间复杂度从 O(ElogE) 优化到 O(ElogV)。

/**
 * 原地升降二叉堆.维护0~n-1的n个元素.用于高速化Dijkstra算法.
 * 默认为小顶堆.
 *
 * @alias SiftHeap
 * @see {@link https://github.com/old-yan/CP-template/blob/a07b6fe0092e9ee890a0e35ada6ea1bb2c83ba05/DS/SiftHeap.md}
 */
class HeapUint32 {
  private readonly _heap: Uint32Array
  private readonly _pos: Int32Array
  private readonly _less: (i: number, j: number) => boolean
  private _heapPtr = 0

  /**
   * @param n [0,n).
   * @param less 下标i处的值是否小于下标j处的值.
   */
  constructor(n: number, less: (i: number, j: number) => boolean = (i, j) => i < j) {
    const pos = new Int32Array(n)
    for (let i = 0; i < n; i++) pos[i] = -1
    this._heap = new Uint32Array(n)
    this._pos = pos
    this._less = less
  }

  /**
   * 0<=i<n.
   */
  push(i: number): void {
    // === -1
    if (!~this._pos[i]) {
      this._heap[(this._pos[i] = this._heapPtr)] = i
      this._heapPtr++
    }
    this._siftUp(i)
  }

  pop(): number | undefined {
    if (!this._heapPtr) return undefined
    const res = this._heap[0]
    this._pos[res] = -1
    const ptr = --this._heapPtr
    if (ptr) {
      const tmp = this._heap[ptr]
      this._pos[tmp] = 0
      this._heap[0] = tmp
      this._siftDown(tmp)
    }
    return res
  }

  peek(): number | undefined {
    return this._heap[0]
  }

  get size(): number {
    return this._heapPtr
  }

  private _siftUp(i: number): void {
    let curPos = this._pos[i]
    for (
      let p = 0;
      // eslint-disable-next-line no-cond-assign
      curPos && this._less(i, (p = this._heap[(curPos - 1) >> 1]));
      curPos = (curPos - 1) >> 1
    ) {
      this._heap[(this._pos[p] = curPos)] = p
    }
    this._heap[(this._pos[i] = curPos)] = i
  }

  private _siftDown(i: number): void {
    let curPos = this._pos[i]
    // eslint-disable-next-line no-cond-assign
    for (let c = 0; (c = (curPos << 1) | 1) < this._heapPtr; curPos = c) {
      if (c + 1 < this._heapPtr && this._less(this._heap[c + 1], this._heap[c])) c++
      if (!this._less(this._heap[c], i)) break
      this._pos[(this._heap[curPos] = this._heap[c])] = curPos
    }
    this._heap[(this._pos[i] = curPos)] = i
  }
}

/**
 * 可删除元素的SiftHeap.
 */
class HeapUint32Erasable {
  private readonly _data: HeapUint32
  private readonly _erased: HeapUint32

  constructor(n: number, less: (i: number, j: number) => boolean = (i, j) => i < j) {
    this._data = new HeapUint32(n, less)
    this._erased = new HeapUint32(n, less)
  }

  /**
   * 0<=i<n.
   */
  push(i: number): void {
    this._data.push(i)
    this._normalize()
  }

  pop(): number | undefined {
    const value = this._data.pop()
    this._normalize()
    return value
  }

  peek(): number | undefined {
    return this._data.peek()
  }

  /**
   * 0<=value<n.
   */
  discard(i: number): void {
    this._erased.push(i)
    this._normalize()
  }

  private _normalize(): void {
    while (this._data.size && this._erased.size && this._data.peek() === this._erased.peek()) {
      this._data.pop()
      this._erased.pop()
    }
  }

  get size(): number {
    return this._data.size
  }
}

export { HeapUint32, HeapUint32Erasable }

if (require.main === module) {
  // time
  const n = 1e7
  const pq2 = new HeapUint32(n)
  console.time('HeapUint32')
  for (let i = 0; i < n; i++) {
    pq2.push(i)
    pq2.push(i)
  }
  for (let i = 0; i < n; i++) {
    pq2.pop()
    pq2.peek()
  }
  console.timeEnd('HeapUint32') // HeapUint32: 1.072s

  const values = ['banana', 'apple', 'orange', 'pear', 'peach']
  const pq = new HeapUint32(values.length, (i, j) => values[i] < values[j])
  pq.push(0)
  pq.push(1)
  pq.push(2)
  console.log(pq.peek()) // 1 => 字典序最小的是 apple
}
