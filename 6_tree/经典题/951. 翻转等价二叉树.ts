import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

// 一个翻转操作:选择任意节点，然后交换它的左子树和右子树。
// 编写一个判断两个二叉树是否是翻转等价的函数

// 注意先判断空 写python写多了就忘记了...
function flipEquiv(root1: BinaryTree | null, root2: BinaryTree | null): boolean {
  if (!root1 && !root2) return true
  if (!root1 || !root2) return false
  // LL,RR/LR，RL
  if (root1.val === root2.val) {
    // 同侧/异侧
    return (
      (flipEquiv(root1.left, root2.left) && flipEquiv(root1.right, root2.right)) ||
      (flipEquiv(root1.left, root2.right) && flipEquiv(root1.right, root2.left))
    )
  }
  return false
}

console.log(
  flipEquiv(
    deserializeNode([1, 2, 3, 4, 5, 6, null, null, null, 7, 8]),
    deserializeNode([1, 3, 2, null, 6, 4, 5, null, null, null, null, 8, 7])
  )
)
