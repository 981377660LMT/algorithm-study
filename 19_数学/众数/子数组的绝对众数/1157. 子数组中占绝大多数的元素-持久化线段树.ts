// 使用主席树可以维护区间[L,R]中每个数值的次数，
// 然后，在这个root[R]-root[L-1]版本的线段树上，
// 查询节点的值大于等于threshold的节点的值。

function usePersistentSegmentTree(nums: number[]) {
  const n = nums.length
  // const tree = Array<SegmentTreeNode>(4 * n + 16 * n)
  const upper = 4 * n + Number(n).toString(2).length * n // 离散化后整个线段树N * 4 + NlogN，索引代表值域
  const treeCount = new Uint32Array(upper).fill(0)
  const treeLeft = new Uint32Array(upper).fill(0)
  const treeRight = new Uint32Array(upper).fill(0)

  const roots = new Uint32Array(n + 1).fill(0) // n+1个版本的根节点的treeId
  let treeId = 1

  // 离散化到0-allNums.length-1
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const mapping = new Uint32Array(allNums.length + 10).fill(0)
  for (const [i, num] of allNums.entries()) mapping[num] = i

  roots[0] = _build(0, allNums.length - 1)
  for (let i = 1; i <= n; i++) {
    roots[i] = _update(roots[i - 1], 0, allNums.length - 1, mapping[nums[i - 1]])
  }

  /**
   * @description 查询区间[`left`,`right`]里的绝对众数
   * @param left left从0开始
   * @param right right从0开始
   */
  function query(left: number, right: number, threshold: number): number {
    if (left >= 0 && left <= right && right + 1 <= n) {
      const rank = _query(roots[left], roots[right + 1], 0, allNums.length - 1, threshold)
      return rank === -1 ? -1 : allNums[rank]
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

  // 把这个新的值插在哪个叶子节点 返回树结点id
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

  // 二分值域查询 根据需要修改
  function _query(
    preRoot: number,
    curRoot: number,
    left: number,
    right: number,
    k: number
  ): number {
    if (left === right) {
      const valueCount = treeCount[curRoot] - treeCount[preRoot] // 这个元素在这个区间里出现的次数
      return valueCount >= k ? left : -1
    }

    // 值域在[left,mid]里的数的总个数
    const leftHalfCount = treeCount[treeLeft[curRoot]] - treeCount[treeLeft[preRoot]]
    const mid = Math.floor((left + right) / 2)
    if (k <= leftHalfCount) return _query(treeLeft[preRoot], treeLeft[curRoot], left, mid, k)
    return _query(treeRight[preRoot], treeRight[curRoot], mid + 1, right, k)
  }

  return {
    query
  }
}

class MajorityChecker {
  private readonly tree: ReturnType<typeof usePersistentSegmentTree>

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
