// from typing import List, Tuple, Optional
// from collections import defaultdict, Counter, deque
// from sortedcontainers import SortedList

// MOD = int(1e9 + 7)
// INF = int(1e20)

// # 给你一个整数 n 和一个在范围 [0, n - 1] 以内的整数 p ，它们表示一个长度为 n 且下标从 0 开始的数组 arr ，数组中除了下标为 p 处是 1 以外，其他所有数都是 0 。

// # 同时给你一个整数数组 banned ，它包含数组中的一些位置。banned 中第 i 个位置表示 arr[banned[i]] = 0 ，题目保证 banned[i] != p 。

// # 你可以对 arr 进行 若干次 操作。一次操作中，你选择大小为 k 的一个 子数组 ，并将它 翻转 。在任何一次翻转操作后，你都需要确保 arr 中唯一的 1 不会到达任何 banned 中的位置。换句话说，arr[banned[i]] 始终 保持 0 。

// # 请你返回一个数组 ans ，对于 [0, n - 1] 之间的任意下标 i ，ans[i] 是将 1 放到位置 i 处的 最少 翻转操作次数，如果无法放到位置 i 处，此数为 -1 。

// # 子数组 指的是一个数组里一段连续 非空 的元素序列。
// # 对于所有的 i ，ans[i] 相互之间独立计算。
// # 将一个数组中的元素 翻转 指的是将数组中的值变成 相反顺序 。
// class Solution:
//     def minReverseOperations(self, n: int, p: int, banned: List[int], k: int) -> List[int]:
//         def getNextPos(cur: int):
//             """反转长度为k的子数组后,cur可以到哪些位置"""
//             for posInArray in range(k):  # 当前子数组中的位置
//                 leftBound = cur - posInArray
//                 rightBound = leftBound + k - 1
//                 cand = leftBound + rightBound - cur
//                 if 0 <= leftBound < n and 0 <= rightBound < n:
//                     if not isBanned[cand]:
//                         yield cand
//                 else:
//                     break

//         if k == 1:
//             res = [-1] * n
//             res[p] = 0
//             return res

//         isBanned = [False] * n
//         for i in banned:
//             isBanned[i] = True

//         # 1可以到哪些位置
//         # bfs??
//         visited = [False] * n
//         queue = deque([(0, p)])
//         dist = [INF] * n
//         visited[p] = True
//         dist[p] = 0
//         while queue:
//             curDist, cur = queue.popleft()
//             if curDist > dist[cur]:
//                 continue
//             for next_ in getNextPos(cur):
//                 if 0 <= next_ < n and not visited[next_] and not isBanned[next_]:
//                     visited[next_] = True
//                     queue.append((curDist + 1, next_))
//                     dist[next_] = curDist + 1

//         return [-1 if dist[i] == INF else dist[i] for i in range(n)]

// # n = 5, p = 0, banned = [2,4], k = 3
// print(Solution().minReverseOperations(5, 0, [2, 4], 3))

const INF = 1 << 30
function minReverseOperations(n: number, p: number, banned: number[], k: number): number[] {
  if (k === 1) {
    const res = Array(n).fill(-1)
    res[p] = 0
    return res
  }
  const isBanned = new Uint8Array(n)
  for (let i = 0; i < banned.length; i++) {
    isBanned[banned[i]] = 1
  }
  const visited = new Uint8Array(n)
  const dist = new Uint32Array(n)
  for (let i = 0; i < n; i++) {
    dist[i] = INF
  }
  visited[p] = 1
  dist[p] = 0
  let queue: [number, number][] = [[0, p]]

  while (queue.length) {
    const nextQueue: [number, number][] = []
    for (let i = 0; i < queue.length; i++) {
      const curDist = queue[i][0]
      const cur = queue[i][1]
      if (curDist > dist[cur]) {
        continue
      }
      for (let posInArray = 0; posInArray < Math.min(k, 1000); posInArray++) {
        const leftBound = cur - posInArray
        const rightBound = leftBound + k - 1
        const cand = leftBound + rightBound - cur
        if (leftBound >= 0 && rightBound < n) {
          if (!isBanned[cand] && !visited[cand]) {
            //  cand
            visited[cand] = 1
            nextQueue.push([curDist + 1, cand])
            dist[cand] = curDist + 1
          }
        }
      }
    }
    queue = nextQueue
  }

  const res = Array(n).fill(-1)
  for (let i = 0; i < n; i++) {
    if (dist[i] !== INF) {
      res[i] = dist[i]
    }
  }
  return res
}

// 2
// 1
// []
// 2
// [1,0]
console.log(minReverseOperations(2, 1, [], 2))
