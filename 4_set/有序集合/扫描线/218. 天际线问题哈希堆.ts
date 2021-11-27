import { HashHeap } from '../../../8_heap/HashHeap'

type Height = number
type Position = number
/**
 * @param {number[][]} buildings  左,右,高
 * @return {number[][]}
   请返回由这些建筑物形成的 天际线
   @summary
   扫描线算法：
   从左到右扫描边缘线
   左边缘则入堆 右边缘则出堆
   堆内高度变化则遇到了`关键点`
   `关键点`坐标为[当前坐标，堆内最大值]
 */
const getSkyline = function (buildings: number[][]): number[][] {
  const res: number[][] = []
  const points: [Position, Height][] = []
  // 技巧：要使坐标相同时先访问左边缘线=>将建筑物左侧高度处理为负数即可
  for (const building of buildings) {
    points.push([building[0], -building[2]])
    points.push([building[1], building[2]])
  }
  points.sort((a, b) => a[0] - b[0] || a[1] - b[1])

  const heightQueue = new HashHeap((a, b) => b - a)
  // 防止heightQueue为空 即底部点的高度为0（右下角)
  heightQueue.push(0)

  let preMaxHeight = 0
  for (const point of points) {
    // 左边缘 入  右边缘 出
    if (point[1] < 0) heightQueue.push(-point[1])
    else heightQueue.remove(point[1])

    // 高度发生变化时加入res
    const curMaxHeight = heightQueue.peek()
    if (preMaxHeight !== curMaxHeight) {
      res.push([point[0], curMaxHeight])
      preMaxHeight = curMaxHeight
    }
  }

  return res
}

console.log(
  getSkyline([
    [2, 9, 10],
    [3, 7, 15],
    [5, 12, 12],
    [15, 20, 10],
    [19, 24, 8],
  ])
)
// [[2,10],[3,15],[7,12],[12,0],[15,10],[20,8],[24,0]]
export {}
