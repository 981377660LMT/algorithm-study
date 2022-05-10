import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297.二叉树的序列化与反序列化'

// 「排序二叉搜索树节点，每一个节点都必须排在它的子孙结点前面:想到了层序遍历」
// 与层序遍历的区别是:层序遍历需要每次处理处理完一层，而这里可以不要，只要保证对于每个节点子孙节点在其后面即可
// 用队列就可以保证这种关系
function BSTSequences(root: BinaryTree | null): number[][] {
  if (!root) return [[]]

  const res: number[][] = []
  const bt = (root: BinaryTree, path: number[], queue: BinaryTree[]) => {
    root.left && queue.push(root.left)
    root.right && queue.push(root.right)
    if (!queue.length) return res.push(path.slice())
    for (let i = 0; i < queue.length; i++) {
      const nextRoot = queue[i]
      path.push(nextRoot.val)
      bt(nextRoot, path, [...queue.slice(0, i), ...queue.slice(i + 1)])
      path.pop()
    }
  }
  bt(root, [root.val], [])

  return res
}

console.log(BSTSequences(deserializeNode([2, 1, 3])))
