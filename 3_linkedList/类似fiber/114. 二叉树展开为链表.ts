// 将二叉树按前序遍历顺序展开为单链表
// 展开后的单链表应该同样使用 TreeNode ，其中 right 子指针指向链表中下一个结点，而左子指针始终为 null
// 要按照 6→5→4→3→2→1 的顺序访问节点，也就是按照右子树 - 左子树 - 根的顺序 DFS 这棵树
export {}

class TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null

  constructor(val?: number, left?: TreeNode | null, right?: TreeNode | null) {
    this.val = val === undefined ? 0 : val
    this.left = left === undefined ? null : left
    this.right = right === undefined ? null : right
  }
}

function flatten(root: TreeNode | null): void {
  let head: TreeNode | null = null

  const dfs = (curRoot: TreeNode | null): void => {
    if (!curRoot) return
    dfs(curRoot.right)
    dfs(curRoot.left)
    curRoot.left = null
    curRoot.right = head
    head = curRoot
  }

  dfs(root)
}
