// # 返回包含 N 个结点的所有可能满二叉树的列表。 答案的每个元素都是一个可能树的根结点。

import { BinaryTree } from '../力扣加加/Tree'

// # 答案中每个树的每个结点都必须有 node.val=0。
function allPossibleFBT(n: number): (BinaryTree | null)[] {
  if ((n & 1) === 0) return []

  return [...gen(n)]

  function* gen(n: number): Generator<BinaryTree | null> {
    if (n === 1) yield new BinaryTree(0)
    for (let i = 1; i < Math.min(20, n); i += 2) {
      for (const left of gen(i)) {
        for (const right of gen(n - 1 - i)) {
          const root = new BinaryTree(0)
          root.left = left
          root.right = right
          yield root
        }
      }
    }
  }
}

console.log(allPossibleFBT(7))
