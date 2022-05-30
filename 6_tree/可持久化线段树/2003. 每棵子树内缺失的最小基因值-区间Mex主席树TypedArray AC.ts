// todo
// dfs遍历树，形成dfs数组。子树的dfs序是连续的，该题就变成：
// 在数组中查询若干区间的mex。
// 主席树找mex的思路
// 如果维护每个权值出现的次数。但是这个不能提供root[L]和root[R+1]作差得到结果。
// 考虑按照权值建线段树，对于每个线段树中的区间，
// 我们存储这个区间中最久没有被插入过的点上一次被修改的时间，
// 如果这个时间在所询问的左端点之前，那么就意味着在我们询问的区间里，
// 存在一个在该范围中的元素没有出现过，那么我们就可以这样查找下去找到这个元素了。
// 这个操作显然可以可持久化，于是建一棵主席树就好了。

// 执行用时：
// 1104 ms
// 内存消耗：
// 135.5 MB

import { useDfsOrder } from '../树的性质/dfs序/useDfsOrder'

// 总共有 1e5 个基因值，基因值 互不相同,每个基因值都用 闭区间 [1, 1e5] 中的一个整数表示
function smallestMissingValueSubtree(parents: number[], nums: number[]): number[] {
  const n = nums.length
  const adjList = Array.from<any, number[]>({ length: n }, () => [])
  for (const [cur, pre] of parents.entries()) {
    if (pre === -1) continue
    adjList[pre].push(cur)
    adjList[cur].push(pre)
  }

  const { queryRange, queryId } = useDfsOrder(n, adjList)
  const newNums = Array(n).fill(0)
  for (let root = 0; root < n; root++) {
    const dfsId = queryId(root)
    newNums[dfsId - 1] = nums[root]
  }

  const segmentTree = usePersistentSegmentTree(newNums)
  const mex = Array(n).fill(1)
  for (let root = 0; root < n; root++) {
    const [left, right] = queryRange(root)
    mex[root] = segmentTree.query(left - 1, right - 1)
  }

  return mex
}

// 不做离散化 因为 题目里 1 <= nums[i] <= 1e5
function usePersistentSegmentTree(nums: number[]) {
  const n = nums.length
  const max = Math.max(n, Math.max(...nums) + 1) // 注意把最大值加1 因为mex可能是max+1
  const upper = 4 * max + Math.ceil(Math.log2(max)) * max
  const treeLeft = new Uint32Array(upper).fill(0)
  const treeRight = new Uint32Array(upper).fill(0)
  const lastPos = new Int32Array(upper).fill(-1) // 这个值上一次被更新的时间

  const roots = new Uint32Array(n + 1).fill(0)
  let treeId = 1

  roots[0] = _build(1, max)
  for (let i = 1; i <= n; i++) {
    roots[i] = _update(roots[i - 1], 1, max, nums[i - 1], i - 1)
  }

  /**
   * @description 查询区间[`left`,`right`]里的mex
   */
  function query(left: number, right: number): number {
    if (0 <= left && left <= right && right + 1 <= n) {
      return _query(roots[right + 1], 1, max, left)
    }

    throw new RangeError(`[left,right] out of range: [${left},${right}]`)
  }

  function _build(left: number, right: number): number {
    const curId = treeId++
    if (left === right) return curId // 到底部了
    const mid = Math.floor((left + right) / 2)
    treeLeft[curId] = _build(left, mid)
    treeRight[curId] = _build(mid + 1, right)
    return curId
  }

  function _update(
    preRoot: number,
    left: number,
    right: number,
    value: number,
    time: number
  ): number {
    const curId = treeId++
    treeLeft[curId] = treeLeft[preRoot]
    treeRight[curId] = treeRight[preRoot]
    lastPos[curId] = lastPos[preRoot]
    if (left === right) {
      lastPos[curId] = time // 更新value值出现的位置
      return curId
    }

    const mid = Math.floor((left + right) / 2)
    if (value <= mid) treeLeft[curId] = _update(treeLeft[preRoot], left, mid, value, time)
    else treeRight[curId] = _update(treeRight[preRoot], mid + 1, right, value, time)
    lastPos[curId] = Math.min(lastPos[treeLeft[curId]], lastPos[treeRight[curId]])
    return curId
  }

  // 二分值域查询 根据需要修改
  function _query(curRoot: number, left: number, right: number, leftTime: number): number {
    if (left === right) {
      return left
    }

    const last = lastPos[treeLeft[curRoot]]
    const mid = Math.floor((left + right) / 2)
    // 若左子树的所有值里存在上次更新时间小于 leftTime 的值，那么就说明答案一定在左子树上
    if (last < leftTime) return _query(treeLeft[curRoot], left, mid, leftTime)
    return _query(treeRight[curRoot], mid + 1, right, leftTime)
  }

  return {
    query,
  }
}

if (require.main === module) {
  // console.log(smallestMissingValueSubtree([-1, 0, 0, 2], [1, 2, 3, 4]))
  // console.log(smallestMissingValueSubtree([-1, 2, 3, 0, 2, 4, 1], [2, 3, 4, 5, 6, 7, 8]))
  console.log(smallestMissingValueSubtree([-1, 0, 1, 0, 3, 3], [5, 4, 6, 2, 1, 3])) // [7,1,1,4,2,1]
}
