import { BinaryTree } from '../Tree'

function findSecondMinimumValue(root: BinaryTree | null): number {
  if (!root) return -1
  const queue = [root]
  const res = new Set<number>()
  while (queue.length) {
    const cur = queue.pop()!
    res.add(cur.val)
    cur.left && queue.push(cur.left)
    cur.right && queue.push(cur.right)
  }
  if (res.size <= 1) return -1
  return [...res].sort((a, b) => a - b)[1]
}
