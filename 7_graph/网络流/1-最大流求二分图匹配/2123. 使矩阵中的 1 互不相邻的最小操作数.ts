import assert from 'assert'
import { MaxFlowPushRelabel } from '../0-最大流模板/MaxFlowPushRelabel'

const DIR4 = [
  [0, 1],
  [0, -1],
  [1, 0],
  [-1, 0]
]

function minimumOperations(grid: number[][]): number {
  const [ROW, COL] = [grid.length, grid[0].length]
  const N = ROW * COL
  const START = N
  const END = START + 1
  const maxFlow = new MaxFlowPushRelabel(END + 1)

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === 0 || (r + c) & 1) continue
      const cur = r * COL + c
      for (const [dr, dc] of DIR4) {
        const [nr, nc] = [r + dr, c + dc]
        if (nr >= 0 && nr < ROW && nc >= 0 && nc < COL && grid[nr][nc] === 1) {
          const next = nr * COL + nc
          maxFlow.addEdge(START, cur, 1)
          maxFlow.addEdge(cur, next, 1)
          maxFlow.addEdge(next, END, 1)
        }
      }
    }
  }

  return maxFlow.maxFlow(START, END)
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
