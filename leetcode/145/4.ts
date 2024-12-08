class UnionFindArray {
  private readonly _n: number
  private readonly _data: Int32Array
  private _part: number

  constructor(n: number) {
    this._n = n
    this._part = n
    this._data = new Int32Array(n).fill(-1)
  }

  union(x: number, y: number): boolean {
    let rootX = this.find(x)
    let rootY = this.find(y)
    if (rootX === rootY) return false
    if (this._data[rootX] > this._data[rootY]) {
      const tmp = rootX
      rootX = rootY
      rootY = tmp
    }
    this._data[rootX] += this._data[rootY]
    this._data[rootY] = rootX
    this._part -= 1
    return true
  }

  find(x: number): number {
    return this._data[x] < 0 ? x : (this._data[x] = this.find(this._data[x]))
  }

  getPart(): number {
    return this._part
  }
}

function countComponents(nums: number[], threshold: number): number {
  const n = nums.length
  const indexMap: { [key: number]: number } = {}
  for (let i = 0; i < n; i++) indexMap[nums[i]] = i
  const uf = new UnionFindArray(n)

  function gcd(a: number, b: number): number {
    while (b !== 0) {
      const t = a % b
      a = b
      b = t
    }
    return a
  }

  function ok(a: number, b: number): boolean {
    const g = gcd(a, b)
    const l = (a / g) * b
    return l <= threshold
  }

  const upper = threshold
  const multi: number[][] = []
  for (let i = 1; i <= upper; i++) multi.push([])

  for (let i = 0; i < n; i++) {
    const val = nums[i]
    if (val <= threshold) {
      for (let j = 1; j * j <= val; j++) {
        if (val % j === 0) {
          multi[j - 1].push(val)
          if (j * j !== val) multi[val / j - 1].push(val)
        }
      }
    }
  }

  for (let i = 0; i < multi.length; i++) {
    const arr = multi[i]
    if (arr.length < 2) continue
    arr.sort((a, b) => a - b)
    const base = arr[0]
    for (let k = 1; k < arr.length; k++) {
      if (ok(base, arr[k])) uf.union(indexMap[base], indexMap[arr[k]])
    }
  }

  return uf.getPart()
}
