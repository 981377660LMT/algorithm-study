// 请你返回你可以参加的 最大 会议数目。
// 1 <= events.length <= 10^5
// [1,2]和 [1,2]可以一起算
// 你可以在满足 startDayi <= d <= endDayi 中的任意一天 d 参加会议 i
// 一天最多参加一个，可以不完整参加会议

import { MinHeap } from '../../../2_queue/minheap'

/**
 *
 * @param events
 * @summary
 * 1. 按会议开始顺序排序
 * 2. 按结束时间加入pq
 * 3. 删去已经结束的会议
 * 4. 加入结束最早的会议
 */
function maxEvents(events: number[][]): number {
  const n = events.length
  let eventIndex = 0
  let res = 0

  events.sort((a, b) => a[0] - b[0])
  const pq = new MinHeap<number>()

  for (let curDay = 1; curDay <= 1e5; curDay++) {
    // 当日开始的会议
    while (eventIndex < n && events[eventIndex][0] === curDay) {
      pq.push(events[eventIndex++][1])
    }

    // 删去已经结束的会议
    while (pq.size > 0 && pq.peek() < curDay) {
      pq.shift()
    }

    // 加入结束最早的会议
    if (pq.size > 0) {
      pq.shift()
      res++
    }
  }

  return res
}

console.log(
  maxEvents([
    [1, 2],
    [2, 3],
    [3, 4],
  ])
)
console.log(
  maxEvents([
    [1, 4],
    [4, 4],
    [2, 2],
    [3, 4],
    [1, 1],
  ])
)

// 注意[1,2] [1,2]可以参加 前一个在d=1参加 后一个在d=2参加
console.log(
  maxEvents([
    [1, 2],
    [2, 3],
    [3, 4],
    [1, 2],
  ])
)

export {}
// console.log(1e5)
