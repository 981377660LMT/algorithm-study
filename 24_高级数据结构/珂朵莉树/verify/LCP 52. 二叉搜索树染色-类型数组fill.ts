// ops[i] = [type, x, y] 表示第 i 次操作为：
// type 等于 0 时，将节点值范围在 [x, y] 的节点均染蓝
// type 等于 1 时，将节点值范围在 [x, y] 的节点均染红

// https://leetcode.cn/problems/QO5KpG/

function getNumber(root: TreeNode | null, ops: number[][]): number {
  const nums: number[] = []
  dfs(root)
  const odt = new Int8Array(nums.length).fill(-1) // 初始时，所有节点都是蓝色的

  const [rank] = sortedSet(nums)
  ops.forEach(([type, x, y]) => {
    const start = rank(x)
    const end = rank(y) + 1
    odt.fill(type, start, end)
  })

  let red = 0
  for (let i = 0; i < odt.length; i++) {
    red += +(odt[i] === 1)
  }
  return red

  function dfs(root: TreeNode | null) {
    if (!root) return
    nums.push(root.val)
    dfs(root.left)
    dfs(root.right)
  }
}

/**
 * (松)离散化.
 * @returns
 * rank: 给定一个数,返回它的排名`(0-count)`.
 * count: 离散化(去重)后的元素个数.
 */
function sortedSet(nums: number[]): [rank: (num: number) => number, count: number] {
  const allNums = [...new Set(nums)].sort((a, b) => a - b)
  const rank = (num: number) => {
    let left = 0
    let right = allNums.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (allNums[mid] >= num) {
        right = mid - 1
      } else {
        left = mid + 1
      }
    }
    return left
  }
  return [rank, allNums.length]
}

export {}
