import { Heap } from './Heap'

const INF = 2e15

/**
 * 用于高速化 Dijkstra. 边权不超过 int32.
 * @warning 实际上比普通的堆慢.
 */
class RadixHeap<T> {
  private readonly _v: [number, T][][] = Array(33)
    .fill(0)
    .map(() => [])
  private _size = 0
  private _last = 0

  top(): [number, T] {
    if (!this._v[0].length) {
      let i = 1
      while (!this._v[i].length) i++
      let min = INF
      this._v[i].forEach(p => {
        min = Math.min(min, p[0])
      })
      this._last = min
      this._v[i].forEach(p => {
        this._aux(p)
      })
      this._v[i] = []
    }
    return this._v[0][this._v[0].length - 1]
  }

  push(key: number, value: T): void {
    this._size++
    this._aux([key, value])
  }

  pop(): void {
    this._size--
    this.top()
    this._v[0].pop()
  }

  empty(): boolean {
    return this._size === 0
  }

  get size(): number {
    return this._size
  }

  private _aux(p: [number, T]): void {
    const bsr = 31 - Math.clz32(p[0] ^ this._last)
    this._v[bsr + 1].push(p)
  }
}

function dijkstraRadixHeap(
  n: number,
  adjList: [next: number, weight: number][][],
  start: number
): number[] {
  const dist = Array(n)
  for (let i = 0; i < n; ++i) dist[i] = INF
  dist[start] = 0
  const queue = new RadixHeap<number>()
  queue.push(0, start)

  while (queue.size) {
    const [curDist, cur] = queue.top()
    queue.pop()
    if (dist[cur] < curDist) continue
    adjList[cur].forEach(([next, w]) => {
      const cand = curDist + w
      if (dist[next] > cand) {
        dist[next] = cand
        queue.push(cand, next)
      }
    })
  }

  return dist
}

function dijkstra(n: number, adjList: [next: number, weight: number][][], start: number): number[] {
  const dist = Array(n)
  for (let i = 0; i < n; ++i) dist[i] = INF
  dist[start] = 0
  const queue = new Heap<[dist: number, vertex: number]>({ data: [], less: (a, b) => a[1] < b[1] })
  queue.push([0, start])

  while (queue.size) {
    const [curDist, cur] = queue.pop()!
    if (dist[cur] < curDist) continue
    adjList[cur].forEach(([next, w]) => {
      const cand = curDist + w
      if (dist[next] > cand) {
        dist[next] = cand
        queue.push([cand, next])
      }
    })
  }

  return dist
}

export {}

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  function networkDelayTime(times: number[][], n: number, k: number): number {
    const adjList: [v: number, w: number][][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    times.forEach(([u, v, w]) => {
      adjList[u - 1].push([v - 1, w])
    })
    const dist = dijkstraRadixHeap(n, adjList, k - 1)
    return Math.max(...dist) === INF ? -1 : Math.max(...dist)
  }

  const n = 2e5
  const edges = Array(n)
  for (let i = 0; i < n; ++i) edges[i] = [i + 1, i + 2, 1]
  const adjList: [number, number][][] = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []
  edges.forEach(([u, v, w]) => {
    adjList[u - 1].push([v - 1, w])
  })

  console.time('radix heap')
  dijkstraRadixHeap(n, adjList, 0)
  console.timeEnd('radix heap') // 51ms

  console.time('binary heap')
  dijkstra(n, adjList, 0)
  console.timeEnd('binary heap') // 20ms

  // RadixHeap更慢了
}
