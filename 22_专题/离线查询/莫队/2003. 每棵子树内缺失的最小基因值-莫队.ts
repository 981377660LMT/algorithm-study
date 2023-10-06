// dfs遍历树，形成dfs数组。子树的dfs序是连续的，该题就变成：
// !在数组中查询若干区间的正整数mex。
// 总共有 1e5 个基因值，基因值 互不相同,每个基因值都用 闭区间 [1, 1e5] 中的一个整数表示

// 执行用时：
// 1448 ms
// 内存消耗：
// 120.5 MB

import { DfsOrder } from '../../../6_tree/树的性质/dfs序/DfsOrder'
import { RangeMexQuery } from './RangeMexQuery-离线查询区间mex-线段树'

function smallestMissingValueSubtree(parents: number[], nums: number[]): number[] {
  const n = nums.length
  const adjList: number[][] = Array(n)
  for (let i = 0; i < n; i++) adjList[i] = []
  parents.forEach((pre, cur) => {
    if (pre === -1) return
    adjList[pre].push(cur)
    adjList[cur].push(pre)
  })

  const D = new DfsOrder(n, adjList)
  const newNums = new Uint32Array(n)
  for (let root = 0; root < n; root++) {
    const dfsId = D.queryId(root)
    newNums[dfsId] = nums[root]
  }

  // !莫队算法求区间mex
  // const rangeMex = new RangeMexQueryMo(newNums, n)
  const rangeMex = new RangeMexQuery(newNums)
  for (let root = 0; root < n; root++) {
    let [left, right] = D.queryRange(root)
    rangeMex.addQuery(left, right)
  }

  return rangeMex.run(1)
}

if (require.main === module) {
  console.log(smallestMissingValueSubtree([-1, 0, 0, 2], [1, 2, 3, 4]))
  console.log(smallestMissingValueSubtree([-1, 0, 1, 1, 1], [5, 4, 1, 3, 2]))
  console.log(smallestMissingValueSubtree([-1, 2, 3, 0, 2, 4, 1], [2, 3, 4, 5, 6, 7, 8]))
}
