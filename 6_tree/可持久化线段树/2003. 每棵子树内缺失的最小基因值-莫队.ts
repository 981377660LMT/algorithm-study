// todo
// dfs遍历树，形成dfs数组。子树的dfs序是连续的，该题就变成：
// 在数组中查询若干区间的mex。
// 主席树找mex的思路
// 利用权值线段树在每个权值上记录`该数最后出现的下标`，再次基础上加上可持续化，便是主席树了

// 执行用时：
// 1448 ms
// 内存消耗：
// 120.5 MB

import { useDfsOrder } from '../树的性质/dfs序/useDfsOrder'
// 总共有 1e5 个基因值，基因值 互不相同,每个基因值都用 闭区间 [1, 1e5] 中的一个整数表示
function smallestMissingValueSubtree(parents: number[], nums: number[]): number[] {
  const n = nums.length
  const adjList = Array.from<any, number[]>({ length: n }, () => [])
  for (const [cur, pre] of parents.entries()) {
    if (pre === -1) continue
    adjList[pre].push(cur)
    adjList[cur].push(pre)
  }

  const { queryRange, queryId } = useDfsOrder(n, adjList)
  const newNums = new Uint32Array(n).fill(0)
  const queries: [left: number, right: number, root: number][] = []

  for (let root = 0; root < n; root++) {
    const dfsId = queryId(root)
    newNums[dfsId - 1] = nums[root]
    const [left, right] = queryRange(root)
    queries.push([left - 1, right - 1, root])
  }

  return moAlgo(newNums, queries)

  function moAlgo(
    nums: Record<number, number>,
    queries: [left: number, right: number, qid: number][]
  ): number[] {
    const res = Array<number>(n).fill(1)
    const chunkLen = Math.max(1, Math.floor(Math.sqrt(n)))
    // 莫队 左端点按块id排序，右端点按编号排序
    queries.sort(
      (q1, q2) => Math.floor(q1[0] / chunkLen) - Math.floor(q2[0] / chunkLen) || q1[1] - q2[1]
    )

    let [left, right, mex] = [0, 0, nums[0] === 1 ? 2 : 1]
    const counter = new Map<number, number>([[nums[0], 1]])

    for (const [qLeft, qRight, qid] of queries) {
      while (left < qLeft) {
        remove(nums[left])
        left++
      }
      while (left > qLeft) {
        add(nums[left - 1])
        left--
      }
      while (right < qRight) {
        add(nums[right + 1])
        right++
      }
      while (right > qRight) {
        remove(nums[right])
        right--
      }

      res[qid] = mex
    }

    return res

    function add(num: number): void {
      // 加一个数 mex增加到缺失
      counter.set(num, (counter.get(num) ?? 0) + 1)
      while ((counter.get(mex) ?? 0) > 0) mex++
    }

    function remove(num: number): void {
      // 减一个数，需要看mex是否变小了
      counter.set(num, counter.get(num)! - 1)
      if (counter.get(num) === 0) mex = Math.min(mex, num) // 少了一个数
    }
  }
}

if (require.main === module) {
  console.log(smallestMissingValueSubtree([-1, 0, 0, 2], [1, 2, 3, 4]))
  console.log(smallestMissingValueSubtree([-1, 2, 3, 0, 2, 4, 1], [2, 3, 4, 5, 6, 7, 8]))
  console.log(smallestMissingValueSubtree([-1, 0, 1, 1, 1], [5, 4, 1, 3, 2]))
}
