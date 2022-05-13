import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297.二叉树的序列化与反序列化'

/**
 * @param voyage 表示 预期 的二叉树 先序遍历 结果。
 * @description 翻转 最少 的树中节点，使二叉树的 先序遍历 与预期的遍历行程 voyage 相匹配 。
 * 思路: 前序遍历，同时用cur指向voyage数组，遍历访问的时候判断数字是否对的上，对不上尝试左右调换。最后判断cur是否走到最后就好了。
 */
const flipMatchVoyage = function (root: BinaryTree | null, voyage: number[]): number[] {
  if (root?.val !== voyage[0]) return [-1]
  const res: number[] = []
  let p = 0

  // 前序遍历多了一个判断当前值想不想等
  const dfs = (root: BinaryTree) => {
    if (root.val === voyage[p]) {
      p++
      if (root.left && root.left.val === voyage[p]) {
        dfs(root.left)
      }
      if (root.right && root.right.val === voyage[p]) {
        dfs(root.right)
        if (root.left && root.left.val === voyage[p]) {
          res.push(root.val)
          dfs(root.left)
        }
      }
    }
  }
  dfs(root)

  return p === voyage.length ? res : [-1]
}

console.dir(flipMatchVoyage(deserializeNode([1, 2, 3])!, [1, 2, 3]), {
  depth: null,
})

export {}
