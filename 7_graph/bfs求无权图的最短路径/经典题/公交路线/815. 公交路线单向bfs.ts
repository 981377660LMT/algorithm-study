// 815. 公交路线
// https://leetcode.cn/problems/bus-routes/description/
// 给你一个数组 routes ，表示一系列公交线路，其中每个 routes[i] 表示一条公交线路，第 i 辆公交车将会在上面循环行驶。
// 例如，路线 routes[0] = [1, 5, 7] 表示第 0 辆公交车会一直按序列 1 -> 5 -> 7 -> 1 -> 5 -> 7 -> 1 -> ... 这样的车站路线行驶。
// 现在从 source 车站出发（初始时不在公交车上），要前往 target 车站。 期间仅可乘坐公交车。
// 求出 最少乘坐的公交车数量 。如果不可能到达终点车站，返回 -1 。

type Station = number
type Bus = number

function numBusesToDestination(routes: number[][], source: number, target: number): number {
  // `每个车站可以乘坐的公交车`
  const busByStation = new Map<Station, Set<Bus>>()
  routes.forEach((route, bus) =>
    route.forEach(station => {
      !busByStation.has(station) && busByStation.set(station, new Set())
      busByStation.get(station)!.add(bus)
    })
  )

  // 已经到达过的车站和已经乘坐过的公交线路不用在遍历了；
  const visitedStation = new Set<number>()
  const visitedBus = new Set<number>()
  let queue: [curStation: number, steps: number][] = [[source, 0]]

  while (queue.length > 0) {
    const nextQueue: [curStation: number, steps: number][] = []
    const len = queue.length
    for (let _ = 0; _ < len; _++) {
      const [curStation, steps] = queue.pop()!
      if (curStation === target) return steps

      for (const nextBus of busByStation.get(curStation) || []) {
        if (visitedBus.has(nextBus)) continue
        visitedBus.add(nextBus)
        for (const nextStation of routes[nextBus]) {
          if (visitedStation.has(nextStation)) continue
          visitedStation.add(nextStation)
          nextQueue.push([nextStation, steps + 1])
        }
      }
    }

    queue = nextQueue
  }

  return -1
}

console.log(
  numBusesToDestination(
    [
      [1, 2, 7],
      [3, 6, 7]
    ],
    1,
    6
  )
)
// 输入：routes = [[1,2,7],[3,6,7]], source = 1, target = 6
// 输出：2
// 解释：最优策略是先乘坐第一辆公交车到达车站 7 , 然后换乘第二辆公交车到车站 6 。
export {}

// 1.
// 为了获取nextStation需要构建 stationToRoute 映射
// 利用题目自带的routeToStation 来获取下一个station
// 2.
// 已经看过的车站和路线必定不在最短路径中重复
