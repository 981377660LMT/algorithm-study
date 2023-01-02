// todo
// dfs遍历树，形成dfs数组。子树的dfs序是连续的，该题就变成：
// !在数组中查询若干区间的正整数mex。

// 执行用时：
// 1448 ms
// 内存消耗：
// 120.5 MB

import { useDfsOrder } from '../../../6_tree/树的性质/dfs序/useDfsOrder'
import { useMoAlgo, WindowManager } from './useMoAlgo'

// 总共有 1e5 个基因值，基因值 互不相同,每个基因值都用 闭区间 [1, 1e5] 中的一个整数表示
function smallestMissingValueSubtree(parents: number[], nums: number[]): number[] {
  const n = nums.length
  const adjList = Array.from<unknown, number[]>({ length: n }, () => [])
  parents.forEach((pre, cur) => {
    if (pre === -1) return
    adjList[pre].push(cur)
    adjList[cur].push(pre)
  })

  const { queryRange, queryId } = useDfsOrder(n, adjList)
  const newNums = new Uint32Array(n)

  for (let root = 0; root < n; root++) {
    const dfsId = queryId(root)
    newNums[dfsId - 1] = nums[root]
  }

  // !莫队算法求区间mex
  let mex = 1
  const counter = new Map<number, number>()
  const windowManager: WindowManager<number> = {
    add(index) {
      const num = newNums[index]
      counter.set(num, (counter.get(num) || 0) + 1)
      while ((counter.get(mex) || 0) > 0) {
        mex++
      }
    },
    remove(index) {
      const num = newNums[index]
      counter.set(num, (counter.get(num) || 0) - 1)
      if ((counter.get(num) || 0) === 0) mex = Math.min(mex, num)
    },
    query() {
      return mex
    }
  }
  const queryMex = useMoAlgo(n, n, windowManager)

  for (let root = 0; root < n; root++) {
    let [left, right] = queryRange(root)
    left--, right--
    queryMex.addQuery(left, right)
  }

  return queryMex.work()
}

if (require.main === module) {
  console.log(smallestMissingValueSubtree([-1, 0, 0, 2], [1, 2, 3, 4]))
  console.log(smallestMissingValueSubtree([-1, 0, 1, 1, 1], [5, 4, 1, 3, 2]))
  console.log(smallestMissingValueSubtree([-1, 2, 3, 0, 2, 4, 1], [2, 3, 4, 5, 6, 7, 8]))
}
