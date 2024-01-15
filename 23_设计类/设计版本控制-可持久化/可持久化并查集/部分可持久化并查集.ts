/**
 * 部分可持久化并查集(初始版本为0).
 */
class PartiallyPersistentUnionFind {
  private _version: number
  private readonly _data: Int32Array
  private readonly _last: Uint32Array
  private readonly _history: { time: number; size: number }[][]

  constructor(n: number) {
    const data = new Int32Array(n)
    const last = new Uint32Array(n)
    const history: { time: number; size: number }[][] = Array(n)
    for (let i = 0; i < n; i++) {
      data[i] = -1
      last[i] = 1e9
      history[i] = [{ time: 0, size: -1 }]
    }
    this._version = 0
    this._data = data
    this._last = last
    this._history = history
  }

  /**
   * 合并x和y所在的集合,返回当前版本号.
   */
  union(x: number, y: number, f?: (big: number, small: number) => void): number {
    this._version++
    x = this.find(this._version, x)
    y = this.find(this._version, y)
    if (x === y) {
      return this._version
    }
    if (this._data[x] > this._data[y]) {
      const tmp = x
      x = y
      y = tmp
    }
    this._data[x] += this._data[y]
    this._history[x].push({ time: this._version, size: this._data[x] })
    this._data[y] = x
    this._last[y] = this._version
    f && f(x, y)
    return this._version
  }

  find(time: number, x: number): number {
    if (time < this._last[x]) {
      return x
    }
    return this.find(time, this._data[x])
  }

  isConnected(time: number, x: number, y: number): boolean {
    return this.find(time, x) === this.find(time, y)
  }

  getSize(time: number, x: number): number {
    x = this.find(time, x)
    const tmp = this._history[x]
    let left = 0
    let right = tmp.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (tmp[mid].time > time) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return -tmp[left - 1].size
  }

  getGroups(time: number): Map<number, number[]> {
    const groups = new Map<number, number[]>()
    for (let i = 0; i < this._data.length; i++) {
      const root = this.find(time, i)
      if (!groups.has(root)) {
        groups.set(root, [])
      }
      groups.get(root)!.push(i)
    }
    return groups
  }

  get version(): number {
    return this._version
  }
}

export { PartiallyPersistentUnionFind }

if (require.main === module) {
  const uf = new PartiallyPersistentUnionFind(10)
  console.log(uf.getSize(0, 1))
  uf.union(1, 2)
  console.log(uf.getSize(1, 1))
  console.log(uf.getSize(0, 1))
  uf.union(2, 3)
  console.log(uf.getSize(1, 1))
  console.log(uf.getSize(2, 1))
  console.log(uf.getSize(0, 1))
}
