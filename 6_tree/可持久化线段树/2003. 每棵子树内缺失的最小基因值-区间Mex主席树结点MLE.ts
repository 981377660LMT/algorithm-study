// 总共有 1e5 个基因值，基因值 互不相同,每个基因值都用 闭区间 [1, 1e5] 中的一个整数表示
function smallestMissingValueSubtree(parents: number[], nums: number[]): number[] {
  const n = nums.length
  const adjList: number[][] = Array(n)
    .fill(0)
    .map(() => [])
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

class SegmentTreeNode {
  left = 0
  right = 0
  lastPos = -1 // 这个值上一次被更新的时间
}

// 不做离散化 因为 题目里 1 <= nums[i] <= 1e5
function usePersistentSegmentTree(nums: number[]) {
  const n = nums.length
  const max = Math.max(...nums) + 1 // 注意把最大值加1 因为mex可能是max+1
  const tree = Array<SegmentTreeNode>(4 * max + Math.ceil(Math.log2(max)) * max)
  tree[0] = new SegmentTreeNode() // 注意这里创建0号 因为值域从1开始
  const roots = new Uint32Array(n + 1)
  let treeId = 1

  roots[0] = _build(1, max)
  for (let i = 1; i <= n; i++) {
    roots[i] = _update(roots[i - 1], 1, max, nums[i - 1], i - 1)
  }

  /**
   * @description 查询区间[`left`,`right`]里的mex
   * @param left left从0开始
   * @param right right从0开始
   */
  function query(left: number, right: number): number {
    if (0 <= left && left <= right && right + 1 <= n) {
      return _query(roots[left], roots[right + 1], 1, max, left)
    }

    throw new RangeError(`[left,right] out of range: [${left},${right}]`)
  }

  function _build(left: number, right: number): number {
    const curId = treeId++
    if (tree[curId] == undefined) tree[curId] = new SegmentTreeNode()
    if (left === right) return curId // 到底部了
    const mid = Math.floor((left + right) / 2)
    tree[curId].left = _build(left, mid)
    tree[curId].right = _build(mid + 1, right)
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
    if (tree[curId] == undefined) tree[curId] = new SegmentTreeNode()
    tree[curId].left = tree[preRoot].left
    tree[curId].right = tree[preRoot].right
    tree[curId].lastPos = tree[preRoot].lastPos
    if (left === right) {
      tree[curId].lastPos = time // 更新value值出现的位置
      return curId
    }

    const mid = Math.floor((left + right) / 2)
    if (value <= mid) tree[curId].left = _update(tree[preRoot].left, left, mid, value, time)
    else tree[curId].right = _update(tree[preRoot].right, mid + 1, right, value, time)
    tree[curId].lastPos = Math.min(tree[tree[curId].left].lastPos, tree[tree[curId].right].lastPos)
    return curId
  }

  // 二分值域查询 根据需要修改
  function _query(
    preRoot: number,
    curRoot: number,
    left: number,
    right: number,
    leftTime: number
  ): number {
    if (left === right) {
      return left
    }

    const lastPos = tree[tree[curRoot].left].lastPos
    const mid = Math.floor((left + right) / 2)
    // 若左子树的所有值里存在上次更新时间小于 leftTime 的值，那么就说明答案一定在左子树上
    if (lastPos < leftTime)
      return _query(tree[preRoot].left, tree[curRoot].left, left, mid, leftTime)
    return _query(tree[preRoot].right, tree[curRoot].right, mid + 1, right, leftTime)
  }

  return {
    query
  }
}

/**
 * @param n 树节点个数
 * @param adjList 无向图邻接表
 */
function useDfsOrder(n: number, adjList: number[][]) {
  const starts = new Uint32Array(n + 1) // 子树中最小的结点序号
  const ends = new Uint32Array(n + 1) // 子树中最大的结点序号，即自己的id
  let dfsId = 1
  dfs(0, -1)

  // 求dfs序
  function dfs(cur: number, pre: number): void {
    starts[cur] = dfsId
    for (const next of adjList[cur]) {
      if (next !== pre) dfs(next, cur)
    }
    ends[cur] = dfsId
    dfsId++
  }

  /**
   * @param root 求root所在子树映射到的区间
   * @returns [start, end] 1 <= start <= end <= n
   */
  function queryRange(root: number): [left: number, right: number] {
    return [starts[root], ends[root]]
  }

  /**
   * @param root 求root自身的dfsId
   * @returns dfsId 1 <= dfsId <= n
   */
  function queryId(root: number): number {
    return ends[root]
  }

  return { queryRange, queryId }
}

export {}
