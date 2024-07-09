type Token<T> = {
  value: T
  heapIndex: number
}

/**
 * 可修改堆.
 * @deprecated 比`ErasableHeap`慢.
 */
class ModifiableHeap<T> {
  private readonly _data: Token<T>[]
  private readonly _less: (a: T, b: T) => boolean

  constructor({ data, less }: { data: T[]; less: (a: T, b: T) => boolean }) {
    this._data = data.map((value, heapIndex) => ({ value, heapIndex }))
    this._less = less
    if (this._data.length > 1) this._heapify()
  }

  push(value: T): Token<T> {
    const res = { value, heapIndex: this._data.length }
    this._data.push(res)
    this._up(this._data.length - 1)
    return res
  }

  pop(): Token<T> {
    if (!this._data.length) throw new Error('pop from an empty heap')
    const n = this._data.length - 1
    this._swap(0, n)
    this._down(0, n)
    return this._data.pop()!
  }

  top(): Token<T> {
    if (!this._data.length) throw new Error('top from an empty heap')
    return this._data[0]
  }

  remove(token: Token<T>): void {
    this._remove(token.heapIndex)
  }

  modify(token: Token<T>, value: T): void {
    token.value = value
    this._fix(token.heapIndex)
  }

  clear(): void {
    this._data.length = 0
  }

  get size(): number {
    return this._data.length
  }

  private _heapify(): void {
    const n = this._data.length
    for (let i = (n >>> 1) - 1; ~i; i--) {
      this._down(i, n)
    }
  }

  private _up(j: number): void {
    const { _data, _less } = this
    while (j) {
      const i = (j - 1) >>> 1 // parent
      if (i === j || !_less(_data[j].value, _data[i].value)) break
      this._swap(i, j)
      j = i
    }
  }

  private _down(i0: number, n: number): boolean {
    const { _data, _less } = this
    let i = i0
    while (true) {
      const j1 = (i << 1) | 1
      if (j1 >= n || j1 < 0) break
      let j = j1
      const j2 = j1 + 1
      if (j2 < n && _less(_data[j2].value, _data[j1].value)) j = j2
      if (!_less(_data[j].value, _data[i].value)) break
      this._swap(i, j)
      i = j
    }
    return i > i0
  }

  private _fix(i: number): void {
    if (!this._down(i, this._data.length)) this._up(i)
  }

  private _remove(i: number): Token<T> {
    const n = this._data.length - 1
    if (n !== i) {
      this._swap(i, n)
      if (!this._down(i, n)) this._up(i)
    }
    return this._data.pop()!
  }

  private _swap(i: number, j: number): void {
    const tmp = this._data[i]
    this._data[i] = this._data[j]
    this._data[j] = tmp
    this._data[i].heapIndex = i
    this._data[j].heapIndex = j
  }
}

export { ModifiableHeap }

if (require.main === module) {
  const pq = new ModifiableHeap<number>({ data: [1, 8, 3, 5], less: (a, b) => a < b })
  console.log(pq.top())

  {
    const token = pq.push(4)
    console.log(pq.pop())
    console.log(pq.pop())
    pq.remove(token)
    console.log(pq.pop())
    console.log(pq.pop())
  }

  {
    const token = pq.push(-1)
    pq.push(111)
    pq.push(1)
    pq.modify(token, 100)
    console.log(pq.pop())
    console.log(pq.pop())
    console.log(pq.pop())
  }

  console.log(pq.size)

  const N = 1e7
  // test perf
  console.time('ErasableHeap')

  const pq2 = new ModifiableHeap<number>({ data: [], less: (a, b) => a < b })
  const tokens = Array<Token<number>>(N)
  for (let i = 0; i < N; ++i) {
    tokens[i] = pq2.push(i)
    pq2.top()
  }
  for (let i = 0; i < N; ++i) {
    pq2.remove(tokens[i])
  }
  console.timeEnd('ErasableHeap') // 1e7: 1.34s
}
