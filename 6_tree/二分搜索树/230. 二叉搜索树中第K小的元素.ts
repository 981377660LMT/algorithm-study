import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297.二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @param {number} k
 * @return {number}
 * 时间复杂度：依次遍历前 k 个节点，因此时间复杂度为 O(k)
   空间复杂度：生成器只需要 O(1) 的空间，如果不考虑递归栈所占用的空间，那么复杂度为 O(1)

 */
const kthSmallest = function (root: BinaryTree, k: number): number {
  function* inOrder(root: BinaryTree | null): Generator<number> {
    if (!root) return
    yield* inOrder(root.left)
    yield root.val
    yield* inOrder(root.right)
  }

  const gen = inOrder(root)
  for (let i = 0; i < k - 1; i++) {
    gen.next()
  }

  return gen.next().value
}

console.log(kthSmallest(deserializeNode([5, 3, 6, 2, 4, null, null, 1])!, 3))

export {}
