interface TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

const bt: TreeNode = {
  val: 3,
  left: {
    val: 4,
    left: {
      val: 1,
      left: null,
      right: null,
    },
    right: {
      val: 3,
      left: null,
      right: null,
    },
  },
  right: {
    val: 5,
    left: null,
    right: {
      val: 1,
      left: null,
      right: null,
    },
  },
}

// 如果抢了的话， 那么我们不能继续抢其左右子节点
// 如果不抢的话，那么我们可以继续抢左右子节点
const rob = (root: TreeNode): number => {
  const dfs = (root: TreeNode | null): [number, number] => {
    // res[0]表示不包括根节点的最大值，res[1]为包含根节点的最大值
    const res = [0, 0] as [number, number]
    if (!root) return res
    const left = dfs(root.left)
    const right = dfs(root.right)
    // 不包含根节点的最大值为左子树最大值加右子树最大值
    res[0] = Math.max.apply(null, left) + Math.max.apply(null, right)
    // 包含根节点的最大值为当前节点值加左子树不包含根节点的值加右子树不包含根节点的值
    res[1] = root.val + left[0] + right[0]
    return res
  }

  const res = dfs(root)
  return Math.max(res[0], res[1])
}

console.dir(rob(bt), { depth: null })
// 小偷一晚能够盗取的最高金额 = 4 + 5 = 9.
export {}
