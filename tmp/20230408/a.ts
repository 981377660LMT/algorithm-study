const INF = 2e15

/**
 * 在线bfs.
 * 不预先给出图，而是通过两个函数 setUsed 和 findUnused 来在线寻找边.
 * @param setUsed 将 u 标记为已访问。
 * @param findUnused 找到和 u 邻接的一个未访问过的点。如果不存在, 返回 `null`。
 * https://leetcode.cn/problems/minimum-reverse-operations/solution/python-zai-xian-bfs-jie-jue-bian-shu-hen-y58m/
 */
function onlineBfs(
  n: number,
  start: number,
  setUsed: (cur: number) => void,
  findUnused: (cur: number) => number | null
): number[] {
  const dist = new Array(n).fill(INF)
  dist[start] = 0
  setUsed(start)
  let queue1 = new Uint32Array(n)
  let queue2 = new Uint32Array(n)

  queue1[0] = start
  let curQueue = queue1
  let nextQueue = queue2
  let curPtr = 1
  let nextPtr = 0

  while (curPtr) {
    for (let i = 0; i < curPtr; i++) {
      const cur = curQueue[i]
      while (true) {
        const next = findUnused(cur)
        if (next == null) {
          break
        }
        dist[next] = dist[cur] + 1 // weight
        nextQueue[nextPtr++] = next
        setUsed(next)
      }
    }
    ;[curQueue, nextQueue] = [nextQueue, curQueue]
    curPtr = nextPtr
    nextPtr = 0
  }

  return dist
}

function minimumVisitedCells(grid: number[][]): number {
  const ROW = grid.length
  const COL = grid[0].length
  const rowVisited: Finder[] = Array(ROW).fill(0)
  const colVisited: Finder[] = Array(COL).fill(0)
  for (let i = 0; i < ROW; i++) {
    rowVisited[i] = new Finder(COL)
    for (let j = 0; j < COL; j++) {
      rowVisited[i].insert(j)
    }
  }
  for (let j = 0; j < COL; j++) {
    colVisited[j] = new Finder(ROW)
    for (let i = 0; i < ROW; i++) {
      colVisited[j].insert(i)
    }
  }

  const dist = onlineBfs(
    ROW * COL,
    0,
    cur => {
      const r = ~~(cur / COL)
      const c = cur % COL
      rowVisited[r].erase(c)
      colVisited[c].erase(r)
    },
    cur => {
      const r = ~~(cur / COL)
      const c = cur % COL
      if (grid[r][c] === 0) return null
      const rightFirst = rowVisited[r].next(c)
      if (rightFirst && rightFirst <= c + grid[r][c]) {
        return r * COL + rightFirst
      }
      const downFirst = colVisited[c].next(r)
      if (downFirst && downFirst <= r + grid[r][c]) {
        return downFirst * COL + c
      }
      return null
    }
  )

  return dist[ROW * COL - 1] < INF ? 1 + dist[ROW * COL - 1] : -1
}

/**
 * 利用位运算寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
 * 初始时,所有位置都未被访问过.
 */
class Finder {
  private readonly _n: number
  private readonly _lg: number
  private readonly _seg: Uint32Array[]

  constructor(n: number) {
    this._n = n
    const seg: Uint32Array[] = []
    while (true) {
      seg.push(new Uint32Array((n + 31) >>> 5))
      n = (n + 31) >>> 5
      if (n <= 1) {
        break
      }
    }
    this._lg = seg.length
    this._seg = seg
  }

  insert(i: number): void {
    for (let h = 0; h < this._lg; h++) {
      this._seg[h][i >>> 5] |= 1 << (i & 31)
      i >>>= 5
    }
  }

  has(i: number): boolean {
    return (this._seg[0][i >>> 5] & (1 << (i & 31))) > 0
  }

  erase(i: number): void {
    for (let h = 0; h < this._lg; h++) {
      this._seg[h][i >>> 5] &= ~(1 << (i & 31))
      if (this._seg[h][i >>> 5] > 0) {
        break
      }
      i >>>= 5
    }
  }

  /**
   * 返回x右侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回null.
   */
  next(i: number): number | null {
    if (i < 0) {
      i = 0
    }
    if (i >= this._n) {
      return null
    }

    for (let h = 0; h < this._lg; h++) {
      if (i >>> 5 === this._seg[h].length) {
        break
      }
      let d = this._seg[h][i >>> 5] >>> (i & 31)
      if (d === 0) {
        i = (i >>> 5) + 1
        continue
      }
      // !trailingZeros32: 31 - Math.clz32(x & -x)
      i += 31 - Math.clz32(d & -d)
      for (let g = h - 1; ~g; g--) {
        i <<= 5
        const tmp = this._seg[g][i >>> 5]
        i += 31 - Math.clz32(tmp & -tmp)
      }
      return i
    }

    return null
  }

  /**
   * 返回x左侧第一个未被访问过的位置(包含x).
   * 如果不存在,返回null.
   */
  prev(i: number): number | null {
    if (i < 0) {
      return null
    }
    if (i >= this._n) {
      i = this._n - 1
    }

    for (let h = 0; h < this._lg; h++) {
      if (i === -1) {
        break
      }
      let d = this._seg[h][i >>> 5] << (31 - (i & 31))
      if (d === 0) {
        i = (i >>> 5) - 1
        continue
      }

      i -= Math.clz32(d)
      for (let g = h - 1; ~g; g--) {
        i <<= 5
        i += 31 - Math.clz32(this._seg[g][i >>> 5])
      }
      return i
    }

    return null
  }

  /**
   * 遍历[start,end)区间内的元素.
   */
  enumerateRange(start: number, end: number, f: (v: number) => void): void {
    let x: number | null = start - 1
    while (true) {
      x = this.next(x + 1)
      if (x == null || x >= end) {
        break
      }
      f(x)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.enumerateRange(0, this._n, v => sb.push(v.toString()))
    return `FastSet(${sb.join(', ')})`
  }
}
