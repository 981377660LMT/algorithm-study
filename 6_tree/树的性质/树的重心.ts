import { BinaryTree } from '../力扣加加/Tree'
// 计算以无根树每个点为根节点时的最大子树大小，这个值最小的点称为无根树的重心。
// 1.某个点是树的重心等价于它最大子树大小不大于整棵树大小的一半。
// 2.树至多有两个重心。如果树有两个重心，那么它们相邻。此时树一定有偶数个节点，且可以被划分为两个大小相等的分支，每个分支各自包含一个重心。
// 3.树中所有点到某个点的距离和中，到重心的距离和是最小的；如果有两个重心，那么到它们的距离和一样。反过来，距离和最小的点一定是重心。

// 找重心利用性质1，一趟dfs即可。
function findCentre(n: number, edges: [next: number, weight: number][][]): number[] {
  const res: number[] = []
  // 最大子树大小,即此节点为割点分割之后两半的最大大小
  const maxSizeOfSubtree = Array<number>(n).fill(Infinity)
  // 树的大小,即向`下面`走可以到多少个结点
  const treeSize = Array<number>(n).fill(Infinity)

  dfs(0, -1)
  return res

  function dfs(cur: number, parent: number): void {
    treeSize[cur] = 1
    maxSizeOfSubtree[cur] = 0

    for (const [next, _] of edges[cur]) {
      if (next === parent) continue
      dfs(next, cur)
      // 后序,更新cur:此时cur可以拿到各个next的信息
      treeSize[cur] += treeSize[next]
      maxSizeOfSubtree[cur] = Math.max(maxSizeOfSubtree[cur], treeSize[next])
    }

    // cur准备回退了，检查cur是否合法
    maxSizeOfSubtree[cur] = Math.max(maxSizeOfSubtree[cur], n - treeSize[cur])
    if (maxSizeOfSubtree[cur] <= n / 2) res.push(cur)
  }
}

export {}
console.log(
  findCentre(6, [
    [
      [1, 1],
      [2, 1],
      [5, 1],
    ],
    [
      [0, 1],
      [3, 1],
      [4, 1],
    ],
    [[0, 1]],
    [[1, 1]],
    [[1, 1]],
    [[0, 1]],
  ])
)
