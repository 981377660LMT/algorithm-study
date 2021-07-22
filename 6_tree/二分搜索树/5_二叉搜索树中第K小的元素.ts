interface TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
}

const bt: TreeNode = {
  val: 6,
  left: {
    val: 2,
    left: {
      val: 0,
      left: null,
      right: null,
    },
    right: {
      val: 4,
      left: {
        val: 3,
        left: null,
        right: null,
      },
      right: {
        val: 5,
        left: null,
        right: null,
      },
    },
  },
  right: {
    val: 8,
    left: {
      val: 7,
      left: null,
      right: null,
    },
    right: {
      val: 9,
      left: null,
      right: null,
    },
  },
}

// 思路:中序遍历
const kthSmallest = (root: TreeNode, k: number): number => {
  const res: number[] = []
  const inorder = (root: TreeNode | null) => {
    //no need to keep going after reach k-th number
    if (!root || res.length >= k) return
    root.left && inorder(root.left)
    res.push(root.val)
    root.right && inorder(root.right)
  }
  inorder(root)
  // root.left && console.log(root.val)
  return res[k - 1]
}

console.dir(kthSmallest(bt, 2), { depth: null })

export {}
