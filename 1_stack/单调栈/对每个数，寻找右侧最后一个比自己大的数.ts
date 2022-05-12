import assert from 'assert'

class SegmentTreeNode {
  left = -1
  right = -1
  isLazy = false
  lazyValue = -1
  value = -1
}

/**
 * @description 线段树区间最大值更新
 * 0 <= A[i] <= 50000
 * 2 <= A.length <= 50000
 */

class SegmentTree {
  private readonly tree: SegmentTreeNode[]

  constructor(size: number) {
    this.tree = Array.from({ length: size << 2 }, () => new SegmentTreeNode())
    this.build(1, 1, size)
  }

  update(root: number, left: number, right: number, maxCand: number): void {
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) {
      node.isLazy = true
      node.lazyValue = Math.max(node.lazyValue, maxCand) // 注意这里
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
    if (left > right) return -Infinity

    const node = this.tree[root]
    if (left <= node.left && node.right <= right) {
      return node.value
    }

    this.pushDown(root)
    let res = -Infinity
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
      left.lazyValue = Math.max(left.lazyValue, node.lazyValue) // 注意这里
      right.lazyValue = Math.max(right.lazyValue, node.lazyValue)
      left.value = Math.max(left.value, node.lazyValue)
      right.value = Math.max(right.value, node.lazyValue)
      node.isLazy = false
      node.lazyValue = -Infinity
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
/**
 * @description 对每个数，寻找右侧最后一个比自己大的数；注意线段树要偏移
 * @param nums 0 <= nums[i] <= 50000  2 <= nums.length <= 50000
 * @returns 每个元素右侧最后一个比自己大的数的索引
 */
function findLastLarge(nums: number[]): number[] {
  const n = nums.length
  const res = Array<number>(n).fill(-1)
  const max = Math.max(...nums)

  const tree = new SegmentTree(max + 10)
  for (let i = n - 1; i >= 0; i--) {
    res[i] = tree.query(1, nums[i] + 1, nums[i] + 1)
    tree.update(1, 1, nums[i] + 1, i)
  }

  return res
}

if (require.main === module) {
  assert.deepStrictEqual(
    findLastLarge([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]),
    [9, 9, 9, 9, 9, 9, 9, 9, 9, -1]
  )

  assert.deepStrictEqual(
    findLastLarge([9, 8, 1, 0, 1, 9, 4, 0, 4, 1]),
    [5, 5, 9, 9, 9, -1, 8, 9, -1, -1]
  )
}

export { findLastLarge }
