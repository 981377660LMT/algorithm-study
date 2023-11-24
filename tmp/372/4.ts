import { MonoStackDynamic } from '../../22_专题/离线查询/根号分治/RightMostLeftMostQuery'

export {}

const INF = 2e9 // !超过int32使用2e15

function leftmostBuildingQueries(heights: number[], queries: number[][]): number[] {
  const finder = new MonoStackDynamic(heights)
  const res: number[] = []
  for (let [alice, bob] of queries) {
    if (alice === bob) {
      res.push(alice)
      continue
    }

    if (alice > bob) {
      const tmp = alice
      alice = bob
      bob = tmp
    }

    if (heights[alice] < heights[bob]) {
      res.push(bob)
      continue
    }

    const cand = finder.rightNearestHigher(bob, heights[alice])
    res.push(cand)
  }

  return res
}
