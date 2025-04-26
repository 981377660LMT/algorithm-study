type Nullish<T> = T | null | undefined

interface ITreeNode {
  left: Nullish<this>
  right: Nullish<this>
}

/**
 * Morris 中序遍历.
 */
function morris<T extends ITreeNode>(root: T, visit: (node: T) => void): void {
  let cur: Nullish<T> = root
  while (cur) {
    if (cur.left) {
      let pred: T | null = cur.left
      while (pred.right && pred.right !== cur) {
        pred = pred.right
      }
      if (!pred.right) {
        pred.right = cur
        cur = cur.left
      } else {
        pred.right = null
        visit(cur)
        cur = cur.right
      }
    } else {
      visit(cur)
      cur = cur.right
    }
  }
}

export { morris }

if (require.main === module) {
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

  // https://leetcode.cn/problems/recover-binary-search-tree/description/
  // eslint-disable-next-line no-inner-declarations
  function recoverTree(root: TreeNode | null): void {
    if (!root) return

    let first: TreeNode | null = null
    let second: TreeNode | null = null
    let prev: TreeNode = new TreeNode(-Infinity)

    morris(root, node => {
      if (prev.val > node.val) {
        if (!first) first = prev
        second = node
      }
      prev = node
    })

    const tmp = first!.val
    first!.val = second!.val
    second!.val = tmp
  }
}
