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

const isValidBST = (root: TreeNode) => {
  if (!root) return true
  let isValidBST = true

  const val = (root: TreeNode) => {
    if (!root) return
    const isValidNode = (root: TreeNode) => {
      const validRight = !root.right || (root.right && root.right.val > root.val)
      const validLeft = !root.left || (root.left && root.left.val < root.val)
      return validRight && validLeft
    }
    if (!isValidNode(root)) return (isValidBST = false)
    root.left && val(root.left)
    root.right && val(root.right)
  }
  val(root)

  return isValidBST
}

console.dir(isValidBST(bt), { depth: null })

export {}
