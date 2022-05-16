// 查询区间[left,right]里的第k小数是几

class SegmentTreeNode {
  left = -1
  right = -1
  count = 0
}

interface StaticSegmentTree {
  /**
   * @description 查询区间[`left`,`right`]信息
   * @param left left从0开始
   * @param right right从0开始
   */
  query(left: number, right: number, ...args: any[]): number
  [other: string]: any
}

function usePersistentSegmentTree(nums: number[]): StaticSegmentTree {
  const n = nums.length
  const tree = Array<SegmentTreeNode>(4 * n + 16 * n) // 整个线段树N * 4 + NlogN，索引代表值域
  const roots = Array(n + 1).fill(0) // n+1个版本的根节点的treeId
  let treeId = 1

  // 离散化到0-allNums.length-1
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const mapping = new Map<number, number>()
  for (const [i, num] of allNums.entries()) mapping.set(num, i)

  // 建树骨架
  roots[0] = _build(0, allNums.length - 1)
  for (let i = 1; i <= n; i++) {
    // 后面的每插入一个点算一个版本, 每次插入都只是比上一个版本多1个数
    roots[i] = _update(roots[i - 1], 0, allNums.length - 1, mapping.get(nums[i - 1])!)
  }

  /**
   * @description 查询区间[`left`,`right`]里的第k小数是几
   * @param left left从0开始
   * @param right right从0开始
   * @param k  k从1开始
   */
  function query(left: number, right: number, k: number): number {
    if (0 <= left && left <= right && right + 1 <= n) {
      const rank = _query(roots[left], roots[right + 1], 0, allNums.length - 1, k)
      return allNums[rank]
    }

    throw new RangeError(`[left,right] out of range: [${left},${right}]`)
  }

  // 递归建树 返回树结点id build是建立好骨架, 每个版本insert改不同数据
  function _build(left: number, right: number): number {
    let curId = treeId++
    if (tree[curId] == undefined) tree[curId] = new SegmentTreeNode()
    if (left === right) return curId // 到底部了

    const mid = Math.floor((left + right) / 2)
    tree[curId].left = _build(left, mid)
    tree[curId].right = _build(mid + 1, right)
    return curId
  }

  // 插入新根节点 返回树结点id
  // 到left===right才插入了新点
  function _update(preRoot: number, left: number, right: number, value: number): number {
    const curId = treeId++
    if (tree[curId] == undefined) tree[curId] = new SegmentTreeNode()
    tree[curId].left = tree[preRoot].left
    tree[curId].right = tree[preRoot].right
    tree[curId].count = tree[preRoot].count
    if (left === right) {
      tree[curId].count++ // 插在这个叶节点上
      return curId
    }

    const mid = Math.floor((left + right) / 2)
    if (value <= mid) tree[curId].left = _update(tree[preRoot].left, left, mid, value)
    else tree[curId].right = _update(tree[preRoot].right, mid + 1, right, value)
    tree[curId].count = tree[tree[curId].left].count + tree[tree[curId].right].count
    return curId
  }

  // 二分值域查询
  function _query(
    preRoot: number,
    curRoot: number,
    left: number,
    right: number,
    k: number
  ): number {
    if (left === right) return left
    const leftHalfCount = tree[tree[curRoot].left].count - tree[tree[preRoot].left].count
    const mid = Math.floor((left + right) / 2)
    if (k <= leftHalfCount) return _query(tree[preRoot].left, tree[curRoot].left, left, mid, k)
    else return _query(tree[preRoot].right, tree[curRoot].right, mid + 1, right, k - leftHalfCount)
  }

  return {
    query,
  }
}

class Solution {
  private readonly tree: StaticSegmentTree

  constructor(nums: number[]) {
    this.tree = usePersistentSegmentTree(nums)
  }

  /**
   * @description 查询区间[`left`,`right`]里的第k小数是几
   */
  query(left: number, right: number, k: number): number {
    return this.tree.query(left - 1, right - 1, k)
  }
}

if (require.main === module) {
  const solution = new Solution([1, 5, 2, 6, 3, 7, 4])
  console.log(solution.query(2, 5, 3))
  console.log(solution.query(4, 4, 1))
  console.log(solution.query(1, 7, 3))
  // 输出样例：
  // 5
  // 6
  // 3
}

export {}
