import { BinaryTree } from '../Tree'

// 偶数下标 层上的所有节点的值都是 奇 整数，从左到右按顺序 严格递增
// 奇数下标 层上的所有节点的值都是 偶 整数，从左到右按顺序 严格递减
function isEvenOddTree(root: BinaryTree | null): boolean {
  if (!root) return false

  let isEven = true
  let queue = [root]

  while (queue.length > 0) {
    const len = queue.length
    let pre: number | null = null

    for (let i = 0; i < len; i++) {
      const cur = queue.shift()!

      if (isEven) {
        if ((cur.val & 1) === 0) return false
        if (pre != void 0 && pre >= cur.val) return false
      } else {
        if ((cur.val & 1) === 1) return false
        if (pre != void 0 && pre <= cur.val) return false
      }

      pre = cur.val
      cur.left && queue.push(cur.left)
      cur.right && queue.push(cur.right)
    }

    isEven = !isEven
  }

  return true
}
