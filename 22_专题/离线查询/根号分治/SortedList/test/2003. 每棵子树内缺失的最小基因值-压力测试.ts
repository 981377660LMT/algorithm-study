/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable generator-star-spacing */

import { SortedListFast } from '../SortedListFast'

// https://leetcode.cn/problems/smallest-missing-genetic-value-in-each-subtree/submissions/
function smallestMissingValueSubtree(parents: number[], nums: number[]): number[] {
  const n = parents.length
  const adjList: number[][] = Array.from({ length: n }, () => [])
  for (let i = 1; i < n; i++) {
    adjList[parents[i]].push(i)
  }
  const res = Array(n).fill(1)
  dfs(0, -1)
  return res

  function findMex(sl: SortedListFast<number>): number {
    let left = 0
    let right = sl.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      const diff = sl.at(mid)! - (mid + 1)
      if (diff >= 1) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left + 1
  }

  function dfs(cur: number, parent: number): SortedListFast<number> {
    let curTree = new SortedListFast<number>()
    for (const next of adjList[cur]) {
      if (next === parent) {
        continue
      }
      const subTree = dfs(next, cur)
      let [big, small] = [curTree, subTree]
      if (big.length < small.length) {
        ;[big, small] = [small, big]
      }
      small.forEach(v => {
        big.add(v)
      })
      curTree = big
    }

    curTree.add(nums[cur])
    res[cur] = findMex(curTree)
    return curTree
  }
}
