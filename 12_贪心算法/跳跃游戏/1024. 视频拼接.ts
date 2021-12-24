import { maxJump } from './maxJump'

/**
 * @param {number[][]} clips  1 <= clips.length <= 100
 * @param {number} time  这些片段来自于一项持续时长为 time 秒的体育赛事
 * @return {number}
 * 返回所需片段的最小数目，如果无法完成该任务，则返回 -1 。
 * @summary 跳跃游戏2
 * 跳跃游戏就相当于一个将起点排序的区间问题
 */
var videoStitching = function (clips: number[][], time: number): number {
  const jumps = Array<number>(time + 1).fill(0)
  for (const [index, value] of clips) {
    jumps[index] = Math.max(jumps[index], value)
  }

  const res = maxJump(jumps, time)
  return res
}

console.log(
  videoStitching(
    [
      [0, 2],
      [4, 6],
      [8, 10],
      [1, 9],
      [1, 5],
      [5, 9],
    ],
    10
  )
)
// console.log(
//   videoStitching(
//     [
//       [0, 1],
//       [1, 2],
//     ],
//     5
//   )
// )
console.log(
  videoStitching(
    [
      [0, 4],
      [2, 8],
    ],
    5
  )
)
