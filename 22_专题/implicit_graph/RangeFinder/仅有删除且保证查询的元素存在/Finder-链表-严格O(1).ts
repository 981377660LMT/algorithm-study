class Finder {
  private readonly _n: number
  private readonly _exist: Int8Array
  private readonly _prev: Int32Array
  private readonly _next: Int32Array

  constructor(n: number) {
    this._n = n
    this._exist = new Int8Array(n)
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
    const l = this._prev[i]
    const r = this._next[i]
    if (l >= 0) this._next[l] = r
    if (r < this._n) this._prev[r] = l
    this._exist[i] = 0
    return true
  }

  /**
   * 返回`严格`小于i的最大元素,如果不存在,返回-1.
   * !调用时需保证i存在.
   */
  prev(i: number): number {
    return this._prev[i]
  }

  /**
   * 返回`严格`大于i的最小元素.如果不存在,返回n.
   * !调用时需保证i存在.
   */
  next(i: number): number {
    return this._next[i]
  }
}

export {}
