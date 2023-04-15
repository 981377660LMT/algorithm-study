import { PriorityQueue } from '../../../../2_queue/PriorityQueue'

/**
 * @param {number[]} heights
 * @param {number} bricks
 * @param {number} ladders
 * @return {number}
 * 如果当前建筑物的高度 大于或等于 下一建筑物的高度，则不需要梯子或砖块
 * 如果当前建筑的高度 小于 下一个建筑的高度，您可以使用 一架梯子 或 (h[i+1] - h[i]) 个砖块
 * 如果以最佳方式使用给定的梯子和砖块，返回你可以到达的最远建筑物的下标（下标 从 0 开始 ）。
 * @summary
 * 先用砖头，再用梯子，等梯子不够用了，我们就要开始事后诸葛亮了，要是前面用砖头就好了
 * 我们将用前面用梯子跨越的建筑物高度差存起来，等到后面梯子用完了，我们将前面被用的梯子“兑换”成砖头继续用
 * 优先兑换高度差大的那次
 */
const furthestBuilding = function (heights: number[], bricks: number, ladders: number): number {
  const len = heights.length - 1
  const pq = new PriorityQueue((a, b) => b - a)

  for (let i = 0; i < len; i++) {
    const diff = heights[i + 1] - heights[i]
    if (diff <= 0) continue

    pq.push(diff)
    bricks -= diff
    while (bricks < 0) {
      if (ladders) {
        ladders--
        bricks += pq.shift()!
      } else {
        return i
      }
    }
  }

  return len
}

console.log(furthestBuilding([4, 12, 2, 7, 3, 18, 20, 3, 19], 10, 2))
// 输出：7

export {}
