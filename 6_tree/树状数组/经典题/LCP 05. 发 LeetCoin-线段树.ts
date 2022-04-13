// 线段树解法
class SegmentTreeNode {
  left = -1
  right = -1
  isLazy = false
  lazyValue = 0
  value = 0
}

/**
 * @description 线段树懒更新模板
 */
class SegmentTree {
  private tree: SegmentTreeNode[]

  constructor(size: number) {
    this.tree = Array.from({ length: size << 2 }, () => new SegmentTreeNode())
    this.build(1, 1, size)
  }

  update(root: number, left: number, right: number, delta: number): void {
    const node = this.tree[root]

    if (left <= node.left && node.right <= right) {
      node.isLazy = true
      node.lazyValue += delta
      node.value += delta * (node.right - node.left + 1)
      return
    }

    this.pushDown(root)
    const mid = (node.left + node.right) >> 1
    if (left <= mid) this.update(root << 1, left, right, delta)
    if (mid < right) this.update((root << 1) | 1, left, right, delta)
    this.pushUp(root)
  }

  query(root: number, left: number, right: number): number {
    const node = this.tree[root]
    if (left <= node.left && node.right <= right) {
      return node.value
    }

    this.pushDown(root)
    let res = 0
    const mid = (node.left + node.right) >> 1
    if (left <= mid) res += this.query(root << 1, left, right)
    if (mid < right) res += this.query((root << 1) | 1, left, right)
    return res
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
    this.pushUp(root)
  }

  /**
   * @param root 向下传递懒标记和懒更新的值 `isLazy`, `lazyValue`，并用 `lazyValue` 更新子区间的值
   */
  private pushDown(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]
    if (node.isLazy) {
      left.isLazy = true
      right.isLazy = true
      left.lazyValue += node.lazyValue
      right.lazyValue += node.lazyValue
      left.value += node.lazyValue * (left.right - left.left + 1)
      right.value += node.lazyValue * (right.right - right.left + 1)
      node.isLazy = false
      node.lazyValue = 0
    }
  }

  /**
   * @param root 用子节点更新父节点的值
   */
  private pushUp(root: number): void {
    const [node, left, right] = [this.tree[root], this.tree[root << 1], this.tree[(root << 1) | 1]]
    node.value = left.value + right.value
  }
}
// 力扣想进行的操作有以下三种：

// 给团队的一个成员（也可以是负责人）发一定数量的LeetCoin；
// 给团队的一个成员（也可以是负责人），以及他/她管理的所有人（即他/她的下属、他/她下属的下属，……），发一定数量的LeetCoin；
// 查询某一个成员（也可以是负责人），以及他/她管理的所有人被发到的LeetCoin之和。

// https://leetcode-cn.com/problems/coin-bonus/solution/xiao-ai-lao-shi-li-kou-bei-li-jie-zhen-t-rut3/
// https://mp.weixin.qq.com/s?__biz=MzkyMzI3ODgzNQ==&mid=2247483674&idx=1&sn=263595b26950ac60e5bf789839d70c9e&chksm=c1e6cd86f691449062d780b96d9af6d9590a71872ebfee980d5b045b9963714043261027c16b&token=1500097142&lang=zh_CN#rd
// 1. dfs将管理和他管理的人映射到一个区间(这部分很巧妙)[a,b] b表示自身的id
// 2. 树状数组区间update/query
const MOD = 1e9 + 7
function bonus(n: number, leadership: number[][], operations: number[][]): number[] {
  const adjList = Array.from<number, number[]>({ length: n + 1 }, () => [])
  const start = Array<number>(n + 1).fill(0) // 子树最开始的结点序号
  const end = Array<number>(n + 1).fill(0) // 本身最后映射到几
  // begin[1] = 1, end[1] = 6，表示编号为 1 的人所管理的团队映射到的区间是 [1, 6]，本身映射到 6
  let id = 1

  for (const [u, v] of leadership) adjList[u].push(v)

  dfs(1)

  const res: number[] = []
  const bit = new SegmentTree(n)
  for (const [optType, optId, optValue] of operations) {
    switch (optType) {
      case 1:
        bit.update(1, end[optId], end[optId], optValue)
        break
      case 2:
        bit.update(1, start[optId], end[optId], optValue)
        break
      case 3:
        const queryRes = bit.query(1, start[optId], end[optId])
        res.push(((queryRes % MOD) + MOD) % MOD)
        break
      default:
        throw new Error('invalid optType')
    }
  }

  return res

  // dfs序
  function dfs(cur: number): void {
    start[cur] = id
    for (const next of adjList[cur]) dfs(next)
    // id在dfs过程中被改变了
    end[cur] = id
    id++
  }
}

console.log(
  bonus(
    6,
    [
      [1, 2],
      [1, 6],
      [2, 3],
      [2, 5],
      [1, 4],
    ],
    [
      [1, 1, 500],
      [2, 2, 50],
      [3, 1],
      [2, 6, 15],
      [3, 1],
    ]
  )
)
// 第一次查询时，每个成员得到的LeetCoin的数量分别为（按编号顺序）：500, 50, 50, 0, 50, 0;
// 第二次查询时，每个成员得到的LeetCoin的数量分别为（按编号顺序）：500, 50, 50, 0, 50, 15.
// 输出：[650, 665]
export {}
