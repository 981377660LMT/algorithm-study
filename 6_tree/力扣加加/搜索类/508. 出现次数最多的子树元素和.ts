import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 *
 * @param root
 * 你需要返回出现次数最多的子树元素和
 * 如果有多个元素出现的次数相同，返回所有出现次数最多的子树元素和
 * 一个结点的「子树元素和」定义为以该结点为根的二叉树上所有结点的元素之和
 */
function findFrequentTreeSum(root: BinaryTree | null): number[] {
  if (!root) return []
  const counter = new Map<number, number>()
  const dfs = (root: BinaryTree | null): number => {
    if (!root) return 0
    const left = dfs(root.left)
    const right = dfs(root.right)
    const sum = root.val + right + left
    counter.set(sum, (counter.get(sum) || 0) + 1) // 只需统计本层
    return sum
  }
  dfs(root)
  const max = Math.max(...counter.values())
  return [...counter.keys()].filter(key => counter.get(key) === max)
}

console.log(findFrequentTreeSum(deserializeNode([5, 2, -3])))
