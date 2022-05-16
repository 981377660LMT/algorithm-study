// 使用主席树可以维护区间[L,R]中每个数值的次数，
// 然后，在这个root[R]-root[L-1]版本的线段树上，
// 查询节点的值大于等于threshold的节点的值。
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
   * @description 查询区间[`left`,`right`]里的绝对众数
   */
  function query(left: number, right: number, k: number): number {
    if (0 <= left && left <= right && right + 1 <= n) {
      const rank = _query(roots[left], roots[right + 1], 0, allNums.length - 1, k)
      return rank === -1 ? -1 : allNums[rank]
    }

    throw new RangeError(`[left,right] out of range: [${left},${right}]`)
  }

  // 递归建树 返回树结点id build是建立好骨架, 每个版本insert改不同数据
  function _build(left: number, right: number): number {
    const curId = treeId++
    if (tree[curId] == undefined) tree[curId] = new SegmentTreeNode()
    if (left === right) return curId // 到底部了

    const mid = Math.floor((left + right) / 2)
    tree[curId].left = _build(left, mid)
    tree[curId].right = _build(mid + 1, right)
    return curId
  }

  // 把这个新的值插在哪个叶子节点 返回树结点id
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
    if (left === right) {
      const valueCount = tree[curRoot].count - tree[preRoot].count // 这个元素在这个区间里出现的次数
      return valueCount >= k ? left : -1
    }

    // 值域在[left,mid]里的数的总个数
    const leftHalfCount = tree[tree[curRoot].left].count - tree[tree[preRoot].left].count
    const mid = Math.floor((left + right) / 2)
    if (k <= leftHalfCount) return _query(tree[preRoot].left, tree[curRoot].left, left, mid, k)
    return _query(tree[preRoot].right, tree[curRoot].right, mid + 1, right, k)
  }

  return {
    query,
  }
}

class MajorityChecker {
  private readonly tree: StaticSegmentTree

  constructor(arr: number[]) {
    this.tree = usePersistentSegmentTree(arr)
  }

  /**
   * @description 返回子数组中的元素  arr[left...right] 至少出现 threshold 次数，
   * 如果不存在这样的元素则返回 -1。
   */
  query(left: number, right: number, threshold: number): number {
    return this.tree.query(left, right, threshold)
  }
}

if (require.main === module) {
  const majorityChecker = new MajorityChecker([1, 1, 2, 2, 1, 1])
  console.log(majorityChecker.query(0, 5, 4))
  console.log(majorityChecker.query(0, 3, 3))
  console.log(majorityChecker.query(2, 3, 2))
  // 1 -1 2
}

export {}
