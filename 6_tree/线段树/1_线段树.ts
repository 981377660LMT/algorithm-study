// 关心的两个区间的操作，这里是计算两个区间元素的和
class Merger<M extends number> {
  merge(a: M, b: M) {
    return a + b
  }
}

class SegmentTree<S extends number> {
  private arr: S[]
  private tree: S[]
  private merger: Merger<S>

  /**
   *
   * @param arr 输入的区间的浅拷贝
   * @param merger 定义两个区间如何融合
   */
  constructor(arr: S[], merger: Merger<S>) {
    this.arr = arr.slice()
    this.tree = Array<S>(this.arr.length * 4)
    this.merger = merger
    this.buildSegementTree(0, 0, this.arr.length - 1)
  }

  get size() {
    return this.arr.length
  }

  getValue(index: number) {
    if (index < 0 || index >= this.arr.length) throw new Error('illegal index')
    return this.arr[index]
  }

  /**
   *
   * @description 返回[l,r]存储的值
   */
  query(queryLeft: number, queryRight: number) {
    return this._query(0, queryLeft, queryRight, 0, this.arr.length - 1)
  }

  /**
   * @description 将index位置的值更新为val
   * @summary 递归更新左右孩子以及根节点的值
   */
  setValue(index: number, val: S) {
    this.tree[index] = val
    this._setValue(0, 0, this.arr.length - 1, index, val)
  }

  private _setValue(
    rootIndex: number,
    left: number,
    right: number,
    index: number,
    val: S
  ): unknown {
    if (left === right) {
      return (this.tree[index] = val)
    }

    const mid = Math.floor((left + right) / 2)
    const leftTreeIndex = this.leftChild(rootIndex)
    const rightTreeIndex = this.rightChild(rootIndex)
    if (index >= mid + 1) {
      return this._setValue(rightTreeIndex, mid + 1, right, index, val)
    } else if (index < mid) {
      return this._setValue(leftTreeIndex, left, mid, index, val)
    }

    this.tree[rootIndex] = this.merger.merge(
      this.tree[leftTreeIndex],
      this.tree[rightTreeIndex]
    ) as S
  }

  /**
   * @description 在rootIndex为根的线段树中的[left,right]的范围里搜寻[queryLeft,queryRight]的值
   * @summary 递归查找分情况三种讨论
   */
  private _query(
    rootIndex: number,
    queryLeft: number,
    queryRight: number,
    left: number,
    right: number
  ): S {
    if (left === queryLeft && right === queryRight) return this.tree[rootIndex]

    const mid = Math.floor((left + right) / 2)
    const leftTreeIndex = this.leftChild(rootIndex)
    const rightTreeIndex = this.rightChild(rootIndex)
    if (queryLeft >= mid + 1)
      return this._query(rightTreeIndex, queryLeft, queryRight, mid + 1, right)
    else if (queryRight <= mid) return this._query(leftTreeIndex, queryLeft, queryRight, left, mid)
    else {
      const leftResult = this._query(leftTreeIndex, queryLeft, mid, left, mid)
      const rightResult = this._query(rightTreeIndex, mid + 1, queryRight, mid + 1, right)
      return this.merger.merge(leftResult, rightResult) as S
    }
  }

  /**
   *
   * @param rootIndex
   * @param left
   * @param right
   * @description 在rootIndex的位置创建表示区间[l,r]的线段树
   */
  private buildSegementTree(rootIndex: number, left: number, right: number) {
    if (left === right) {
      this.tree[rootIndex] = this.arr[left]
      return
    }

    const mid = Math.floor((left + right) / 2)
    const leftTreeIndex = this.leftChild(rootIndex)
    const rightTreeIndex = this.rightChild(rootIndex)
    this.buildSegementTree(leftTreeIndex, left, mid)
    this.buildSegementTree(rightTreeIndex, mid + 1, right)
    this.tree[rootIndex] = this.merger.merge(
      this.tree[leftTreeIndex],
      this.tree[rightTreeIndex]
    ) as S
  }

  private leftChild(index: number) {
    return 2 * index + 1
  }

  private rightChild(index: number) {
    return 2 * index + 2
  }
}

const merger = new Merger<number>()
const st = new SegmentTree<number>([-2, 0, 3, -5, 2, -1], merger)
console.log(st)
// 计算1到3区间元素的和
console.log(st.query(1, 3))
st.setValue(3, 20)
console.log(st)

export { SegmentTree, Merger }
