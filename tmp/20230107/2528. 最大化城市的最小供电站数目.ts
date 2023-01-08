// 2528. 最大化城市的最小供电站数目

class BITArray {
  private static _build(arrayLike: ArrayLike<number>): number[] {
    const tree = Array(arrayLike.length + 1).fill(0)
    for (let i = 1; i < tree.length; i++) {
      tree[i] += arrayLike[i - 1]
      const parent = i + (i & -i)
      if (parent < tree.length) tree[parent] += tree[i]
    }
    return tree
  }

  readonly length: number
  private readonly _tree: number[]

  constructor(lengthOrArrayLike: number | ArrayLike<number>) {
    if (typeof lengthOrArrayLike === 'number') {
      this.length = lengthOrArrayLike
      this._tree = Array(lengthOrArrayLike + 1).fill(0)
    } else {
      this.length = lengthOrArrayLike.length
      this._tree = BITArray._build(lengthOrArrayLike)
    }
  }

  add(index: number, delta: number): void {
    if (index <= 0) throw new RangeError(`add index must be greater than 0, but got ${index}`)
    for (let i = index; i <= this.length; i += i & -i) {
      this._tree[i] += delta
    }
  }

  query(right: number): number {
    if (right > this.length) right = this.length
    let res = 0
    for (let i = right; i > 0; i -= i & -i) {
      res += this._tree[i]
    }
    return res
  }

  queryRange(left: number, right: number): number {
    return this.query(right) - this.query(left - 1)
  }

  toString(): string {
    const sb: string[] = []
    sb.push('BITArray: [')
    for (let i = 1; i < this._tree.length; i++) {
      sb.push(String(this.queryRange(i, i)))
      if (i < this._tree.length - 1) sb.push(', ')
    }
    sb.push(']')
    return sb.join('')
  }
}

function maxPower(stations: number[], r: number, k: number): number {
  let left = 1
  let right = 2e15
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (check(mid)) left = mid + 1
    else right = mid - 1
  }

  return right

  // !可以用滑窗优化
  function check(mid: number): boolean {
    const bit = new BITArray(stations)
    let curK = k
    for (let i = 1; i <= stations.length; i++) {
      const cur = bit.queryRange(Math.max(1, i - r), Math.min(stations.length, i + r))
      if (cur < mid) {
        const diff = mid - cur
        bit.add(Math.min(stations.length, i + r), diff)
        curK -= diff
        if (curK < 0) return false
      }
    }
    return true
  }
}

// stations = [1,2,4,5,0], r = 1, k = 2
console.log(maxPower([1, 2, 4, 5, 0], 1, 2))
const bit = new BITArray([1, 2, 3, 4, 5])
console.log(bit.queryRange(1, 2))
bit.add(1, 1)
console.log(bit.queryRange(1, 1))
export {}
