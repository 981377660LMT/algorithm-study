import { BinaryTree } from '../6_tree/分类/Tree'

/**
 * @param {number} n
 * @return {TreeNode[]}
 * 请你生成并返回所有由 n 个节点组成且节点值从 1 到 n 互不相同的不同 二叉搜索树
 */
const generateTrees = (n: number): (BinaryTree | null)[] => {
  // 所有数字都有可能作为根，因此遍历 num 作为根
  // 根的左子树由比根小的数字构成，递归 num[:i] 就是左子树所有的可能结构，同理可获得右子树所有可能的结构
  function* gen(first: number, last: number): Generator<BinaryTree | null> {
    if (first > last) yield null
    for (let root = first; root <= last; root++) {
      for (const left of gen(first, root - 1)) {
        for (const right of gen(root + 1, last)) {
          yield new BinaryTree(root, left, right)
        }
      }
    }
  }

  return [...gen(1, n)]
}

console.dir(generateTrees(1), { depth: null })
export default 1
