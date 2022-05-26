import { useUnionFindArray } from '../useUnionFind'

/**
 * @param {string[]} grid
 * @return {number}
 * @description 在由 1 x 1 方格组成的 N x N 网格 grid 中，每个 1 x 1 方块由 /、\ 或空格构成。
 * 返回区域的数目。
 * 并查集
 * @link https://leetcode-cn.com/problems/regions-cut-by-slashes/solution/js-bing-cha-ji-kan-wo-de-wen-zi-tu-jiu-hao-dong-li/
 * @summary
 * 大概把每一个小方块看成如下图
 *      0
      3   1
        2
        
 */
const regionsBySlashes = function (grid: string[]): number {
  const N = grid.length
  const uf = useUnionFindArray(4 * N ** 2)

  for (let i = 0; i < N; i++) {
    for (let j = 0; j < N; j++) {
      const cur = grid[i][j]
      const pos = i * N + j // 索引转换为 0 1 2 3 4 5...

      switch (cur) {
        case ' ':
          uf.union(pos * 4 + 0, pos * 4 + 1)
          uf.union(pos * 4 + 1, pos * 4 + 2)
          uf.union(pos * 4 + 2, pos * 4 + 3)
          uf.union(pos * 4 + 3, pos * 4 + 0)
          break
        case '/':
          uf.union(pos * 4 + 1, pos * 4 + 2)
          uf.union(pos * 4 + 0, pos * 4 + 3)
          break
        case '\\':
          uf.union(pos * 4 + 0, pos * 4 + 1)
          uf.union(pos * 4 + 2, pos * 4 + 3)
          break
      }

      const top = i - 1 >= 0 ? grid[i - 1][j] : null
      const left = j - 1 >= 0 ? grid[i][j - 1] : null

      if (top) {
        // 连接当前方块和上边的方块
        let topPos = pos - N
        uf.union(pos * 4 + 0, topPos * 4 + 2)
      }

      if (left) {
        // 连接当前方块和左边的方块
        let leftPos = pos - 1
        uf.union(pos * 4 + 3, leftPos * 4 + 1)
      }
    }
  }

  return uf.getCount()
}

console.log(regionsBySlashes(['/\\', '\\/']))
// 5

export {}
