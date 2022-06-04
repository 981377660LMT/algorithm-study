import assert from 'assert'
import { useDinic } from '../0-最大流模板/useDinic'
import { useDinicArray } from '../0-最大流模板/useDinicArray'

const DIR2 = [
  [0, 1],
  [1, 0],
]

// Map 1528 ms
function minimumOperations(grid: number[][]): number {
  const [ROW, COL] = [grid.length, grid[0].length]
  const [START, END] = [-1, -2]
  const maxFlow = useDinic(START, END)

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === 0) continue

      const cur = r * COL + c
      for (const [dr, dc] of DIR2) {
        const [nr, nc] = [r + dr, c + dc]
        if (0 <= nr && nr < ROW && 0 <= nc && nc < COL && grid[nr][nc] === 1) {
          const next = nr * COL + nc
          const [v1, v2] = (r + c) & 1 ? [cur, next] : [next, cur]
          maxFlow.addEdge(START, v1, 1)
          maxFlow.addEdge(v1, v2, 1)
          maxFlow.addEdge(v2, END, 1)
        }
      }
    }
  }

  return maxFlow.work()
}

// Array 1144 ms
function minimumOperations2(grid: number[][]): number {
  const [ROW, COL] = [grid.length, grid[0].length]
  const N = ROW * COL
  const [START, END] = [N + 2, N + 3]
  const maxFlow = useDinicArray(N, START, END)

  for (let r = 0; r < ROW; r++) {
    for (let c = 0; c < COL; c++) {
      if (grid[r][c] === 0) continue

      const cur = r * COL + c
      for (const [dr, dc] of DIR2) {
        const [nr, nc] = [r + dr, c + dc]
        if (0 <= nr && nr < ROW && 0 <= nc && nc < COL && grid[nr][nc] === 1) {
          const next = nr * COL + nc
          const [v1, v2] = (r + c) & 1 ? [cur, next] : [next, cur]
          maxFlow.addEdge(START, v1, 1)
          maxFlow.addEdge(v1, v2, 1)
          maxFlow.addEdge(v2, END, 1)
        }
      }
    }
  }

  return maxFlow.work()
}

if (require.main === module) {
  assert.strictEqual(
    minimumOperations([
      [1, 1, 0],
      [0, 1, 1],
      [1, 1, 1],
    ]),
    3
  )
  assert.strictEqual(
    minimumOperations2([
      [1, 1, 0],
      [0, 1, 1],
      [1, 1, 1],
    ]),
    3
  )
}
export {}
