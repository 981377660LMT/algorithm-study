import { deserializeNode } from './297二叉树的序列化与反序列化'

class TreeNode {
  val: number
  left?: TreeNode
  right?: TreeNode
  constructor(val: number = 0, left?: TreeNode, right?: TreeNode) {
    this.val = val
    this.left = left
    this.right = right
  }
}

// 对二叉搜索树采用中序遍历就能得到一个升序序列。
// 那么如果我们在遍历过程中，修改每一个根节点的左右指向，不就实现了原址修改了吗。
// 中序遍历，同时记录前驱节点
// cur.left = None
// pre.right = cur
// pre = cur
// 展开后的单链表应该与二叉树 中序遍历 顺序相同。

/**
 * @summary 二叉树转链表 关键是要记录pre节点
 */
function convertBiNode(root: TreeNode | undefined): TreeNode | undefined {
  if (!root) return root
  let res: TreeNode | undefined = undefined
  let pre: TreeNode | undefined = undefined

  const dfs = (root: TreeNode | undefined) => {
    if (!root) return
    dfs(root.left)

    // 当第一次执行到下面这一行代码，恰好是在最左下角:此时res是最左叶子节点
    !res && (res = root)
    root.left = undefined
    pre && (pre.right = root)
    pre = root

    dfs(root.right)
  }
  dfs(root)

  return res
}

console.dir(convertBiNode(deserializeNode([4, 2, 5, 1, 3, null, 6, 0]) as TreeNode))
