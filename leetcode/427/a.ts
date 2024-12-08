class Fenwick {
  private tree: number[]

  constructor(n: number) {
    this.tree = new Array(n + 1).fill(0)
  }

  add(i: number): void {
    while (i < this.tree.length) {
      this.tree[i] += 1
      i += i & -i
    }
  }

  pre(i: number): number {
    let res = 0
    while (i > 0) {
      res += this.tree[i]
      i &= i - 1
    }
    return res
  }

  query(l: number, r: number): number {
    return this.pre(r) - this.pre(l - 1)
  }
}

function pairwise(arr: number[]): [number, number][] {
  const pairs: [number, number][] = []
  for (let i = 0; i < arr.length - 1; i++) {
    pairs.push([arr[i], arr[i + 1]])
  }
  return pairs
}

function bisectLeft(arr: number[], x: number): number {
  let left = 0
  let right = arr.length
  while (left < right) {
    const mid = (left + right) >> 1
    if (arr[mid] < x) {
      left = mid + 1
    } else {
      right = mid
    }
  }
  return left
}

function maxRectangleArea(xCoord: number[], yCoord: number[]): number {
  const xMap = new Map<number, number[]>()
  const yMap = new Map<number, number[]>()

  for (let i = 0; i < xCoord.length; i++) {
    const x = xCoord[i]
    const y = yCoord[i]
    if (!xMap.has(x)) xMap.set(x, [])
    xMap.get(x)!.push(y)
    if (!yMap.has(y)) yMap.set(y, [])
    yMap.get(y)!.push(x)
  }

  const below = new Map<string, number>()
  for (const [x, ys] of xMap) {
    ys.sort((a, b) => a - b)
    for (const [y1, y2] of pairwise(ys)) {
      below.set(`${x},${y2}`, y1)
    }
  }

  const left = new Map<string, number>()
  for (const [y, xs] of yMap) {
    xs.sort((a, b) => a - b)
    for (const [x1, x2] of pairwise(xs)) {
      left.set(`${x2},${y}`, x1)
    }
  }

  const xs = Array.from(xMap.keys()).sort((a, b) => a - b)
  const ys = Array.from(yMap.keys()).sort((a, b) => a - b)

  const queries: Array<[number, number, number, number, number]> = []

  for (const [x2, listY] of xMap) {
    listY.sort((a, b) => a - b)
    for (const [y1, y2] of pairwise(listY)) {
      const x1 = left.get(`${x2},${y2}`)
      if (x1 !== undefined) {
        if (left.get(`${x2},${y1}`) === x1 && below.get(`${x1},${y2}`) === y1) {
          const dx1 = bisectLeft(xs, x1)
          const dx2 = bisectLeft(xs, x2)
          const dy1 = bisectLeft(ys, y1)
          const dy2 = bisectLeft(ys, y2)
          const area = (x2 - x1) * (y2 - y1)
          queries.push([dx1, dx2, dy1, dy2, area])
        }
      }
    }
  }

  const groupedQueries: Array<Array<[number, number, number, number]>> = []
  for (let i = 0; i < xs.length; i++) {
    groupedQueries.push([])
  }

  for (let i = 0; i < queries.length; i++) {
    const [x1, x2, y1, y2, _] = queries[i]
    if (x1 > 0) {
      groupedQueries[x1 - 1].push([i, -1, y1, y2])
    }
    groupedQueries[x2].push([i, 1, y1, y2])
  }

  const res = new Array<number>(queries.length).fill(0)
  const tree = new Fenwick(ys.length)

  for (let i = 0; i < xs.length; i++) {
    const x = xs[i]
    const qs = groupedQueries[i]
    const listY = xMap.get(x)!
    for (const y of listY) {
      tree.add(bisectLeft(ys, y) + 1)
    }
    for (const [qid, sign, Y1, Y2] of qs) {
      res[qid] += sign * tree.query(Y1 + 1, Y2 + 1)
    }
  }

  let ans = -1
  for (let i = 0; i < queries.length; i++) {
    const cnt = res[i]
    const area = queries[i][4]
    if (cnt === 4) {
      ans = Math.max(ans, area)
    }
  }

  return ans
}
