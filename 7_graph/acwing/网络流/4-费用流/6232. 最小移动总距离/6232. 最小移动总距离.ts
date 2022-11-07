import { useMinCostMaxFlow } from '../useMinCostMaxFlow'

function minimumTotalDistance(robot: number[], factory: number[][]): number {
  const n = robot.length
  const m = factory.length
  const START = n + m
  const END = n + m + 1
  const mcmf = useMinCostMaxFlow(n + m + 2, START, END)

  // 源点到机器人
  for (let i = 0; i < n; i++) {
    mcmf.addEdge(START, i, 1, 0)
  }

  // 机器人到工厂
  for (let i = 0; i < n; i++) {
    for (let j = 0; j < m; j++) {
      mcmf.addEdge(i, n + j, 1, Math.abs(robot[i] - factory[j][0]))
    }
  }

  // 工厂到汇点
  for (let i = 0; i < m; i++) {
    mcmf.addEdge(n + i, END, factory[i][1], 0)
  }

  return mcmf.work()[1]
}

console.log(
  minimumTotalDistance(
    [0, 4, 6],
    [
      [2, 2],
      [6, 2]
    ]
  )
)
