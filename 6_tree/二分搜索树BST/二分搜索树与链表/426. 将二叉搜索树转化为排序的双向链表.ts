// 将一个 二叉搜索树 就地转化为一个 已排序的双向循环链表 。

import { BinaryTree } from '../../分类/Tree'

// 树中节点的左指针需要指向前驱，树中节点的右指针需要指向后继。还需要返回链表中最小元素的指针。
function treeToDoublyList(root: BinaryTree | null): BinaryTree | null {
  if (!root) return root
  let head: BinaryTree | null = null
  let pre: BinaryTree | null = null

  const dfs = (root: BinaryTree | null) => {
    if (!root) return
    dfs(root.left)

    // 当第一次执行到下面这一行代码，恰好是在最左下角:此时res是最左叶子节点
    if (pre) {
      pre.right = root
      root.left = pre
    } else {
      head = root
    }

    pre = root

    dfs(root.right)
  }
  dfs(root)

  // 最后pre为最右侧节点
  head!.left = pre
  pre!.right = head

  return head
}

export {}
