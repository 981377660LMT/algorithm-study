class Node {
  val: number
  children: Node[]

  constructor(val = 0, children: Node[] = []) {
    this.val = val
    this.children = children
  }
}

// 给定一棵 N 叉树的根节点 root ，计算这棵树的直径长度。
// 思路一:建图 两次bfs，适用于所有的树

// 思路二：
// 套用模板，后序遍历，先获得各子树返回来的信息，这些信息包含着你要求的问题答案的子问题
// 类似于二叉树，此处只需保存子节点里两条最长的(二叉树本身就是left和right)

function diameter(root: Node): number {
  let res = 0
  dfs(root)
  return res

  function dfs(root: Node): number {
    if (!root) return 0

    // 维护两个最大
    let [max1, max2] = [0, 0]
    for (const child of root.children) {
      const childMax = dfs(child)
      if (childMax > max1) {
        max2 = max1
        max1 = childMax
      } else if (childMax > max2) max2 = childMax
    }

    res = Math.max(res, max1 + max2)
    return Math.max(max1, max2) + 1
  }
}

export {}
