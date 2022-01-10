import { ArrayDeque } from '../../2_queue/Deque/ArrayDeque'

/**
 * @param {number[]} jump
 * @return {number}
 * 通过按动弹簧，可以选择把小球向右弹射 jump[i] 的距离，或者向左弹射到任意左侧弹簧的位置。
 * 请求出最少需要按动多少次弹簧，可以将小球从编号 0 弹簧弹出整个机器，即向右越过编号 N-1 的弹簧
 */
const minJump = function (jump: number[]): number {
  const visited = new Set()
  const queue = new ArrayDeque(10 ** 6)
  queue.push(0)
  let steps = 0
  let visitedMaxIndex = 0

  while (queue.length) {
    const len = queue.length

    for (let i = 0; i < len; i++) {
      const curIndex = queue.shift()!
      // 可以弹出了
      if (curIndex + jump[curIndex] >= jump.length) return steps + 1

      visited.add(curIndex)

      for (const next of [...range(visitedMaxIndex, curIndex + 1), curIndex + jump[curIndex]]) {
        if (next >= 0 && next < jump.length && !visited.has(next)) {
          queue.push(next)
        }
      }

      // 更新看过的最大位置 剪枝
      visitedMaxIndex = Math.max(curIndex, visitedMaxIndex)
    }

    steps++
  }

  return -1

  function range(start: number, end: number): number[] {
    const res: number[] = []
    for (let i = start; i < end; i++) {
      res.push(i)
    }
    return res
  }
}

console.log(minJump([2, 5, 1, 1, 1, 1]))

// 输出：3
// 解释：小 Z 最少需要按动 3 次弹簧，小球依次到达的顺序为 0 -> 2 -> 1 -> 6，最终小球弹出了机器。

// bfs 超时
// 假设我们jump数组的长度为6
// 当我们走到下标为3的位置，我们需要访问下标为0,1,2的位置
// 当我们走到下边为4的位置，我们需要访问下标为0,1,2,3的位置
// 0,1,2这些位置会在之后的位置被不断地访问，这就是bfs算法超时的原因
// 时间复杂度大概是O(n*n)

// 记录已访问过的最大index，该index左边的弹簧均已被遍历，从而避免重复访问
