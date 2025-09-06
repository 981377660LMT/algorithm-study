import { onlineBfs } from '../../../22_专题/implicit_graph/OnlineBfs-在线bfs'
import { Finder } from '../../../22_专题/implicit_graph/RangeFinder/Finder-fastset'

const INF = 2e15

// 2612. 最少翻转操作数
// https://leetcode.cn/problems/minimum-reverse-operations/description/
// 给定一个整数 n 和一个整数 p，它们表示一个长度为 n 且除了下标为 p 处是 1 以外，其他所有数都是 0 的数组 arr。
// 同时给定一个整数数组 banned ，它包含数组中的一些限制位置。在 arr 上进行下列操作：
// 如果单个 1 不在 banned 中的位置上，反转大小为 k 的 子数组。
// 返回一个包含 n 个结果的整数数组 answer，其中第 i 个结果是将 1 放到位置 i 处所需的 最少 翻转操作次数，如果无法放到位置 i 处，此数为 -1 。
function minReverseOperations(n: number, p: number, banned: number[], k: number): number[] {
  const finder = [new Finder(n), new Finder(n)]
  for (let i = 0; i < n; i++) {
    finder[(i & 1) ^ 1].erase(i)
  }
  banned.forEach(i => {
    finder[i & 1].erase(i)
  })

  // 反转长度为k的子数组,从i可以到达的左右边界闭区间
  const getRange = (i: number): [number, number] => [
    Math.max(i - k + 1, k - i - 1),
    Math.min(i + k - 1, 2 * n - k - i - 1)
  ]

  // 将u位置标记为已经访问过.
  const setUsed = (u: number) => {
    finder[u & 1].erase(u)
  }

  // 找到一个未使用的位置.如果不存在,返回null.
  const findUnused = (u: number): number | null => {
    const [left, right] = getRange(u)
    const next = finder[(u + k + 1) & 1].next(left)
    if (next != null && left <= next && next <= right) {
      return next
    }
    return null
  }

  const dist = onlineBfs(n, p, setUsed, findUnused)[0]
  return dist.map(d => (d === INF ? -1 : d))
}
