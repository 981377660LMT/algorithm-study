class TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
  constructor(val: number = 0) {
    this.val = val
    this.left = null
    this.right = null
  }
}

class Node {
  value: number
  next: Node | undefined
  constructor(value: number, next?: Node) {
    this.value = value
    this.next = next
  }
}

const a = new Node(-10)
const b = new Node(-3)
const c = new Node(0)
const d = new Node(5)
const e = new Node(9)
a.next = b
b.next = c
c.next = d
d.next = e

// 关键:根节点是中位数,这样分配可以保证左右子树的节点数目差不超过 1
// 选择中点作为根节点，根节点左侧的作为左子树，右侧的作为右子树即可。原因很简单，这样分配可以保证左右子树的节点数目差不超过 1。因此高度差自然也不会超过 1 了。

// 链表选择中点需要快慢指针

const sortedListToBST = (head: Node): TreeNode | null => {
  if (!head) return null
  const toTree = (
    head: Node | null | undefined,
    tail: Node | null | undefined
  ): TreeNode | null => {
    if (!head || head === tail) return null
    let slow: Node | null | undefined = head
    let fast: Node | null | undefined = head
    while (fast && fast.next && fast !== tail && fast.next !== tail) {
      fast = fast.next.next
      slow = slow!.next
    }
    const root = new TreeNode(slow!.value)
    root.left = toTree(head, slow!)
    root.right = toTree(slow!.next, tail)
    return root
  }

  return toTree(head, null)
}

console.dir(sortedListToBST(a), { depth: null })
// 输出：[0,-3,9,-10,null,5]
export {}
