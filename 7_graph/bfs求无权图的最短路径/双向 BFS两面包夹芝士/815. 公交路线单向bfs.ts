import { useUnionFindArray } from '../../../14_并查集/推荐使用并查集精简版'
import { ArrayDeque } from '../../../2_queue/Deque/ArrayDeque'

type Station = number
type RouteIndex = number
type Cost = number

/**
 * @param {number[][]} routes routes[i] 中的所有值 互不相同
 * @param {number} source
 * @param {number} target  0 <= source, target < 10**6
 * @return {number}  求出 最少乘坐的公交车数量 。如果不可能到达终点车站，返回 -1 。
 */
const numBusesToDestination = function (
  routes: number[][],
  source: number,
  target: number
): number {
  // 每个车站可以乘坐的公交车
  const stationToRoute = new Map<Station, Set<RouteIndex>>()
  routes.forEach((route, routeIndex) =>
    route.forEach(station => {
      !stationToRoute.has(station) && stationToRoute.set(station, new Set())
      stationToRoute.get(station)!.add(routeIndex)
    })
  )

  // 每个公交车线路可以到达的车站
  const routeToStation = routes.map(route => new Set(route))

  // 已经到达过的车站和已经乘坐过的公交线路不用在遍历了；
  const visitedStation = new Set<number>()
  const visitedRoute = new Set<number>()
  const queue = new ArrayDeque<[Station, Cost]>(10 ** 4)
  queue.push([source, 0])

  while (queue.length) {
    const [curStation, cost] = queue.shift()!
    if (curStation === target) return cost

    for (const nextRoute of stationToRoute.get(curStation)!) {
      if (visitedRoute.has(nextRoute)) continue
      visitedRoute.add(nextRoute)

      for (const nextStation of routeToStation[nextRoute]) {
        if (visitedStation.has(nextStation)) continue
        visitedStation.add(nextStation)

        queue.push([nextStation, cost + 1])
      }
    }
  }

  return -1
}

console.log(
  numBusesToDestination(
    [
      [1, 2, 7],
      [3, 6, 7],
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
