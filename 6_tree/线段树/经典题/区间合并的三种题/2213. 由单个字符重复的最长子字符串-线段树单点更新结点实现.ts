// #region SegmentTree
class SegmentTreeNode {
  left = -1
  right = -1
  max = 1 // [left,right]区间内的最大连续数
  pre = 1 // 区间前缀的连续数
  suf = 1 // 区间后缀的连续数
}

class SegmentTree {
  private readonly tree: SegmentTreeNode[]
  readonly chars: string[]

  constructor(input: string) {
    const n = input.length
    this.tree = Array.from({ length: n << 2 }, () => new SegmentTreeNode())
    this.chars = input.split('')
    this.build(1, 1, n)
  }

  update(root: number, left: number, right: number, value: string): void {
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) {
      node.pre = 1
      node.suf = 1
      node.max = 1
      this.chars[left - 1] = value
      return
    }

    const mid = Math.floor((node.left + node.right) / 2)
    if (left <= mid) this.update(root << 1, left, right, value)
    if (mid < right) this.update((root << 1) | 1, left, right, value)
    this.pushUp(root)
  }

  query(root: number, left: number, right: number): number {
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) return node.max

    let res = 0
    const mid = Math.floor((left + right) / 2)
    if (left <= mid) res = Math.max(res, this.query(root << 1, left, right))
    if (mid < right) res = Math.max(res, this.query((root << 1) | 1, left, right))
    return res
  }

  queryAll(): number {
    return this.tree[1].max
  }

  private build(root: number, left: number, right: number): void {
    const node = this.tree[root]
    node.left = left
    node.right = right

    if (left === right) {
      node.pre = 1
      node.suf = 1
      node.max = 1
      return
    }

    const mid = Math.floor((left + right) / 2)
    this.build(root << 1, left, mid)
    this.build((root << 1) | 1, mid + 1, right)
    this.pushUp(root)
  }

  private pushUp(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]
    node.pre = left.pre
    node.suf = right.suf

    const mid = Math.floor((node.left + node.right) / 2)
    if (this.chars[mid - 1] === this.chars[mid]) {
      node.max = Math.max(left.max, right.max, left.suf + right.pre)
      if (left.pre === left.right - left.left + 1) node.pre += right.pre
      if (right.suf === right.right - right.left + 1) node.suf += left.suf
    } else {
      node.max = Math.max(left.max, right.max)
    }
  }
}

// #endregion
function longestRepeating(s: string, queryCharacters: string, queryIndices: number[]): number[] {
  const segmentTree = new SegmentTree(s)
  const res = Array<number>(queryIndices.length).fill(0)

  for (let i = 0; i < queryIndices.length; i++) {
    const [qc, qi] = [queryCharacters[i], queryIndices[i]]
    if (qc !== segmentTree.chars[qi]) segmentTree.update(1, qi + 1, qi + 1, qc)
    res[i] = segmentTree.queryAll()
  }

  return res
}

console.log(longestRepeating('babacc', 'bcb', [1, 3, 3]))
export {}
