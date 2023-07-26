/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { onlineBfs } from '../../../22_专题/implicit_graph/OnlineBfs-在线bfs'

// 给你一个整数数组 arr ，你一开始在数组的第一个元素处（下标为 0）。
// 每一步，你可以从下标 i 跳到下标 i + 1 、i - 1 或者 j ：

// i + 1 需满足：i + 1 < arr.length
// i - 1 需满足：i - 1 >= 0
// j 需满足：arr[i] == arr[j] 且 i != j
// 请你返回到达数组最后一个元素的下标处所需的 最少操作次数 。

function minJumps(arr: number[]): number {
  const n = arr.length
  const finder = new Map<number, number[]>()
  const visited = new Uint8Array(n)
  arr.forEach((v, i) => {
    !finder.has(v) && finder.set(v, [])
    finder.get(v)!.push(i)
  })

  const dist = onlineBfs(
    n,
    0,
    cur => {
      visited[cur] = 1 // 标记
    },
    cur => {
      if (cur - 1 >= 0 && !visited[cur - 1]) {
        return cur - 1
      }
      if (cur + 1 < n && !visited[cur + 1]) {
        return cur + 1
      }

      const num = arr[cur]
      const nexts = finder.get(num)!
      // 延迟删除
      while (nexts.length && visited[nexts[nexts.length - 1]]) {
        nexts.pop()
      }
      if (nexts.length) {
        return nexts[nexts.length - 1]
      }

      return null
    }
  )[0]

  return dist[n - 1]
}
