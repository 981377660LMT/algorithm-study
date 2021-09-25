// import { MinHeap } from '../../../2_queue/minheap'

// type Cur = number
// type Weight = number
// type Level = number
// type Edge = [Cur, Weight, Level]

// /**
//  * @param {number} n
//  * @param {number[][]} flights
//  * @param {number} src
//  * @param {number} dst
//  * @param {number} k  最多经过 k 站中转的路线
//  * @return {number}
//  * 找到出一条最多经过 k 站中转的路线，使得从 src 到 dst 的 价格最便宜 ，
//  * 并返回该价格。 如果不存在这样的路线，则输出 -1。
//  * @summary
//  * 带限制的最短路径
//  */
// const findCheapestPrice = function (
//   n: number,
//   flights: number[][],
//   src: number,
//   dst: number,
//   k: number
// ): number {
//   const adjList = Array.from<number, [number, number][]>({ length: n }, () => [])
//   flights.forEach(([u, v, w]) => adjList[u].push([v, w]))
//   const dist = Array<number>(n).fill(Infinity)
//   dist[src] = 0
//   const levels = Array<number>(n).fill(0)

//   const compareFunction = (a: Edge, b: Edge) => a[1] - b[1]
//   const queue = new MinHeap<Edge>(compareFunction)
//   queue.push([src, 0, 0])

//   while (queue.size) {
//     // 1.每次都从离原点最近的没更新过的点开始更新(性能瓶颈：可使用优先队列优化成ElogE)
//     const [cur, weightSum, level] = queue.shift()!
//     if (level > k + 1) continue
//     if (cur === dst) return weightSum
//     // 2.加入visited
//     levels[cur] = level

//     // 3.利用cur点来更新其相邻节点next与原点的距离
//     for (const [next, weight] of adjList[cur]) {
//       if (dist[cur] + weight < dist[next]) {
//         dist[next] = dist[cur] + weight
//         queue.push([next, dist[next], level + 1])
//       }
//     }
//   }

//   return -1
// }

// // console.log(
// //   findCheapestPrice(
// //     3,
// //     [
// //       [0, 1, 100],
// //       [1, 2, 100],
// //       [0, 2, 500],
// //     ],
// //     0,
// //     2,
// //     1
// //   )
// // )

// console.log(
//   findCheapestPrice(
//     3,
//     [
//       [0, 1, 100],
//       [1, 2, 100],
//       [0, 2, 500],
//     ],
//     0,
//     2,
//     0
//   )
// )
// export {}
// 有问题·
