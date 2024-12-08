class Heap<T> {
  static mergeDestructively<E>(heap1: Heap<E>, heap2: Heap<E>): Heap<E> {
    if (heap1.size < heap2.size) {
      const tmp = heap1
      heap1 = heap2
      heap2 = tmp
    }
    for (let i = 0; i < heap2.size; i++) heap1.push(heap2._data[i])
    return heap1
  }

  private readonly _data: T[]
  private readonly _less: (a: T, b: T) => boolean

  constructor({ data, less }: { data: T[]; less: (a: T, b: T) => boolean }) {
    this._data = data.slice()
    this._less = less
    if (this._data.length > 1) this._heapify()
  }

  push(value: T): void {
    this._data.push(value)
    this._up(this._data.length - 1)
  }

  pop(): T {
    if (!this._data.length) throw new Error('pop from an empty heap')
    const n = this._data.length - 1
    this._swap(0, n)
    this._down(0, n)
    return this._data.pop()!
  }

  top(): T {
    if (!this._data.length) throw new Error('top from an empty heap')
    return this._data[0]
  }

  replace(value: T): T {
    if (!this._data.length) throw new Error('replace from an empty heap')
    const top = this._data[0]
    this._data[0] = value
    this._fix(0)
    return top
  }

  pushPop(value: T): T {
    if (this._data.length && this._less(this._data[0], value)) {
      const tmp = this._data[0]
      this._data[0] = value
      value = tmp
      this._fix(0)
    }
    return value
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
      const i = (j - 1) >>> 1
      if (i === j || !_less(_data[j], _data[i])) break
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
      if (j2 < n && _less(_data[j2], _data[j1])) j = j2
      if (!_less(_data[j], _data[i])) break
      this._swap(i, j)
      i = j
    }
    return i > i0
  }

  private _fix(i: number): void {
    if (!this._down(i, this._data.length)) this._up(i)
  }

  private _swap(i: number, j: number): void {
    const tmp = this._data[i]
    this._data[i] = this._data[j]
    this._data[j] = tmp
  }
}

function minOperations(n: number, m: number): number {
  const sN = n.toString()
  const sM = m.toString()
  if (sN.length !== sM.length) return -1
  if (isPrime(n) || isPrime(m)) return -1

  const length = sN.length
  const start = padNumber(n, length)
  const end = padNumber(m, length)
  const dist: { [k: string]: number } = {}
  dist[start] = n

  const heap = new Heap<[number, string]>({ data: [[n, start]], less: (a, b) => a[0] < b[0] })

  while (heap.size) {
    const [curCost, cur] = heap.pop()
    if (curCost > dist[cur]) continue
    if (cur === end) return curCost
    for (let i = 0; i < cur.length; i++) {
      const d = Number(cur[i])
      if (d < 9) {
        const up = cur.slice(0, i) + (d + 1) + cur.slice(i + 1)
        if (up[0] !== '0') {
          const x = Number(up)
          if (!isPrime(x)) {
            const val = curCost + x
            if (val < (dist[up] ?? Infinity)) {
              dist[up] = val
              heap.push([val, up])
            }
          }
        }
      }
      if (d > 0) {
        const down = cur.slice(0, i) + (d - 1) + cur.slice(i + 1)
        if (down[0] !== '0') {
          const y = Number(down)
          if (!isPrime(y)) {
            const val = curCost + y
            if (val < (dist[down] ?? Infinity)) {
              dist[down] = val
              heap.push([val, down])
            }
          }
        }
      }
    }
  }

  return -1
}

function isPrime(x: number): boolean {
  if (x <= 1) return false
  let i = 2
  while (i * i <= x) {
    if (x % i === 0) return false
    i++
  }
  return true
}

function padNumber(x: number, length: number): string {
  let s = x.toString()
  while (s.length < length) s = '0' + s
  return s
}
