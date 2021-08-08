import { BinaryTree, bt } from '../Tree'

// 选择二叉树中 任意 节点和一个方向（左或者右）。
// 如果前进方向为右，那么移动到当前节点的的右子节点，否则移动到它的左子节点。
// const longestZigZag = (root: BinaryTree): number => {
//   if (!root) return 0
//   let max = 0

//   const dfs = (root: BinaryTree, direction: 0 | 1, depth: number) => {
//     root.right && direction === 0 && dfs(root.right, 1, depth + 1)
//     root.left && direction === 1 && dfs(root.left, 0, depth + 1)
//     max = Math.max(max, depth)
//   }
//   root.left && dfs(root.left, 0, 1)
//   root.right && dfs(root.right, 1, 1)

//   return max
// }

// 注意起点是任意的
// 要用weakMap缓存
const longestZigZag = (root: BinaryTree): number => {
  if (!root) return 0

  const dfs = (root: BinaryTree | null, isLeft: boolean, len: number): number => {
    if (!root) return len
    if (isLeft) {
      return Math.max(dfs(root.right, false, len + 1), dfs(root.left, true, 0))
    } else {
      return Math.max(dfs(root.left, true, len + 1), dfs(root.right, false, 0))
    }
  }
  return Math.max(dfs(root.left, true, 0), dfs(root.right, false, 0))
}

console.dir(longestZigZag(bt), { depth: null })

export {}
