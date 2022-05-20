// # 区间内查询数字的频率
// # 求出给定子数组内一个给定值的 频率
class SegmentTreeNode {
  left = -1
  right = -1
  counter = new Map<number, number>()
}

class SegmentTree {
  private readonly tree: SegmentTreeNode[]

  constructor(nums: number[]) {
    const size = nums.length
    this.tree = Array.from({ length: size << 2 })
    this.build(1, 1, size, nums)
  }

  query(root: number, left: number, right: number, value: number): number {
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) {
      return node.counter.get(value) ?? 0
    }

    let res = 0
    const mid = Math.floor((node.left + node.right) / 2)
    if (left <= mid) res += this.query(root << 1, left, right, value)
    if (mid < right) res += this.query((root << 1) | 1, left, right, value)
    return res
  }

  private build(root: number, left: number, right: number, nums: number[]): void {
    if (this.tree[root] == undefined) this.tree[root] = new SegmentTreeNode()
    const node = this.tree[root]
    node.left = left
    node.right = right

    if (left === right) {
      
      node.counter.set(nums[left - 1], (node.counter.get(nums[left - 1]) ?? 0) + 1)
      return
    }

    const mid = Math.floor((node.left + node.right) / 2)
    this.build(root << 1, left, mid, nums)
    this.build((root << 1) | 1, mid + 1, right, nums)
    this.pushUp(root)
  }

  private pushUp(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]
    for (const [k, v] of left.counter) node.counter.set(k, (node.counter.get(k) ?? 0) + v)
    for (const [k, v] of right.counter) node.counter.set(k, (node.counter.get(k) ?? 0) + v)
  }
}

class RangeFreqQuery {
  private readonly tree: SegmentTree

  constructor(arr: number[]) {
    this.tree = new SegmentTree(arr)
  }

  // 1 <= arr[i], value <= 1e4
  // 注意线段树里的left从1开始
  query(left: number, right: number, value: number): number {
    return this.tree.query(1, left + 1, right + 1, value)
  }
}

if (require.main === module) {
  const rangeFreqQuery = new RangeFreqQuery([12, 33, 4, 56, 22, 2, 34, 33, 22, 12, 34, 56])
  console.log(rangeFreqQuery.query(1, 2, 4))
  console.log(rangeFreqQuery.query(0, 11, 33))
}

export {}
