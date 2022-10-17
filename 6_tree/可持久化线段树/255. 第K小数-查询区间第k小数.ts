// 查询区间[left,right]里的第k小数是几

function usePersistentSegmentTree(nums: number[]) {
  const n = nums.length
  const upper = 4 * n + Math.ceil(Math.log2(n)) * n // 离散化后整个线段树N * 4 + NlogN，索引代表值域
  const treeLeft = new Uint32Array(upper).fill(0)
  const treeRight = new Uint32Array(upper).fill(0)
  const treeCount = new Uint32Array(upper).fill(0)

  const roots = new Uint32Array(n + 1).fill(0) // !n+1个版本的根节点
  let treeId = 1 // !节点的版本号(唯一标识)

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
    if (left >= 0 && left <= right && right + 1 <= n) {
      const rank = _query(roots[left], roots[right + 1], 0, allNums.length - 1, k)
      return allNums[rank]
    }

    throw new RangeError(`[left,right] out of range: [${left},${right}]`)
  }

  // 递归建树 返回树结点id build是建立好骨架, 每个版本insert改不同数据
  function _build(left: number, right: number): number {
    const curId = treeId++
    if (left === right) return curId
    const mid = Math.floor((left + right) / 2)
    treeLeft[curId] = _build(left, mid)
    treeRight[curId] = _build(mid + 1, right)
    return curId
  }

  // 插入新根节点 返回树结点id
  // 到left===right才插入了新点
  function _update(preRoot: number, left: number, right: number, value: number): number {
    const curId = treeId++
    treeLeft[curId] = treeLeft[preRoot]
    treeRight[curId] = treeRight[preRoot]
    treeCount[curId] = treeCount[preRoot]
    if (left === right) {
      treeCount[curId]++ // value 插在这个叶节点上 频率加1
      return curId
    }

    const mid = Math.floor((left + right) / 2)
    if (value <= mid) treeLeft[curId] = _update(treeLeft[preRoot], left, mid, value)
    else treeRight[curId] = _update(treeRight[preRoot], mid + 1, right, value)
    treeCount[curId] = treeCount[treeLeft[curId]] + treeCount[treeRight[curId]]
    return curId
  }

  // 树上二分值域查询
  function _query(
    preRoot: number,
    curRoot: number,
    left: number,
    right: number,
    k: number
  ): number {
    if (left === right) return left
    const leftHalfCount = treeCount[treeLeft[curRoot]] - treeCount[treeLeft[preRoot]]
    const mid = Math.floor((left + right) / 2)
    if (k <= leftHalfCount) return _query(treeLeft[preRoot], treeLeft[curRoot], left, mid, k)
    return _query(treeRight[preRoot], treeRight[curRoot], mid + 1, right, k - leftHalfCount)
  }

  return {
    query
  }
}

/**
 * 查询区间第k小数是几
 * !把数都变成负数就变成区间第k大了
 */
class KthTree {
  private readonly _tree: ReturnType<typeof usePersistentSegmentTree>
  private readonly _isMin: boolean

  constructor(nums: number[], isMin: boolean) {
    nums = isMin ? nums.slice() : nums.map(num => -num)
    this._tree = usePersistentSegmentTree(nums)
    this._isMin = isMin
  }

  /**
   * 查询区间[`left`,`right`]里的第k小数是几
   *
   * @param left left >= 0
   * @param right right >= 0
   * @param k k >= 1
   */
  query(left: number, right: number, k: number): number {
    const res = this._tree.query(left, right, k)
    return this._isMin ? res : -res
  }
}

if (require.main === module) {
  const solution = new KthTree([1, 5, 2, 6, 3, 7, 4], false)
  console.log(solution.query(2 - 1, 5 - 1, 3))
  console.log(solution.query(4 - 1, 4 - 1, 1))
  console.log(solution.query(1 - 1, 7 - 1, 3))
  // 输出样例：
  // 5
  // 6
  // 3
}

export { KthTree }
