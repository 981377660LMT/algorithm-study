class Node {
  val: number
  children: Node[]
  constructor(val = 0) {
    this.val = val
    this.children = []
  }
}

// function preorder(root: Node | null): number[] {
//   if (!root) return []
//   const res: number[] = []
//   const helper = (root: Node) => {
//     res.push(root.val)
//     for (let child of root.children) helper(child)
//   }
//   helper(root)
//   return res
// }

// 与二叉树类似的思路:
// 栈中存放未遍历的结点，遍历完当前结点后，将孩子结点逆序入栈（保证出栈顺序是顺序的）
// 先根后孩子
function preorder(root: Node | null): number[] {
  if (!root) return []

  const res: number[] = []
  const stack: Node[] = [root]

  while (stack.length) {
    const cur = stack.pop()!
    res.push(cur.val)
    for (let i = cur.children.length - 1; ~i; i--) {
      stack.push(cur.children[i])
    }
  }

  return res
}

// 先孩子后根
function inorder(root: Node | null): number[] {
  if (!root) return []

  const res: number[] = []
  const stack: Node[] = [root]

  while (stack.length) {
    const cur = stack.pop()!
    res.push(cur.val)
    stack.push(...cur.children)
  }

  return res.reverse()
}
export {}
