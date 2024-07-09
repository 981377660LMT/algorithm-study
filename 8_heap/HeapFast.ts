// HeapFast

export {}

if (require.main === module) {
  const data = [3, 1, 2, 4, 5]
  const innerHeap: IInnerHeap<number> = {
    len: () => data.length,
    less: (i, j) => data[i] < data[j],
    swap: (i, j) => {
      const tmp = data[i]
      data[i] = data[j]
      data[j] = tmp
    },
    push: x => data.push(x),
    pop: () => data.pop()!
  }

  {
    const N = 1e7
    console.time('Heap')
    heapify(innerHeap)
    for (let i = 0; i < N; i++) {
      heapPush(innerHeap, i)
    }
    for (let i = 0; i < N; i++) {
      heapPop(innerHeap)
    }
    console.timeEnd('Heap')
  }
}

/**
 * golang 风格的 heap 接口.
 * @link https://pkg.go.dev/container/heap#Interface
 */
interface IInnerHeap<T = any> {
  /** sort interface. */
  len: () => number
  less: (i: number, j: number) => boolean
  swap: (i: number, j: number) => void

  push: (x: T) => void
  pop: () => T
}

function heapify(h: IInnerHeap): void {
  // heapify
  const n = h.len()
  for (let i = (n >> 1) - 1; ~i; i--) {
    _down(h, i, n)
  }
}

function heapPush<T>(h: IInnerHeap<T>, x: T): void {
  h.push(x)
  _up(h, h.len() - 1)
}

function heapPop<T>(h: IInnerHeap<T>): T {
  const n = h.len() - 1
  h.swap(0, n)
  _down(h, 0, n)
  return h.pop()
}

function heapRemove<T>(h: IInnerHeap<T>, i: number): T {
  const n = h.len() - 1
  if (n !== i) {
    h.swap(i, n)
    if (!_down(h, i, n)) {
      _up(h, i)
    }
  }
  return h.pop()
}

function heapFix(h: IInnerHeap, i: number): void {
  if (!_down(h, i, h.len())) {
    _up(h, i)
  }
}

function _up(h: IInnerHeap, j: number): void {
  while (j) {
    const i = (j - 1) >>> 1 // parent
    if (i === j || !h.less(j, i)) break
    h.swap(i, j)
    j = i
  }
}

function _down(h: IInnerHeap, i0: number, n: number): boolean {
  let i = i0
  while (true) {
    const j1 = (i << 1) | 1
    if (j1 >= n || j1 < 0) break // j1 < 0 after int overflow
    let j = j1 // left child
    const j2 = j1 + 1
    if (j2 < n && h.less(j2, j1)) j = j2 // = 2*i + 2  // right child
    if (!h.less(j, i)) break
    h.swap(i, j)
    i = j
  }
  return i > i0
}
