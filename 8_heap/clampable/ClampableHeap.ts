import { Heap } from '../Heap'

type Pair = { value: number; count: number }

class ClampableHeap {
  private readonly _clampMin: boolean
  private readonly _heap: Heap<Pair>
  private _total = 0
  private _count = 0

  /**
   * @param clampMin 为true时，支持容器内所有数与x取min；为false时，支持容器内所有数与x取max.
   */
  constructor(clampMin: boolean) {
    this._clampMin = clampMin
    this._heap = this._clampMin
      ? new Heap<Pair>((a, b) => b.value - a.value)
      : new Heap<Pair>((a, b) => a.value - b.value)
  }

  add(x: number): void {
    this._heap.push({ value: x, count: 1 })
    this._total += x
    this._count++
  }

  clamp(x: number): void {
    let newCount = 0
    if (this._clampMin) {
      while (this._heap.size > 0) {
        const top = this._heap.peek()!
        if (top.value < x) break
        this._heap.pop()
        this._total -= top.value * top.count
        newCount += top.count
      }
    } else {
      while (this._heap.size > 0) {
        const top = this._heap.peek()!
        if (top.value > x) break
        this._heap.pop()
        this._total -= top.value * top.count
        newCount += top.count
      }
    }
    this._total += x * newCount
    this._heap.push({ value: x, count: newCount })
  }

  sum(): number {
    return this._total
  }

  get size(): number {
    return this._count
  }
}

export { ClampableHeap }

if (require.main === module) {
  const pq1 = new ClampableHeap(true)
  pq1.add(1)
  pq1.add(2)
  pq1.add(3)
  console.assert(pq1.sum() === 6)
  pq1.clamp(2)
  console.assert(pq1.sum() === 5)
  pq1.clamp(1)
  pq1.add(2)
  console.assert(pq1.sum() === 5)

  const pq2 = new ClampableHeap(false)
  pq2.add(1)
  pq2.add(2)
  pq2.add(3)
  console.assert(pq2.sum() === 6)
  pq2.clamp(2)
  console.assert(pq2.sum() === 7)
  pq2.clamp(3)
  pq2.add(2)
  console.assert(pq2.sum() === 11)
}
