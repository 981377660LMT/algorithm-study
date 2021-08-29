// morris 遍历 是可以在常数的空间复杂度完成树的遍历的一种算法。
interface BinaryTree {
  val: number
  left: BinaryTree | undefined
  right: BinaryTree | undefined
}

const bt: BinaryTree = {
  val: 4,
  left: {
    val: 2,
    left: {
      val: 1,
      left: undefined,
      right: undefined,
    },
    right: {
      val: 3,
      left: undefined,
      right: undefined,
    },
  },
  right: {
    val: 6,
    left: {
      val: 5,
      left: undefined,
      right: undefined,
    },
    right: {
      val: 7,
      left: undefined,
      right: undefined,
    },
  },
}

// morris中序遍历:对于没有左子树的节点只到达一次，对于有左子树的节点会到达两次
// https://leetcode-cn.com/problems/recover-binary-search-tree/solution/yi-wen-zhang-wo-morrisbian-li-suan-fa-by-a-fei-8/
const morris = (root: BinaryTree) => {
  const getPredecessor = (root: BinaryTree) => {
    let rootP = root
    if (rootP.left) {
      rootP = rootP.left
      // rootP.right !== root是什么意思
      // 因为每个pre节点都会指向root 所以这里不能指回root
      while (rootP.right && rootP.right !== root) {
        rootP = rootP.right
      }
    }
    return rootP
  }

  let rootP: BinaryTree | undefined = root
  while (rootP) {
    // 1.如果当前节点的左孩子为空，则输出当前节点并将其右孩子作为当前节点。 rootP = rootP.right
    if (!rootP.left) {
      console.log(`${rootP.val} `)
      rootP = rootP.right
    } else {
      // 2.如果当前节点的左孩子不为空，在当前节点的左子树中找到当前节点在中序遍历下的前驱节点
      let pre = getPredecessor(rootP)
      // 2.1 如果前驱节点（mostRight）的右孩子为空，将它的右孩子指向当前节点
      if (!pre.right) {
        pre.right = rootP
        rootP = rootP.left
        // 2.2 如果前驱节点（mostRight）的右孩子为当前节点
      } else if (pre.right === rootP) {
        pre.right = undefined
        console.log(`${rootP.val} `)
        rootP = rootP.right
      }
    }
  }
}
morris(bt)

export {}
