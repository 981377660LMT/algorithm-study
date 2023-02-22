import { useDinic } from '../0-最大流模板/0-useDinic'

const DIR2 = [
  [0, 1],
  [1, 0]
]
const INF = 2e15

function guardCastle(grid: string[]): number {
  const [ROW, COL] = [grid.length, grid[0].length]
  const OFFSET = ROW * COL
  const [START, END, TELEPORT] = [2 * OFFSET + 5, 2 * OFFSET + 6, 2 * OFFSET + 7]
  const maxFlow = useDinic(2 * OFFSET + 10, START, END)

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === '#') continue

      const cur = r * COL + c

      // 0. 所有点拆成 入点 和 出点 两个点
      maxFlow.addEdgeIfAbsent(cur, cur + OFFSET, grid[r][c] === '.' ? 1 : INF)

      // 1. 源点连接恶魔出生点
      if (grid[r][c] === 'S') {
        maxFlow.addEdgeIfAbsent(START, cur, INF)
      }

      // 2. 城堡连接汇点
      if (grid[r][c] === 'C') {
        maxFlow.addEdgeIfAbsent(cur, END, INF)
      }

      // 3. 虚拟点连通所有传送门
      if (grid[r][c] === 'P') {
        maxFlow.addEdgeIfAbsent(cur + OFFSET, TELEPORT, INF)
        maxFlow.addEdgeIfAbsent(TELEPORT, cur, INF)
      }

      // 4. 所有出点连通周围的入点
      for (const next of genNext(cur)) {
        maxFlow.addEdgeIfAbsent(cur + OFFSET, next, INF)
        maxFlow.addEdgeIfAbsent(next + OFFSET, cur, INF)
      }
    }
  }

  const minCut = maxFlow.calMaxFlow()
  return minCut < INF ? minCut : -1

  function* genNext(cur: number): Generator<number> {
    const [curRow, curCol] = [Math.floor(cur / COL), cur % COL]
    for (const [dr, dc] of DIR2) {
      const [nextRow, nextCol] = [curRow + dr, curCol + dc]
      if (
        nextRow >= 0 &&
        nextRow < ROW &&
        nextCol >= 0 &&
        nextCol < COL &&
        grid[nextRow][nextCol] !== '#'
      ) {
        yield nextRow * COL + nextCol
      }
    }
  }
}

if (require.main === module) {
  const grid = ['S.C.P#P.', '.....#.S']
  console.log(guardCastle(grid))
}

export {}
