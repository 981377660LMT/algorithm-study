// #region SegmentTree

class SegmentTree {
  private readonly _max: Uint32Array
  private readonly _preMax: Uint32Array
  private readonly _sufMax: Uint32Array
  private readonly _n: number
  readonly chars: string[]

  constructor(input: string) {
    this._n = input.length
    const cap = 1 << (32 - Math.clz32(this._n - 1) + 1)
    this._max = new Uint32Array(cap)
    this._preMax = new Uint32Array(cap)
    this._sufMax = new Uint32Array(cap)
    this.chars = input.split('')
    this.build(1, 1, this._n)
  }

  query(left: number, right: number): number {
    return this._query(1, left, right, 1, this._n)
  }

  update(left: number, right: number, lazy: string): void {
    this._update(1, left, right, 1, this._n, lazy)
  }

  queryAll(): number {
    return this._max[1]
  }

  private _update(root: number, L: number, R: number, l: number, r: number, target: string): void {
    if (L <= l && r <= R) {
      this.chars[l - 1] = target // !propagate
      return
    }

    const mid = Math.floor((l + r) / 2)
    if (L <= mid) this._update(root << 1, L, R, l, mid, target)
    if (mid < R) this._update((root << 1) | 1, L, R, mid + 1, r, target)
    this.pushUp(root, l, r, mid)
  }

  private _query(root: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) {
      return this._max[root]
    }

    let res = 0
    const mid = Math.floor((l + r) / 2)
    if (L <= mid) res = Math.max(res, this._query(root << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this._query((root << 1) | 1, L, R, mid + 1, r))
    return res
  }

  private build(root: number, l: number, r: number): void {
    if (l === r) {
      this._max[root] = 1
      this._preMax[root] = 1
      this._sufMax[root] = 1
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.build(root << 1, l, mid)
    this.build((root << 1) | 1, mid + 1, r)
    this.pushUp(root, l, r, mid)
  }

  // !op
  private pushUp(root: number, l: number, r: number, mid: number): void {
    const [leftPre, rightPre] = [this._preMax[root << 1], this._preMax[(root << 1) | 1]]
    const [leftSuf, rightSuf] = [this._sufMax[root << 1], this._sufMax[(root << 1) | 1]]
    const [leftMax, rightMax] = [this._max[root << 1], this._max[(root << 1) | 1]]

    this._preMax[root] = leftPre
    this._sufMax[root] = rightSuf

    if (this.chars[mid - 1] === this.chars[mid]) {
      this._max[root] = Math.max(leftMax, rightMax, leftSuf + rightPre)
      if (leftPre === mid - l + 1) this._preMax[root] += rightPre
      if (rightSuf === r - mid) this._sufMax[root] += leftSuf
    } else {
      this._max[root] = Math.max(leftMax, rightMax)
    }
  }
}

// #endregion
function longestRepeating(s: string, queryCharacters: string, queryIndices: number[]): number[] {
  const segmentTree = new SegmentTree(s)
  const res = Array<number>(queryIndices.length).fill(0)

  for (let i = 0; i < queryIndices.length; i++) {
    const [qc, qi] = [queryCharacters[i], queryIndices[i]]
    if (qc !== segmentTree.chars[qi]) {
      segmentTree.update(qi + 1, qi + 1, qc)
    }
    res[i] = segmentTree.queryAll()
  }

  return res
}

if (require.main === module) {
  console.log(longestRepeating('babacc', 'bcb', [1, 3, 3]))
}

export {}
