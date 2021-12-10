class Node {
  val: number
  children: Node[]

  constructor(val = 0, children: Node[] = []) {
    this.val = val
    this.children = children
  }
}

// 给定一棵没有重复值的 N 叉树的根节点 root ，以及其中的两个节点 p 和 q。
// 移动节点 p 及其子树，使节点 p 成为节点 q 的直接子节点
// 如果 p 已经是 q 的直接子节点，则请勿改动任何节点。节点 p 必须是节点 q 的子节点列表的最后一项。

// 总结:要么q在p的子树里，要么不在，只要处理这两种情况就行。
// dfs过一遍`找到p和q的parent，顺便判断q是不是在p的子树里`。
// 有了p和q的parent节点就能随意移动p和q了(知道parent就可以定位、操作pq了)。
function moveSubTree(root: Node | null, p: Node | null, q: Node | null): Node | null {
  if (!root || !p || !q) return null

  const dummy = new Node(0, [root])
  let [parentOfP, parentOfQ] = [dummy, dummy]
  let qIsUnderP = false

  dfs(root, dummy, false)

  if (parentOfP === q) return root

  if (qIsUnderP) {
    // 把q接到p父亲下，再把p接到q下
    parentOfQ.children = parentOfQ.children.filter(node => node !== q)
    const pIndex = parentOfP.children.indexOf(p)
    parentOfP.children[pIndex] = q
    q.children.push(p)
  } else {
    parentOfP.children = parentOfP.children.filter(node => node !== p)
    q.children.push(p)
  }

  return dummy.children[0]

  /**
   * @description 找到p和q的parent，顺便判断q是不是在p的子树里
   */
  function dfs(root: Node, parent: Node, isUnderP: boolean) {
    if (root === p) {
      parentOfP = parent
      isUnderP = true
    }

    if (root === q) {
      if (isUnderP) qIsUnderP = true
      parentOfQ = parent
    }

    for (const child of root.children) {
      dfs(child, root, isUnderP)
    }
  }
}

export {}
