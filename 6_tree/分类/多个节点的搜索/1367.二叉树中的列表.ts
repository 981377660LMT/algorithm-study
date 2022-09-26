import { BinaryTree } from '../Tree'

class Node {
  value: number
  next: Node | null
  constructor(value: number, next: Node | null = null) {
    this.value = value
    this.next = next
  }
}

/**
 * @param {ListNode} head
 * @param {TreeNode} root
 * @return {boolean}
 * @description
 * 如果在二叉树中，存在一条一直向下的路径，
 * 且每个点的数值恰好一一对应以 head 为首的链表中每个节点的值，那么请你返回 True
 * @summary 注意这是以每个节点而不是根开始搜索
 */
function isSubPath(head: Node, root: BinaryTree | null): boolean {
  if (!root) {
    return false
  }

  return dfs(head, root) || isSubPath(head, root.left) || isSubPath(head, root.right)

  // 以root为节点时是否存在
  function dfs(head: Node | null, root: BinaryTree | null): boolean {
    if (!head) {
      return true
    }

    if (!root) {
      return false
    }

    if (head.value !== root.val) {
      return false
    }

    return dfs(head.next, root.left) || dfs(head.next, root.right)
  }
}

export {}
