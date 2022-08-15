/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { isBipartite } from './二分图检测'

/**
 *
 * @param n
 * @param dislikes
 * @description
 * 给定一组 N 人（编号为 1, 2, ..., N）， 我们想把每个人分进任意大小的两组。
 * 每个人都可能不喜欢其他人，那么他们不应该属于同一组。
 * 当可以用这种方法将每个人分进两组时，返回 true；否则返回 false。
 * @summary
 * @link https://leetcode-cn.com/problems/possible-bipartition/solution/dfs-jin-xing-er-fen-tu-ran-se-wo-lai-gei-l2p3/
 * 👆节省空间的做法
 * 考虑由给定的 “不喜欢” 边缘形成的 N 人的图表。我们要检查这个图的每个连通分支是否为二分的。
 */
function possibleBipartition(n: number, dislikes: number[][]): boolean {
  // 邻接表
  const adjMap = new Map<number, Set<number>>()
  for (const [a, b] of dislikes) {
    !adjMap.has(a) && adjMap.set(a, new Set())
    !adjMap.has(b) && adjMap.set(b, new Set())
    adjMap.get(a)!.add(b)
    adjMap.get(b)!.add(a)
  }

  return isBipartite(adjMap)
}

if (require.main === module) {
  console.log(
    possibleBipartition(4, [
      [1, 2],
      [1, 3],
      [2, 4]
    ])
  )
}

export {}
