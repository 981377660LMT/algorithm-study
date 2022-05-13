class Node {
  value: number
  next: Node | null
  constructor(value: number, next: Node | null = null) {
    this.value = value
    this.next = next
  }
}

const a = new Node(1)
const b = new Node(3)
const c = new Node(7)
a.next = b
b.next = c

interface BinaryTree {
  val: number
  left: BinaryTree | null
  right: BinaryTree | null
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
    left: {
      val: 6,
      left: null,
      right: null,
    },
    right: {
      val: 7,
      left: null,
      right: null,
    },
  },
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
const isSubPath = function (head: Node, root: BinaryTree | null): boolean {
  if (!root) return false

  // 以root为节点是否存在
  const dfs = (head: Node | null, root: BinaryTree | null): boolean => {
    if (!head) return true
    if (!root) return false
    if (head.value !== root.val) return false
    return dfs(head.next, root.left) || dfs(head.next, root.right)
  }

  return dfs(head, root) || isSubPath(head, root.left) || isSubPath(head, root.right)
}

export {}

console.log(isSubPath(a, bt))
