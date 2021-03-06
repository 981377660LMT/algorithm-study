const { readFileSync } = require('fs')
const iter = readlines()
const input = (): string => iter.next().value
function* readlines(path = 0) {
  const lines = readFileSync(path)
    .toString()
    .trim()
    .split(/\r\n|\r|\n/)

  yield* lines
}
/////////////////////////////////////////

class SegmentTreeNode {
  left = -1
  right = -1
  isLazy = false
  lazyValue = 0 // 保存在A序列中的下标
  value = -1 // 初始全部为 -1 序列 B
}

class SegmentTree {
  private tree: SegmentTreeNode[]

  constructor(size: number) {
    this.tree = Array.from({ length: size << 2 }, () => new SegmentTreeNode())
    this.build(1, 1, size)
  }

  /**
   * @param aIndex 要更新的区间的A数组中的开始位置
   */
  update(root: number, left: number, right: number, aIndex: number): void {
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) {
      node.isLazy = true
      node.lazyValue = node.left - left + aIndex // 非常特别的更新，根据偏移量算一下当前结点对应A序列中那个位置
      node.value = nums[node.left - left + aIndex]
      return
    }

    this.pushDown(root)
    const mid = (node.left + node.right) >> 1
    if (left <= mid) this.update(root << 1, left, right, aIndex)
    if (mid < right) this.update((root << 1) | 1, left, right, aIndex)
    // this.pushUp(root)
  }

  query(root: number, left: number, right: number): number {
    const node = this.tree[root]
    if (left <= node.left && node.right <= right) {
      return node.value
    }

    this.pushDown(root)
    const mid = (node.left + node.right) >> 1

    // 单点查询
    if (left <= mid) return this.query(root << 1, left, right)
    else return this.query((root << 1) | 1, left, right)
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
    // this.pushUp(root)
  }

  /**
   * @param root 向下传递懒标记和懒更新的值 `isLazy`, `lazyValue`，并用 `lazyValue` 更新子区间的值
   */
  private pushDown(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]

    if (node.isLazy) {
      const mid = (node.left + node.right) >> 1

      left.isLazy = true
      left.value = nums[node.lazyValue]
      left.lazyValue = node.lazyValue

      right.isLazy = true
      // 左边是[node.left,mid] 右边是 [mid+1,node.right]
      right.value = nums[node.lazyValue + (mid + 1 - node.left)]
      right.lazyValue = node.lazyValue + (mid + 1 - node.left)

      node.isLazy = false
      node.lazyValue = 0
    }
  }

  /**
   * @param root 没有区间修改，所以并不需要pushup操作
   */
  private pushUp(root: number): void {}
}

export {}

const n = Number(input())
const nums = [0, ...input().split(' ').map(Number)] // 配合题目，线段树，index从1开始
const m = Number(input())

const tree = new SegmentTree(n)

for (let i = 0; i < m; i++) {
  const [opt, ...rest] = input().split(' ').map(Number)

  if (opt === 1) {
    // 把 A 序列中从下标 x 位置开始的连续 k 个元素粘贴到 B 序列中从下标 y 开始的连续 k 个位置上。
    // 输入数据可能会出现粘贴后 k 个元素超出 B 序列原有长度的情况，超出部分可忽略
    const [length, xStart, yStart] = rest
    tree.update(1, yStart, yStart + length - 1, xStart)
  } else {
    // 表示询问B序列下标 x 处的值是多少
    const [qi] = rest
    console.log(tree.query(1, qi, qi))
  }
}
