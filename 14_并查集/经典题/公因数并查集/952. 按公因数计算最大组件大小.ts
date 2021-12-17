// from typing import List

import { useUnionFindArray } from '../../推荐使用并查集精简版'

// # 只有当 A[i] 和 A[j] 共用一个大于 1 的公因数时，A[i] 和 A[j] 之间才有一条边。
// # 返回图中最大连通组件的大小。

// print(Solution().largestComponentSize([2, 3, 6, 7, 4, 12, 21, 39]))
// # 输出：8

// 1 <= A.length <= 20000
// 1 <= A[i] <= 100000

// 将每个数和自己的所有(质)因子进行合并=>看有哪些门派
function largestComponentSize(nums: number[]): number {
  const n = Math.max(...nums)
  const uf = useUnionFindArray(n + 1)
  for (const num of nums) {
    for (let factor = 2; factor <= ~~Math.sqrt(num); factor++) {
      if (num % factor === 0) {
        uf.union(num, factor)
        uf.union(num, num / factor)
      }
    }
  }

  const counter = new Map<number, number>()
  for (const num of nums) {
    const root = uf.find(num)
    counter.set(root, (counter.get(root) || 0) + 1)
  }

  return Math.max(...counter.values())
}

console.log(largestComponentSize([2, 3, 6, 7, 4, 12, 21, 39]))
// 输出：8
