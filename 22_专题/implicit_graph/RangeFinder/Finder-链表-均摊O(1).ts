class Finder {
  private readonly _n: number
  private readonly _exist: Uint8Array
  private readonly _prev: Int32Array
  private readonly _next: Int32Array

  /**
   * 建立一个包含0~n-1的集合.
   */
  constructor(n: number) {
    this._n = n
    this._exist = new Uint8Array(n)
    this._prev = new Int32Array(n)
    this._next = new Int32Array(n)
    for (let i = 0; i < n; i++) {
      this._exist[i] = 1
      this._prev[i] = i - 1
      this._next[i] = i + 1
    }
  }

  has(i: number): boolean {
    return i >= 0 && i < this._n && !!this._exist[i]
  }

  erase(i: number): boolean {
    if (!this.has(i)) return false
    this._exist[i] = 0
    return true
  }

  /**
   * 返回小于等于i的最大元素.如果不存在,返回-1.
   */
  prev(i: number): number {
    if (i < 0) return -1
    if (i >= this._n) i = this._n - 1
    if (this._exist[i]) return i

    let realPrev = this._prev[i]
    while (realPrev >= 0 && !this._exist[realPrev]) {
      realPrev = this._prev[realPrev]
    }
    let cur = i
    while (cur >= 0 && cur !== realPrev) {
      const tmp = this._prev[cur]
      this._prev[cur] = realPrev
      cur = tmp
    }
    return realPrev
  }

  /**
   * 返回大于等于i的最小元素.如果不存在,返回n.
   */
  next(i: number): number {
    if (i < 0) i = 0
    if (i >= this._n) return this._n
    if (this._exist[i]) return i

    let realNext = this._next[i]
    while (realNext < this._n && !this._exist[realNext]) {
      realNext = this._next[realNext]
    }
    let cur = i
    while (cur < this._n && cur !== realNext) {
      const tmp = this._next[cur]
      this._next[cur] = realNext
      cur = tmp
    }
    return realNext
  }

  /**
   * 遍历[start,end)区间内的元素.
   */
  enumerate(start: number, end: number, f: (i: number) => void): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    for (let x = this.next(start); x < end; x = this.next(x + 1)) {
      f(x)
    }
  }

  toString(): string {
    const res: string[] = []
    this.enumerate(0, this._n, (i: number) => {
      res.push(i.toString())
    })
    return `Finder{${res}}`
  }
}

export {}

if (require.main === module) {
  const finder = new Finder(10)
  finder.erase(1)
  console.log(finder.toString())
  console.log(finder.prev(1))
  console.log(finder.next(1))
}
