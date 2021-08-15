interface BinaryTree {
  val: number
  left: BinaryTree | null
  right: BinaryTree | null
}

const node = {
  val: 6,
  left: null,
  right: null,
}
const bt: BinaryTree = {
  val: 1,
  left: {
    val: 2,
    left: {
      val: 4,
      left: null,
      right: null,
    },
    right: {
      val: 5,
      left: null,
      right: null,
    },
  },
  right: {
    val: 3,
    left: node,
    right: {
      val: 6,
      left: null,
      right: null,
    },
  },
}
const btcopy: BinaryTree = {
  val: 1,
  left: {
    val: 2,
    left: {
      val: 4,
      left: null,
      right: null,
    },
    right: {
      val: 5,
      left: null,
      right: null,
    },
  },
  right: {
    val: 3,
    left: {
      val: 6,
      left: null,
      right: null,
    },
    right: {
      val: 6,
      left: null,
      right: null,
    },
  },
}

// 你 不能 对两棵二叉树，以及 target 节点进行更改。
// 只能 返回对克隆树 cloned 中已有的节点的引用。
// 树中允许出现值相同的节点 (需要使用===判断original的节点是否与target相等)
const getTargetCopy = function (
  original: BinaryTree | null,
  cloned: BinaryTree | null,
  target: BinaryTree
): BinaryTree | null {
  if (!cloned || !original) return null
  if (original === target) return cloned
  // original 与 cloned 同调
  const l = getTargetCopy(original.left, cloned.left, target)
  const r = getTargetCopy(original.right, cloned.right, target)
  return l || r
}

console.log(getTargetCopy(bt, btcopy, node))
export {}
