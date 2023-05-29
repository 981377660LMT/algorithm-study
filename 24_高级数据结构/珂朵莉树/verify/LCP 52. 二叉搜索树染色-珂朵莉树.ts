// ops[i] = [type, x, y] 表示第 i 次操作为：
// type 等于 0 时，将节点值范围在 [x, y] 的节点均染蓝
// type 等于 1 时，将节点值范围在 [x, y] 的节点均染红

// https://leetcode.cn/problems/QO5KpG/

import { ODT } from '../ODT-fastset'

function getNumber(root: TreeNode | null, ops: number[][]): number {
  const nums: number[] = []
  dfs(root)
  const odt = new ODT(nums.length, -1) // 初始时，所有节点都是蓝色的

  const [_, rank] = sortedSet(nums)
  ops.forEach(([type, x, y]) => {
    const start = rank(x)
    const end = rank(y) + 1
    odt.set(start, end, type)
  })

  let red = 0
  odt.enumerateAll((start, end, color) => {
    if (color === 1) red += end - start
  })
  return red

  function dfs(root: TreeNode | null) {
    if (!root) return
    nums.push(root.val)
    dfs(root.left)
    dfs(root.right)
  }
}

/**
 * 离散化.
 * @returns 给定一个数，返回其在离散化后的数组中的位置 (0 - {@link nums.length})
 */
function sortedSet(nums: number[]): readonly [sorted: number[], getRank: (num: number) => number] {
  const sorted = [...new Set(nums)].sort((a, b) => a - b)
  const getRank = (num: number) => {
    let left = 0
    let right = sorted.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (sorted[mid] >= num) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }
  return [sorted, getRank]
}

export {}
