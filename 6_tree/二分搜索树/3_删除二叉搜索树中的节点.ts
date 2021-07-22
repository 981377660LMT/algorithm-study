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
// 删除二叉搜索树中的 key 对应的节点，
// 并保证二叉搜索树的性质不变。返回二叉搜索树（有可能被更新）的根节点的引用。
// 注意这个函数最后返回的是对根节点的引用
// 策略:
// 0.先递归法找到要删除的节点
// 1.没有左右孩子直接返回null
// 2.左右孩子有一个则顶上
// 3.左右孩子都有则用前驱节点(左子树中最大节点)顶上，然后递归继续删除前驱节点
const deleteNode = (root: TreeNode, key: number): TreeNode | null => {
  if (!root) return null
  // 先找到点
  if (root.val < key) {
    root.right = deleteNode(root.right!, key)
  } else if (root.val > key) {
    root.left = deleteNode(root.left!, key)
  } else {
    // 这就是我要找的点
    // 如果目标点没有孩子
    if (!root.left && !root.right) {
      return null
    } else if (!root.left && root.right) {
      // 如果目标点只有右孩子 直接用右孩子代替自己
      return root.right
    } else if (!root.right && root.left) {
      // 如果目标点有左孩子 没有右孩子 直接用左孩子代替自己
      return root.left
    } else {
      // 如果目标点左右孩子都有 用前驱节点代替自己
      // 指针
      let pre = root.left
      while (pre?.right) {
        pre = pre.right
      }
      // 拿值过来
      root.val = pre?.val!
      // 把原本那个删了
      root.left = deleteNode(root.left!, pre?.val!)
    }
  }

  return root
}

console.dir(deleteNode(bt, 8), { depth: null })

export {}
