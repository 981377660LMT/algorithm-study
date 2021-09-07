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
// 1.可以直接中序遍历，并一边遍历一边判断遍历结果是否是单调递增的，如果不是则提前返回 False 即可。
// 2. 容易想到的做法
// const isValidBST = (root: TreeNode) => {
//   if (!root) return true
//   let isValidBST = true

//   const val = (root: TreeNode) => {
//     if (!root) return
//     const isValidNode = (root: TreeNode) => {
//       const validRight = !root.right || (root.right && root.right.val > root.val)
//       const validLeft = !root.left || (root.left && root.left.val < root.val)
//       return validRight && validLeft
//     }
//     if (!isValidNode(root)) return (isValidBST = false)
//     root.left && val(root.left)
//     root.right && val(root.right)
//   }
//   val(root)

//   return isValidBST
// }

// 3.比较难理解的递归中序遍历
const isValidBST = (root: TreeNode) => {
  if (!root) return true
  let pre: TreeNode | null = null
  const dfs = (root: TreeNode | null): boolean => {
    if (!root) return true
    if (!dfs(root.left)) return false
    if (pre && pre.val >= root.val) return false
    // pre最开始是在最左下角
    pre = root
    if (!dfs(root.right)) return false
    return true
  }
  return dfs(root)
}
console.dir(isValidBST(bt), { depth: null })

export {}
