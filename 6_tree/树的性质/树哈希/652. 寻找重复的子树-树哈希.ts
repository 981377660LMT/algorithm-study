/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-shadow */

// 如果两棵树具有 相同的结构 和 相同的结点值 ，则认为二者是 重复 的。
// 对于同一类的重复子树，你只需要返回其中任意 一棵 的根结点即可。
// !注意题目要求子树对应位置也要一样
// !n<=5000
// !-200<=Node.val<=200

// 解法1：n元组字符串表示子树 长链(n为1e5)时会TLE
function findDuplicateSubtrees(root: TreeNode): Array<TreeNode | null> {
  // 获取每个节点的唯一识别
  const counter = new Map<unknown, TreeNode[]>()
  const res: TreeNode[] = []
  dfs(root)
  for (const nodes of counter.values()) {
    if (nodes.length > 1) res.push(nodes[0])
  }

  return res

  function dfs(root: TreeNode | null): string {
    if (!root) return ''
    const subTree = [String(root.val)] // !子树顺序对结果有影响,需要`#`分隔成n元组
    subTree.push(dfs(root.left))
    subTree.push(dfs(root.right))

    const key = subTree.join('#')
    !counter.has(key) && counter.set(key, [])
    counter.get(key)!.push(root)
    return key
  }
}

// 解法2：优化
// !使用`哈希值的编号`来代替很长的哈希字符串 减少哈希值长度
function findDuplicateSubtrees2(root: TreeNode | null): Array<TreeNode | null> {
  const visited = new Map<string, [node: TreeNode, hashId: number][]>()
  const res: TreeNode[] = []
  let hashId = 0 // !代表不同哈希值的编号
  dfs(root)
  for (const nodes of visited.values()) {
    if (nodes.length > 1) res.push(nodes[0][0])
  }

  return res

  function dfs(root: TreeNode | null): number {
    if (!root) return 0
    const subTree = [String(root.val)]
    subTree.push(String(dfs(root.left)))
    subTree.push(String(dfs(root.right)))

    const key = subTree.join('#')
    if (visited.has(key)) {
      const pairs = visited.get(key)!
      const curId = pairs[0][1]
      pairs.push([root, curId])
      return curId
    }

    hashId++
    visited.set(key, [])
    visited.get(key)!.push([root, hashId])
    return hashId
  }
}
