import { useMinCostMaxFlow } from './useMinCostMaxFlow'

const INF = 2e15
function minimumTotalDistance(robot: number[], factory: number[][]): number {
  const n = robot.length
  const m = factory.length
  const START = n + m
  const END = n + m + 1
  const mcmf = useMinCostMaxFlow(n + m + 2, START, END)

  const positions: [pos: number, vertex: number][] = []

  // 源点到机器人
  for (let i = 0; i < n; i++) {
    mcmf.addEdge(START, i, 1, 0)
    positions.push([robot[i], i])
  }

  // 工厂到汇点
  for (let i = 0; i < m; i++) {
    const [pos, limit] = factory[i]
    mcmf.addEdge(n + i, END, limit, 0)
    positions.push([pos, n + i])
  }

  positions.sort((a, b) => a[0] - b[0])
  // 将工厂和机器人按照坐标排序,只连接X轴上相邻的两点,
  // 费用为两点间距离,容量为n(或者大于n的任意数值),总边数O(n+m)
  for (let i = 0; i < positions.length - 1; i++) {
    const [pos1, vertex1] = positions[i]
    const [pos2, vertex2] = positions[i + 1]
    mcmf.addEdge(vertex1, vertex2, INF, pos2 - pos1)
    mcmf.addEdge(vertex2, vertex1, INF, pos2 - pos1)
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

export {}
