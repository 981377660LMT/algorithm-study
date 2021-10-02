import { BinaryTree } from '../力扣加加/Tree'

// left right换个位置即可
function kthLargest(root: BinaryTree | null, k: number): number {
  function* inOrder(root: BinaryTree | null): Generator<number> {
    if (!root) return
    yield* inOrder(root.right)
    yield root.val
    yield* inOrder(root.left)
  }

  const gen = inOrder(root)
  for (let i = 0; i < k - 1; i++) {
    gen.next()
  }

  return gen.next().value
}
