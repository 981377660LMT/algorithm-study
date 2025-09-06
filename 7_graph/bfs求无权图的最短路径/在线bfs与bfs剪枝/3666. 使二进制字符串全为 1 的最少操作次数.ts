// 类似 2612. 最少翻转操作数

import { onlineBfs } from '../../../22_专题/implicit_graph/OnlineBfs-在线bfs'
import { Finder } from '../../../22_专题/implicit_graph/RangeFinder/Finder-fastset'

const INF = 2e15

// 3666. 使二进制字符串全为 1 的最少操作次数
// https://leetcode.cn/problems/minimum-operations-to-equalize-binary-string/description/
// 给你一个二进制字符串 s 和一个整数 k。
// 在一次操作中，你必须选择 恰好 k 个 不同的 下标，并将每个 '0' 翻转 为 '1'，每个 '1' 翻转为 '0'。
// 返回使字符串中所有字符都等于 '1' 所需的 最少 操作次数。如果不可能，则返回 -1。
// https://leetcode.cn/problems/minimum-operations-to-equalize-binary-string/solutions/3767987/bfs-you-xu-ji-he-by-tsreaper-1isk/
function minOperations(s: string, k: number): number {
  const n = s.length
  let zeros = 0
  for (let i = 0; i < n; i++) {
    if (s[i] === '0') {
      zeros++
    }
  }

  if (zeros === 0) {
    return 0
  }

  const finder = [new Finder(n + 1), new Finder(n + 1)]
  for (let i = 0; i <= n; i++) {
    if ((i & 1) === 0) {
      finder[1].erase(i)
    } else {
      finder[0].erase(i)
    }
  }

  // 从有 i 个 '0' 的状态，可以转移到的邻接点(新的 '0' 的数量)的范围 [min, max]。
  const getRange = (i: number): [number, number] => {
    const minZ = Math.max(0, k - (n - i))
    const maxZ = Math.min(i, k)
    const minNextZeros = i + k - 2 * maxZ
    const maxNextZeros = i + k - 2 * minZ
    return [minNextZeros, maxNextZeros]
  }

  // 将节点 u 标记为已访问。
  const setUsed = (u: number) => {
    finder[u & 1].erase(u)
  }

  // 从节点 u 出发，找到一个未访问的邻居。
  const findUnused = (u: number): number | null => {
    const [left, right] = getRange(u)
    const targetParity = (u & 1) ^ (k & 1)
    const next = finder[targetParity].next(left)
    if (next != null && next <= right) {
      return next
    }
    return null
  }

  const [dist] = onlineBfs(n + 1, zeros, setUsed, findUnused)
  return dist[0] === INF ? -1 : dist[0]
}
