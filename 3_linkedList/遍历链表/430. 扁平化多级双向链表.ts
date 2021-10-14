class Node {
  val: number
  prev: Node | undefined
  next: Node | undefined
  child: Node | undefined
  constructor(val: number, prev?: Node, next?: Node, child?: Node) {
    this.val = val
    this.prev = prev
    this.next = next
    this.child = child
  }
}
/**
 * @link https://leetcode-cn.com/problems/flatten-a-multilevel-doubly-linked-list/solution/bian-ping-hua-duo-ji-shuang-xiang-lian-biao-by-lee/
 * @param head 多级双向链表中，除了指向下一个节点和前一个节点指针之外，它还有一个子链表指针
 * 
 * @description
 *我们可能会疑问什么情况下会使用这样的数据结构。
  其中一个场景就是 git 分支的简化版本。
  通过扁平化多级列表，可以认为将所有 git 的分支合并在一起。

  我们将列表顺时针转 90 °，那么就会看到一颗二叉树，
  则扁平化的操作也就是对二叉树进行先序遍历（深度优先搜索）
  我们可以将 child 指针当作二叉树中指向左子树的 left 指针。
  同样，next 指针可以当作是二叉树中的 right 指针
  @summary 很像 6_tree\力扣加加\构建类\114. 二叉树展开为链表.ts
  使用pre来穿针引线 因为最后链表起点是根节点 所以先序遍历
 */
function flatten(head: Node | undefined): Node | undefined {
  const dummy = new Node(0, undefined, head, undefined)
  let pre = dummy

  const dfs = (root: Node | undefined) => {
    if (!root) return

    pre.next = root
    root.prev = pre
    const [left, right] = [root.child, root.next]
    root.child = root.next = undefined // 最后只保留next 和 prev
    pre = root

    dfs(left)
    dfs(right)
  }
  dfs(head)

  // 记得断开连接
  const res = dummy.next
  res && (res.prev = undefined)
  return res
}

export {}

const a = new Node(1)
const b = new Node(2)
const c = new Node(3)
const d = new Node(4)
const e = new Node(5)
a.next = b
b.prev = a
b.next = c
c.prev = b
c.child = d
d.next = e
e.prev = d

console.dir(flatten(a), { depth: null })
