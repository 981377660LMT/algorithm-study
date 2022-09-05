/* eslint-disable prefer-destructuring */
class TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
  constructor(val?: number, left?: TreeNode | null, right?: TreeNode | null) {
    this.val = val === undefined ? 0 : val
    this.left = left === undefined ? null : left
    this.right = right === undefined ? null : right
  }
}

function levelOrder(root: TreeNode | null): number[][] {
  if (!root) return []
  const queue: TreeNode[] = [root]
  const res: number[][] = []

  while (queue.length) {
    const level: number[] = []
    const step = queue.length
    for (let _ = 0; _ < step; _++) {
      const cur = queue.shift()! // eslint-disable-line @typescript-eslint/no-non-null-assertion
      level.push(cur.val)
      cur.left && queue.push(cur.left)
      cur.right && queue.push(cur.right)
    }
    res.push(level)
  }

  return res
}

export {}
