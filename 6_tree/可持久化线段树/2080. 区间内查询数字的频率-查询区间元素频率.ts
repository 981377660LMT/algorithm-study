function usePersistentSegmentTree(nums: number[]) {
  const n = nums.length
  const max = Math.max(n, Math.max(...nums)) // 注意这里不能只用Math.max(...nums) (如果不离散化的话)
  const upper = 4 * max + Math.ceil(Math.log2(max)) * max
  const treeCount = new Uint32Array(upper).fill(0)
  const treeLeft = new Uint32Array(upper).fill(0)
  const treeRight = new Uint32Array(upper).fill(0)

  const roots = new Uint32Array(n + 1).fill(0) // n+1个版本的根节点的treeId
  let treeId = 1

  // 不离散化
  roots[0] = _build(1, max)
  for (let i = 1; i <= n; i++) {
    // 后面的每插入一个点算一个版本, 每次插入都只是比上一个版本多1个数
    roots[i] = _update(roots[i - 1], 1, max, nums[i - 1])
  }

  function query(left: number, right: number, value: number): number {
    if (0 <= left && left <= right && right + 1 <= n) {
      return _query(roots[left], roots[right + 1], 1, max, value)
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
    value: number
  ): number {
    if (left === right) {
      return treeCount[curRoot] - treeCount[preRoot] // 这个元素在这个区间里出现的次数
    }

    const mid = Math.floor((left + right) / 2)
    if (value <= mid) return _query(treeLeft[preRoot], treeLeft[curRoot], left, mid, value)
    return _query(treeRight[preRoot], treeRight[curRoot], mid + 1, right, value)
  }

  return {
    query,
  }
}

class RangeFreqQuery {
  private readonly tree: ReturnType<typeof usePersistentSegmentTree>

  constructor(arr: number[]) {
    this.tree = usePersistentSegmentTree(arr)
  }

  query(left: number, right: number, value: number): number {
    return this.tree.query(left, right, value)
  }
}

if (require.main === module) {
  const Q = new RangeFreqQuery([2, 2, 1, 2, 2])
  console.log(Q.query(2, 4, 1))
}

export {}
