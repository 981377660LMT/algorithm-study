/* eslint-disable no-restricted-syntax */
// todo
// dfs遍历树，形成dfs数组。子树的dfs序是连续的，该题就变成：
// !在数组中查询若干区间的正整数mex。

// 执行用时：
// 1448 ms
// 内存消耗：
// 120.5 MB

import { useDfsOrder } from '../../../6_tree/树的性质/dfs序/useDfsOrder'

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
  const queries: [left: number, right: number, root: number][] = []

  for (let root = 0; root < n; root++) {
    const dfsId = queryId(root)
    newNums[dfsId - 1] = nums[root]
    let [left, right] = queryRange(root)
    left--, right--
    queries.push([left, right + 1, root]) // !注意这里right+1
  }

  return moAlgo(newNums, queries)

  function moAlgo(
    arr: Record<number, number>,
    Q: [left: number, right: number, qid: number][]
  ): number[] {
    const res = Array<number>(n).fill(1)
    const chunkSize = Math.max(1, Math.floor(Math.sqrt(n)))
    // !以询问左端点所在的分块的序号为第一关键字，右端点的大小为第二关键字进行排序
    Q.sort(
      (q1, q2) => Math.floor(q1[0] / chunkSize) - Math.floor(q2[0] / chunkSize) || q1[1] - q2[1]
    )

    let [left, right] = [0, 0]
    let mex = 1
    const counter = new Map<number, number>()

    for (const [qLeft, qRight, qIndex] of Q) {
      // !窗口收缩
      while (right > qRight) {
        right--
        remove(arr[right])
      }

      while (left < qLeft) {
        remove(arr[left])
        left++
      }

      // !窗口扩张
      while (right < qRight) {
        add(arr[right])
        right++
      }

      while (left > qLeft) {
        left--
        add(arr[left])
      }

      res[qIndex] = mex
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
