import { BinaryTree } from './分类/Tree'
import { deserializeNode } from './重构json/297.二叉树的序列化与反序列化'

/**
 * @description
 * 现在这个二叉树受到「污染」，所有的 treeNode.val 都变成了 -1。
 * @link https://leetcode-solution-leetcode-pp.gitbook.io/leetcode-solution/medium/1261.find-elements-in-a-contaminated-binary-tree
 */
class FindElements {
  private root: BinaryTree | null

  /**
   *
   * @param root 用受污染的二叉树初始化对象，你需要先把它还原。
   *  0
   * 1 2
   * ....
   */
  constructor(root: BinaryTree | null) {
    root && (root.val = 0)
    this.root = this.recover(root)
  }

  /**
   *
   * @param target 判断目标值 target 是否存在于还原后的二叉树中并返回结果。
   * @summary
   * 加一转二进制后
   * 所有子节点和父节点都有相同的前缀,当最后一位是0时则走左侧，是1时则走右侧
   * 直接用 target + 1 的二进制表示进行二叉树寻路 即可。
   */
  find(target: number): boolean {
    let root = this.root
    let directions = (target + 1).toString(2).split('')

    for (let i = 1; i < directions.length; i++) {
      if (directions[i] === '0') root = root!.left
      else root = root!.right
      if (!root) return false
    }

    return true
  }

  private recover(root: BinaryTree | null) {
    if (!root) return null
    root.left && (root.left.val = 2 * root.val + 1)
    root.right && (root.right.val = 2 * root.val + 2)
    this.recover(root.left)
    this.recover(root.right)
    return root
  }
}

const FE = new FindElements(deserializeNode([-1, null, -1]))
console.log(FE.find(1))
console.log(FE.find(2))
