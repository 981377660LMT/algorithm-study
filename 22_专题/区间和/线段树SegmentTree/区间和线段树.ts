interface ISegmentTree<S> {
  query: (left: number, right: number) => S
  update: (index: number, val: S) => void
}

// 线段树的操作和堆非常像 可以结合堆理解
// 叶子节点存储着单一的值 叶子结点的父节点是merge的子节点的境地
// 左<=mid 右>=mid+1
// 比如[0,1,2,3,4] 左节点存储的[0,1,2] 右节点存储的[3,4]
// 比如[0,1,2,3,4,5] 左节点存储的[0,1,2] 右节点存储的[3,4,5]
// 动态查询 O(logn)
// 注意建树/查询/更新最后都要后序merge
class SegmentTree<S> implements ISegmentTree<S> {
  private data: S[]
  private tree: S[]
  private mergeFunc: (a: S, b: S) => S

  /**
   *
   * @param arr 输入的原数据区间的浅拷贝
   * @param mergeFunc 定义两个区间如何融合
   * @description 线段树对于长度为n的区间至多需要4n个节点(最好情况n=2^k 最坏情况n=2^k+1)
   */
  constructor(arr: S[], mergeFunc: (a: S, b: S) => S) {
    this.data = arr.slice()
    this.tree = Array<S>(this.data.length * 4)
    this.mergeFunc = mergeFunc
    this.buildSegementTree(0, 0, this.data.length - 1)
  }

  /**
   *
   * @param rootIndex
   * @param left
   * @param right
   * @description 以rootIndex的位置为根，创建表示区间[l,r]的线段树
   * 1. 找到mid划以分区间
   * 2. 左子树递归 右子树递归 递归终点为区间长度为1
   * 3. 融合左右子树
   * 不断二分 直到区间长度为1(叶子节点)
   */
  private buildSegementTree(rootIndex: number, left: number, right: number): void {
    if (left === right) {
      this.tree[rootIndex] = this.data[left]
      return
    }

    const mid = (left + right) >> 1
    const leftTreeIndex = this.leftChild(rootIndex)
    const rightTreeIndex = this.rightChild(rootIndex)
    this.buildSegementTree(leftTreeIndex, left, mid)
    this.buildSegementTree(rightTreeIndex, mid + 1, right)

    this.tree[rootIndex] = this.mergeFunc(this.tree[leftTreeIndex], this.tree[rightTreeIndex])
  }

  /**
   *
   * @description 查询[left,right]区间的值
   */
  query(left: number, right: number): S {
    return this._query(0, 0, this.data.length - 1, left, right)
  }

  /**
   * @description 在rootIndex为根的线段树的存储的区间的[left,right](大)的范围里搜寻[queryLeft,queryRight](小)的值
   * @summary 递归查找分情况三种讨论
   * 1. 要查找的区间在mid左边 查左边
   * 2. 要查找的区间在mid右边 查右边
   * 3. 要查找的区间在mid左边和右边 分别查询两边并merge
   * 4. 递归终点为当前节点存储的区间与查询区间匹配
   */
  private _query(
    rootIndex: number,
    left: number,
    right: number,
    queryLeft: number,
    queryRight: number
  ): S {
    if (left === queryLeft && right === queryRight) return this.tree[rootIndex]

    const mid = (left + right) >> 1
    const leftTreeIndex = this.leftChild(rootIndex)
    const rightTreeIndex = this.rightChild(rootIndex)
    if (queryLeft >= mid + 1) {
      // query的整个范围在右边
      return this._query(rightTreeIndex, mid + 1, right, queryLeft, queryRight)
    } else if (queryRight <= mid) {
      // query的整个范围在左边
      return this._query(leftTreeIndex, left, mid, queryLeft, queryRight)
    } else {
      // 左边找左边的query范围
      const leftResult = this._query(leftTreeIndex, left, mid, queryLeft, mid)
      // 右边找右边的query范围
      const rightResult = this._query(rightTreeIndex, mid + 1, right, mid + 1, queryRight)

      return this.mergeFunc(leftResult, rightResult)
    }
  }

  /**
   * @description 将index位置的值更新为val
   * @summary 递归找到index叶子节点，后序冒泡merge节点的值
   */
  update(index: number, val: S): void {
    this.data[index] = val
    this._update(0, 0, this.data.length - 1, index, val)
  }

  /**
   * @description 
   * 在[left,right] 这棵线段树中更新 index 叶子节点的值
     寻找index在哪里 很像二分搜索树查找节点
     二分搜索树是寻找节点值大于还是小于target 线段树寻找index对应哪个叶子节点是二分区间
     1.二分寻找index 在哪边
     2.递归终点为找到叶子节点(区间长度为1)
     3.后续merge 更新
   */
  private _update(rootIndex: number, left: number, right: number, index: number, val: S): void {
    if (left === right) {
      this.tree[rootIndex] = val
      return
    }

    const mid = (left + right) >> 1
    const leftTreeIndex = this.leftChild(rootIndex)
    const rightTreeIndex = this.rightChild(rootIndex)
    if (index >= mid + 1) {
      this._update(rightTreeIndex, mid + 1, right, index, val)
    } else if (index <= mid) {
      this._update(leftTreeIndex, left, mid, index, val)
    }

    // 后序，下面的节点已经处理完了
    this.tree[rootIndex] = this.mergeFunc(this.tree[leftTreeIndex], this.tree[rightTreeIndex])
  }

  private leftChild(index: number) {
    return 2 * index + 1
  }

  private rightChild(index: number) {
    return 2 * index + 2
  }
}

if (require.main === module) {
  const mergeFunc = (a: number, b: number) => a + b
  const st = new SegmentTree<number>([-2, 0, 3, -5, 2, -1], mergeFunc)
  console.log(st)
  // 计算1到3区间元素的和
  console.log(st.query(1, 3))
  st.update(3, 20)
  console.log(st)
}

export { SegmentTree }
