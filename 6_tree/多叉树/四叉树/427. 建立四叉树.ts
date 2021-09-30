class Node {
  val: boolean // val：储存叶子结点所代表的区域的值。1 对应 True，0 对应 False
  isLeaf: boolean // 当这个节点是一个叶子结点时为 True，如果它有 4 个子节点则为 False
  topLeft: Node | null
  topRight: Node | null
  bottomLeft: Node | null
  bottomRight: Node | null
  constructor(
    val: boolean,
    isLeaf: boolean,
    topLeft: Node | null,
    topRight: Node | null,
    bottomLeft: Node | null,
    bottomRight: Node | null
  ) {
    this.val = val
    this.isLeaf = isLeaf
    this.topLeft = topLeft
    this.topRight = topRight
    this.bottomLeft = bottomLeft
    this.bottomRight = bottomRight
  }
}

function construct(grid: (0 | 1)[][]): Node | null {
  return makeTree(0, 0, grid.length)

  function makeTree(x: number, y: number, edgeLength: number): Node | null {
    const val = getMatrixValue(x, y, edgeLength)
    let node: Node

    if (val === 2) {
      const half = ~~(edgeLength / 2)
      node = new Node(
        true,
        false,
        makeTree(x, y, half),
        makeTree(x, y + half, half),
        makeTree(x + half, y, half),
        makeTree(x + half, y + half, half)
      )
    } else {
      node = new Node(val === 1 ? true : false, true, null, null, null, null)
    }

    return node
  }

  function getMatrixValue(x: number, y: number, edgeLength: number): 0 | 1 | 2 {
    const val = grid[x][y]

    for (let i = x; i < x + edgeLength; i++) {
      for (let j = y; j < y + edgeLength; j++) {
        if (grid[i][j] !== val) return 2
      }
    }

    return val
  }
}

console.log(
  construct([
    [0, 1],
    [1, 0],
  ])
)
// 输出：[[0,1],[1,0],[1,1],[1,1],[1,0]]   [0,1] 表示虚拟根节点不是叶子 值为1

// 判断当前矩阵是否全为1或全为0，如果是则创建四叉树叶子节点，
// 否则创建四叉树非叶子节点，4个子分别是将矩阵分为4份构成的子树
// 时间复杂度：O(n^4)
// 空间复杂度：O(logn)

export { construct }
