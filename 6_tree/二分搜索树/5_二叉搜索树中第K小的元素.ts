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
  function* inorder(root: TreeNode | null): Generator<number> {
    //no need to keep going after reach k-th number
    if (!root) return
    root.left && inorder(root.left)

    if (k === 1) yield root.val
    console.log(k, 666, root)
    k--
    root.right && inorder(root.right)
  }
  console.log(k)
  return inorder(root).next().value
}

console.dir(kthSmallest(bt, 4), { depth: null })

export {}
