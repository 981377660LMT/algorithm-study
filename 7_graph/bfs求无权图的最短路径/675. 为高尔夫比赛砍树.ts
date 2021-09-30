import { ArrayDeque } from '../../2_queue/Deque'

// 把所有的树然后按照高度排序，利用BFS寻找两树之间的最短距离。
function cutOffTree(forest: number[][]): number {
  const trees: number[][] = []
  for (let i = 0; i < forest.length; i++) {
    for (let j = 0; j < forest[i].length; j++) {
      if (forest[i][j] > 1) trees.push([forest[i][j], i, j])
    }
  }
  trees.sort((a, b) => a[0] - b[0])

  let res = 0
  let [startRow, startCol] = [0, 0]

  for (const [_, treeRow, treeCol] of trees) {
    const dis = getDist(forest, startRow, startCol, treeRow, treeCol)
    if (dis < 0) return -1
    res += dis
    ;[startRow, startCol] = [treeRow, treeCol]
  }

  return res

  /**
   * @description 求无权图两点最短路径
   */
  function getDist(
    forest: number[][],
    startRow: number,
    startCol: number,
    treeRow: number,
    treeCol: number
  ): number {
    const [m, n] = [forest.length, forest[0].length]
    const dir = [
      [1, 0],
      [-1, 0],
      [0, 1],
      [0, -1],
    ]
    const queue = new ArrayDeque<[number, number, number]>(10 ** 4)
    queue.push([startRow, startCol, 0])
    const visited = new Set<number>([startRow * n + startCol])

    while (queue.length) {
      const [row, col, dis] = queue.shift()!
      if (row === treeRow && col === treeCol) return dis
      for (let [dx, dy] of dir) {
        const nextRow = row + dx
        const nextCol = col + dy
        if (
          nextRow >= 0 &&
          nextRow < m &&
          nextCol >= 0 &&
          nextCol < n &&
          forest[nextRow][nextCol] &&
          !visited.has(nextRow * n + nextCol)
        ) {
          queue.push([nextRow, nextCol, dis + 1])
          visited.add(nextRow * n + nextCol)
        }
      }
    }

    return -1
  }
}

console.log(
  cutOffTree([
    [1, 2, 3],
    [0, 0, 4],
    [7, 6, 5],
  ])
)

// 输出：6
// 解释：沿着上面的路径，你可以用 6 步，按从最矮到最高的顺序砍掉这些树。

// 0 表示障碍，无法触碰
// 1 表示地面，可以行走
// 比 1 大的数 表示有树的单元格，可以行走，数值表示树的高度
// 你需要按照树的高度从低向高砍掉所有的树，每砍过一颗树，该单元格的值变为 1（即变为地面）
// 你将从 (0, 0) 点开始工作，返回你砍完所有树需要走的最小步数。 如果你无法砍完所有的树，返回 -1 。

// (0,0) 位置的树，可以直接砍去，不用算步数。
// 只要不是0的地方都可以走，只是轮不到砍这个高度的树的时候不能砍
