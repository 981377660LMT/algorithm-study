class BITArray {
  readonly n: number
  private readonly _data: number[]
  private _total = 0

  constructor(n: number, f?: (i: number) => number) {
    if (f === undefined) {
      this.n = n
      this._data = Array(n).fill(0)
    } else {
      this.n = n
      this._data = Array(n)
      for (let i = 0; i < n; i++) {
        this._data[i] = f(i)
        this._total += this._data[i]
      }
      for (let i = 1; i <= n; i++) {
        let j = i + (i & -i)
        if (j <= n) {
          this._data[j - 1] += this._data[i - 1]
        }
      }
    }
  }

  add(index: number, v: number): void {
    this._total += v
    index += 1
    while (index <= this.n) {
      this._data[index - 1] += v
      index += index & -index
    }
  }

  queryPrefix(end: number): number {
    if (end > this.n) {
      end = this.n
    }
    let res = 0
    while (end > 0) {
      res += this._data[end - 1]
      end -= end & -end
    }
    return res
  }

  queryRange(start: number, end: number): number {
    if (start < 0) start = 0
    if (end > this.n) end = this.n
    if (start >= end) return 0
    if (start === 0) return this.queryPrefix(end)
    let pos = 0
    let neg = 0
    while (end > start) {
      pos += this._data[end - 1]
      end &= end - 1
    }
    while (start > end) {
      neg += this._data[start - 1]
      start &= start - 1
    }
    return pos - neg
  }

  queryAll(): number {
    return this._total
  }
}

function bisectLeft(arrayLike: ArrayLike<number>, target: number): number {
  const n = arrayLike.length
  if (n === 0) return 0

  let left = 0
  let right = n - 1
  while (left <= right) {
    const mid = (left + right) >>> 1
    const midElement = arrayLike[mid]
    if (midElement < target) {
      left = mid + 1
    } else {
      right = mid - 1
    }
  }

  return left
}

function maxRectangleArea(xCoord: number[], yCoord: number[]): number {
  const h = (x: number, y: number) => `${x},${y}`
  const n = xCoord.length
  const points: number[][] = Array(n)
  for (let i = 0; i < n; i++) points[i] = [xCoord[i], yCoord[i]]
  points.sort((a, b) => a[0] - b[0] || a[1] - b[1])

  const xToYs: Map<number, number[]> = new Map()
  for (const [x, y] of points) {
    if (!xToYs.has(x)) xToYs.set(x, [])
    xToYs.get(x)!.push(y)
  }

  const yToXs: Map<number, number[]> = new Map()
  for (const [x, y] of points) {
    if (!yToXs.has(y)) yToXs.set(y, [])
    yToXs.get(y)!.push(x)
  }

  const pairToYs: Map<string, number[]> = new Map()
  for (const [y, xs] of yToXs) {
    xs.sort((a, b) => a - b)
    for (let i = 0; i + 1 < xs.length; i++) {
      const x1 = xs[i]
      const x2 = xs[i + 1]
      const key = h(x1, x2)
      if (!pairToYs.has(key)) pairToYs.set(key, [])
      pairToYs.get(key)!.push(y)
    }
  }

  const isAdjacent = (x: number, y1: number, y2: number): boolean => {
    const ys = xToYs.get(x)!
    let i1 = bisectLeft(ys, y1)
    let i2 = bisectLeft(ys, y2)
    return i1 >= 0 && i2 >= 0 && Math.abs(i1 - i2) === 1
  }

  const candidates: { x1: number; x2: number; y1: number; y2: number }[] = []
  for (const [key, ys] of pairToYs) {
    if (ys.length < 2) continue
    ys.sort((a, b) => a - b)
    const [x1, x2] = key.split(',').map(Number)
    for (let i = 0; i + 1 < ys.length; i++) {
      const y1 = ys[i]
      const y2 = ys[i + 1]
      if (x2 > x1 && y2 > y1) {
        if (isAdjacent(x1, y1, y2) && isAdjacent(x2, y1, y2)) {
          candidates.push({ x1, x2, y1, y2 })
        }
      }
    }
  }

  if (candidates.length === 0) return -1

  const allYs: number[] = []
  for (const [_, y] of points) {
    allYs.push(y)
  }
  for (const c of candidates) {
    allYs.push(c.y1, c.y2)
  }
  const sortedY = Array.from(new Set(allYs)).sort((a, b) => a - b)
  const compressY = new Map<number, number>()
  for (let i = 0; i < sortedY.length; i++) {
    compressY.set(sortedY[i], i + 1)
  }

  type Event1 = {
    x: number
    y1: number
    y2: number
    rectId: number
    type: 0 | 1
  }

  const queries: Event1[] = []
  const countX1: number[] = Array(candidates.length).fill(0)
  const countX2: number[] = Array(candidates.length).fill(0)
  for (let i = 0; i < candidates.length; i++) {
    const { x1, x2, y1, y2 } = candidates[i]
    const cy1 = compressY.get(y1)!
    const cy2 = compressY.get(y2)!
    queries.push({ x: x1, y1: cy1, y2: cy2, rectId: i, type: 0 })
    queries.push({ x: x2, y1: cy1, y2: cy2, rectId: i, type: 1 })
  }

  queries.sort((a, b) => a.x - b.x)
  type Event2 = { x: number; y: number }
  const pointEvents: Event2[] = points.map(p => ({ x: p[0], y: p[1] }))

  let qi = 0
  let pi = 0
  const events: {
    type: 0 | 1
    x: number
    y?: number
    y1?: number
    y2?: number
    rectId?: number
    queryType?: 0 | 1
  }[] = []

  while (qi < queries.length || pi < pointEvents.length) {
    if (qi < queries.length && (pi >= pointEvents.length || queries[qi].x < pointEvents[pi].x)) {
      const q = queries[qi++]
      events.push({
        type: 0,
        x: q.x,
        y1: q.y1,
        y2: q.y2,
        rectId: q.rectId,
        queryType: q.type
      })
    } else if (
      qi < queries.length &&
      pi < pointEvents.length &&
      queries[qi].x === pointEvents[pi].x
    ) {
      const q = queries[qi++]
      events.push({
        type: 0,
        x: q.x,
        y1: q.y1,
        y2: q.y2,
        rectId: q.rectId,
        queryType: q.type
      })
    } else {
      const p = pointEvents[pi++]
      events.push({ type: 1, x: p.x, y: p.y })
    }
  }

  const bit = new BITArray(sortedY.length)
  for (const e of events) {
    if (e.type === 0) {
      const cy1 = e.y1!
      const cy2 = e.y2!
      const upper = cy2 - 1
      const lower = cy1
      const count = bit.queryRange(lower, upper)
      if (e.queryType === 0) {
        countX1[e.rectId!] = count
      } else {
        countX2[e.rectId!] = count
      }
    } else {
      const cy = compressY.get(e.y!)!
      bit.add(cy - 1, 1)
    }
  }

  let res = -1
  for (let i = 0; i < candidates.length; i++) {
    const { x1, x2, y1, y2 } = candidates[i]
    const inner = countX2[i] - countX1[i]
    if (inner === 0) {
      const area = (x2 - x1) * (y2 - y1)
      if (area > res) res = area
    }
  }
  return res
}
