// ![l,r]区间的最大连续长度就是左区间的最大连续长度,右区间最大连续长度,以及左右两区间结合在一起中间的最大连续长度.
// #region SegmentTree

class SegmentTree {
  private readonly tree: Uint32Array // 区间最大连续长度
  private readonly pre: Uint32Array // 前缀最大连续长度
  private readonly suf: Uint32Array // 后缀最大连续长度
  readonly chars: string[]

  constructor(input: string) {
    const n = input.length
    this.tree = new Uint32Array(n << 2)
    this.pre = new Uint32Array(n << 2)
    this.suf = new Uint32Array(n << 2)
    this.chars = input.split('')
    this.build(1, 1, n)
  }

  update(root: number, L: number, R: number, l: number, r: number, target: string): void {
    if (L <= l && r <= R) {
      this.chars[l - 1] = target
      return
    }

    const mid = Math.floor((l + r) / 2)
    if (L <= mid) this.update(root << 1, L, R, l, mid, target)
    if (mid < R) this.update((root << 1) | 1, L, R, mid + 1, r, target)
    this.pushUp(root, l, r, mid)
  }

  query(root: number, L: number, R: number, l: number, r: number): number {
    if (L <= l && r <= R) {
      return this.tree[root]
    }

    let res = 0
    const mid = Math.floor((l + r) / 2)
    if (L <= mid) res = Math.max(res, this.query(root << 1, L, R, l, mid))
    if (mid < R) res = Math.max(res, this.query((root << 1) | 1, L, R, mid + 1, r))
    return res
  }

  queryAll(): number {
    return this.tree[1]
  }

  private build(root: number, l: number, r: number): void {
    if (l === r) {
      this.tree[root] = 1
      this.pre[root] = 1
      this.suf[root] = 1
      return
    }

    const mid = Math.floor((l + r) / 2)
    this.build(root << 1, l, mid)
    this.build((root << 1) | 1, mid + 1, r)
    this.pushUp(root, l, r, mid)
  }

  private pushUp(root: number, l: number, r: number, mid: number): void {
    const [leftPre, rightPre] = [this.pre[root << 1], this.pre[(root << 1) | 1]]
    const [leftSuf, rightSuf] = [this.suf[root << 1], this.suf[(root << 1) | 1]]
    const [leftMax, rightMax] = [this.tree[root << 1], this.tree[(root << 1) | 1]]

    this.pre[root] = leftPre
    this.suf[root] = rightSuf

    // 是否需要跨越区间连接前后缀
    // 注意要先更新tree,再更新pre和suf
    if (this.chars[mid - 1] === this.chars[mid]) {
      this.tree[root] = Math.max(leftMax, rightMax, leftSuf + rightPre)
      if (leftPre === mid - l + 1) this.pre[root] += rightPre
      if (rightSuf === r - mid) this.suf[root] += leftSuf
    } else {
      this.tree[root] = Math.max(leftMax, rightMax)
    }
  }
}

// #endregion
function longestRepeating(s: string, queryCharacters: string, queryIndices: number[]): number[] {
  const n = s.length
  const segmentTree = new SegmentTree(s)
  const res = Array<number>(queryIndices.length).fill(0)

  for (let i = 0; i < queryIndices.length; i++) {
    const [qc, qi] = [queryCharacters[i], queryIndices[i]]
    if (qc !== segmentTree.chars[qi]) segmentTree.update(1, qi + 1, qi + 1, 1, n, qc)
    res[i] = segmentTree.queryAll()
  }

  return res
}

if (require.main === module) {
  console.log(longestRepeating('babacc', 'bcb', [1, 3, 3]))
}

export {}
