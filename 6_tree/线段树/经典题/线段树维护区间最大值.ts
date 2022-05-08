// 给出一个数组，最多删除一个连续子数组，求剩下数组的严格递增连续子数组的最大长度。
// n<=1e6

// 对位置为i的元素，最大长度为在它之前的所有比它小的元素的值的pre[j]的最大值 + suf[i] nlogn
// 用线段树维护小于该元素的最大的pre[j]

class SegmentTreeNode {
  left = -1
  right = -1
  isLazy = false
  lazyValue = 0
  value = 0
}

/**
 * @description 线段树懒更新模板
 */
class SegmentTree {
  private readonly tree: SegmentTreeNode[]

  constructor(size: number) {
    this.tree = Array.from({ length: size << 2 }, () => new SegmentTreeNode())
    this.build(1, 1, size)
  }

  update(root: number, left: number, right: number, maxCand: number): void {
    if (left > right) return
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) {
      node.isLazy = true
      node.lazyValue = maxCand
      node.value = Math.max(node.value, maxCand)
      return
    }

    this.pushDown(root)
    const mid = (node.left + node.right) >> 1
    if (left <= mid) this.update(root << 1, left, right, maxCand)
    if (mid < right) this.update((root << 1) | 1, left, right, maxCand)
    this.pushUp(root)
  }

  query(root: number, left: number, right: number): number {
    if (left > right) return 0
    const node = this.tree[root]
    if (left <= node.left && node.right <= right) {
      return node.value
    }

    this.pushDown(root)
    let res = 0
    const mid = (node.left + node.right) >> 1
    if (left <= mid) res = Math.max(res, this.query(root << 1, left, right))
    if (mid < right) res = Math.max(res, this.query((root << 1) | 1, left, right))
    return res
  }

  private build(root: number, left: number, right: number): void {
    const node = this.tree[root]
    node.left = left
    node.right = right

    if (left === right) {
      return
    }

    const mid = (node.left + node.right) >> 1
    this.build(root << 1, left, mid)
    this.build((root << 1) | 1, mid + 1, right)
    this.pushUp(root)
  }

  /**
   * @param root 向下传递懒标记和懒更新的值 `isLazy`, `lazyValue`，并用 `lazyValue` 更新子区间的值
   */
  private pushDown(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]
    if (node.isLazy) {
      left.isLazy = true
      right.isLazy = true
      left.lazyValue = node.lazyValue
      right.lazyValue = node.lazyValue
      left.value = Math.max(left.value, node.lazyValue)
      right.value = Math.max(right.value, node.lazyValue)
      node.isLazy = false
      node.lazyValue = 0
    }
  }

  /**
   * @param root 用子节点更新父节点的值
   */
  private pushUp(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]
    node.value = Math.max(left.value, right.value)
  }
}

// 假设数组都是正整数
function maxLenAfterRemove(nums: number[]): number {
  const n = nums.length
  const pre = Array<number>(n).fill(1)
  const suf = Array<number>(n).fill(1)

  for (let i = 0; i < n; i++) {
    if (nums[i] > nums[i - 1]) {
      pre[i] = pre[i - 1] + 1
    }
  }

  for (let i = n - 2; ~i; i--) {
    if (nums[i] < nums[i + 1]) {
      suf[i] = suf[i + 1] + 1
    }
  }

  let res = 1
  const max = Math.max(...nums)
  const tree = new SegmentTree(max + 10)
  for (let i = 0; i < n; i++) {
    const leftMax = tree.query(1, 1, nums[i] - 1)
    const right = suf[i]
    res = Math.max(res, leftMax + right)
    tree.update(1, nums[i], max, pre[i])
  }

  return res
}

console.log(maxLenAfterRemove([5, 3, 4, 9, 2, 8, 6, 7, 1])) // 4

export {}
