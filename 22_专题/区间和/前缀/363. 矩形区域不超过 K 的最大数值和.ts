// 找出并返回矩阵内部矩形区域的不超过 k 的最大数值和。
// 核心就是怎么把O(n^4)降到O(n^3)，
// 方法是用上下定限O(n^2)来降维，
// 然后在O(n)或O(nlogn)时间内解决降维后的子问题  定三移一的思路
// 矩阵其实可能为负数，因此不满足单调性。这里我们可以手动维护 pres 单调递增，

import { bisectLeft } from '../../../9_排序和搜索/二分api/7_二分搜索寻找最左插入位置'
import { bisectInsort } from '../../../9_排序和搜索/二分api/7_二分搜索插入元素'

// 这样就可以使用二分法在 $logN$ 的时间求出以当前项 i 结尾的不大于 k 的最大矩形和
// 1 <= m, n <= 100
function maxSumSubmatrix(matrix: number[][], k: number): number {
  const m = matrix.length
  const n = matrix[0].length
  let res = -Infinity
  const pre = Array.from<any, number[]>({ length: m + 1 }, () => Array(n + 1).fill(0))
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      pre[i][j] = matrix[i - 1][j - 1] + pre[i - 1][j] + pre[i][j - 1] - pre[i - 1][j - 1]
    }
  }

  // 固定上下
  for (let top = 1; top <= m; top++) {
    for (let bottom = top; bottom <= m; bottom++) {
      const treeSet = [0]
      // 遍历子矩阵的右边界
      for (let right = 1; right <= n; right++) {
        const sum = pre[bottom][right] - pre[top - 1][right]
        const leftIndex = bisectLeft(treeSet, sum - k)
        if (leftIndex < treeSet.length) {
          res = Math.max(res, sum - treeSet[leftIndex])
        }
        bisectInsort(treeSet, sum)
      }
    }
  }

  // console.table(pre)

  return res === -Infinity ? -1 : res
}

console.log(
  maxSumSubmatrix(
    [
      [1, 0, 1],
      [0, -2, 3],
    ],
    2
  )
)

// 时间复杂度：$O(m^ 2 * n )$
// 空间复杂度：$O(m)$

// 时间复杂度
// 进阶：如果行数远大于列数，该如何设计解决方案？
// 我们可以将行列兑换 复杂度为 O(m * n^2)
// 空间复杂度
// 我们可以将计算前缀和的逻辑下放到搜索子矩阵的循环里去做，
// 从而将 O(m * n)O(m∗n) 的空间复杂度下降到 O(\max(m,n))O(max(m,n))。
