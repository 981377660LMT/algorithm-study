import assert from 'assert'
import { useDinic } from '../0-最大流模板/0-useDinic'

const DIR2 = [
  [0, 1],
  [1, 0]
]

function minimumOperations(grid: number[][]): number {
  const [ROW, COL] = [grid.length, grid[0].length]
  const N = ROW * COL
  const [START, END] = [N + 2, N + 3]
  const maxFlow = useDinic(N + 5, START, END)

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === 0) continue

      const cur = r * COL + c
      for (const [dr, dc] of DIR2) {
        const [nr, nc] = [r + dr, c + dc]
        if (0 <= nr && nr < ROW && 0 <= nc && nc < COL && grid[nr][nc] === 1) {
          const next = nr * COL + nc
          const [v1, v2] = (r + c) & 1 ? [cur, next] : [next, cur]
          maxFlow.addEdgeIfAbsent(START, v1, 1)
          maxFlow.addEdgeIfAbsent(v1, v2, 1)
          maxFlow.addEdgeIfAbsent(v2, END, 1)
        }
      }
    }
  }

  return maxFlow.calMaxFlow()
}

if (require.main === module) {
  const G = [
    [1, 1, 0],
    [0, 1, 1],
    [1, 1, 1]
  ]
  assert.strictEqual(minimumOperations(G), 3)
}

export {}
