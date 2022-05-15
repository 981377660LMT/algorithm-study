// 线段树解法
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
  private readonly size: number

  constructor(size: number) {
    this.size = size
    this.tree = Array.from({ length: size << 2 }, () => new SegmentTreeNode())
    this.build(1, 1, size)
  }

  update(root: number, left: number, right: number, delta: number): void {
    this.checkRange(left, right)
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) {
      node.isLazy = true
      node.lazyValue += delta
      node.value += delta * (node.right - node.left + 1)
      return
    }

    this.pushDown(root)
    const mid = (node.left + node.right) >> 1
    if (left <= mid) this.update(root << 1, left, right, delta)
    if (mid < right) this.update((root << 1) | 1, left, right, delta)
    this.pushUp(root)
  }

  // 对于不在线段树中的点，应按照题意来返回
  query(root: number, left: number, right: number): number {
    this.checkRange(left, right)
    const node = this.tree[root]
    if (left <= node.left && node.right <= right) {
      return node.value
    }

    this.pushDown(root)
    let res = 0
    const mid = (node.left + node.right) >> 1
    if (left <= mid) res += this.query(root << 1, left, right)
    if (mid < right) res += this.query((root << 1) | 1, left, right)
    return res
  }

  queryAll(): number {
    return this.tree[1].value
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
      left.lazyValue += node.lazyValue
      right.lazyValue += node.lazyValue
      left.value += node.lazyValue * (left.right - left.left + 1)
      right.value += node.lazyValue * (right.right - right.left + 1)
      node.isLazy = false
      node.lazyValue = 0
    }
  }

  /**
   * @param root 用子节点更新父节点的值
   */
  private pushUp(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]
    node.value = left.value + right.value
  }

  private checkRange(left: number, right: number): void {
    if (1 <= left && left <= right && right <= this.size) return
    throw new RangeError(`[SegmentTree] range error: [${left}, ${right}]`)
  }
}

export { SegmentTree, SegmentTreeNode }

if (require.main === module) {
  const sg = new SegmentTree(10)
  sg.update(1, 2, 3, 2)
  console.log(sg.query(1, 1, 8))
  console.log(sg.query(1, 1, 1))
  console.log(sg.queryAll())
}
