import { PriorityQueue } from '../../../../2_queue/PriorityQueue'

/**
 * @param {number} target
 * @param {number} startFuel
 * @param {number[][]} stations
 * @return {number}
 * @description 每个 station[i] 代表一个加油站，它位于出发位置东面 station[i][0] 英里处，并且有 station[i][1] 升汽油
 * 为了到达目的地，汽车所必要的最低加油次数是多少？如果无法到达目的地，则返回 -1 。
 * 什么时候必须加油呢？答案应该是如果你不加油，就无法到达下一个目的地的时候。
 * 
 * 每经过一个站，就将其油量加到堆。
   尽可能往前开，油只要不小于 0 就继续开。
   如果油量小于 0 ，就从堆中取最大的加到油箱中去，如果油量还是小于 0 继续重复取堆中的最大油量。
   如果加完油之后油量大于 0 ，继续开，重复上面的步骤。否则返回 -1，表示无法到达目的地。
 * 现实中你无论如何都无法知道在当前站，我是应该加油还是不加油的，因为信息太少了。
   这个事后诸葛亮体现在我们是等到没油了才去想应该在之前的某个站加油。
 */
const minRefuelStops = function (target: number, startFuel: number, stations: number[][]): number {
  // 终点看作最后一个加油站
  stations.push([target, 0])
  // 最大堆存储
  const pq = new PriorityQueue((a, b) => b - a)
  let res = 0
  let curFuel = startFuel
  let curPosition = 0

  // 到每一个车站，就存起来
  for (let i = 0; i < stations.length; i++) {
    const [stationPosition, stationFuel] = stations[i]
    curFuel -= stationPosition - curPosition

    // 如果不够就取出油
    while (curFuel < 0 && pq.length) {
      curFuel += pq.shift()!
      res++
    }

    // 还是不够油
    if (curFuel < 0) return -1

    curPosition = stationPosition
    pq.push(stationFuel)
  }

  return res
}

console.log(
  minRefuelStops(100, 10, [
    [10, 60],
    [20, 30],
    [30, 30],
    [60, 40],
  ])
)

export {}
