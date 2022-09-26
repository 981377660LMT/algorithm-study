/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-shadow */

// 你需要以列表的形式返回上述重复子树的根结点。
// 树哈希与树的同构
function findDuplicateSubtrees(root: TreeNode | null): Array<TreeNode | null> {
  // 获取每个节点的唯一识别
  const counter = new Map<string, TreeNode[]>()
  const res: TreeNode[] = []
  dfs(root)
  for (const nodes of counter.values()) {
    if (nodes.length > 1) res.push(nodes[0])
  }

  return res

  function dfs(root: TreeNode | null): string {
    if (!root) return ''

    const left = dfs(root.left)
    const right = dfs(root.right)
    const key = `${root.val}#${left}#${right}`
    !counter.has(key) && counter.set(key, [])
    counter.get(key)!.push(root)
    return key
  }
}

// 如果是多叉树呢？1948
// function genHash(root: TrieNode, counter: Map<string, number>): void {
//   // 这句话很关键
//   if (root.children.size === 0) return

//   const sb: string[] = []
//   for (const [childName, child] of root.children.entries()) {
//     genHash(child, counter)
//     sb.push(`${childName}(${child.subtreeHash})`)
//   }

//   sb.sort()
//   root.subtreeHash = sb.join('')
//   counter.set(root.subtreeHash, (counter.get(root.subtreeHash) || 0) + 1)
// }
