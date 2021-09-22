import { BinaryTree } from './力扣加加/Tree'

// 你需要以列表的形式返回上述重复子树的根结点。
function findDuplicateSubtrees(root: BinaryTree | null): Array<BinaryTree | null> {
  // 获取每个节点的唯一识别
  const counter = new Map<string, BinaryTree[]>()
  const dfs = (root: BinaryTree | null): string => {
    if (!root) return ''
    const l = dfs(root.left)
    const r = dfs(root.right)
    const key = `${root.val}#${l}#${r}`
    !counter.has(key) && counter.set(key, [])
    counter.get(key)!.push(root)
    return key
  }
  dfs(root)

  const res: BinaryTree[] = []
  for (const nodes of counter.values()) {
    if (nodes.length > 1) res.push(nodes[0])
  }
  return res
}
