import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

// 请你删除所有值为 target 的 叶子节点 。
// 注意，一旦删除值为 target 的叶子节点，它的父节点就可能变成叶子节点；如果新叶子节点的值恰好也是 target ，那么这个节点也应该被删除。
function removeLeafNodes(root: BinaryTree | null, target: number): BinaryTree | null {
  const dummy = new BinaryTree(Infinity, root)
  dfs(dummy)
  return dummy.left

  // function dfs(root: BinaryTree | null): void {
  //   if (!root) return
  //   dfs(root.left)
  //   dfs(root.right)
  //   错误解法:只是指针置为null,并不能删除对象
  //   if (!root.left && !root.right && root.val === target) root = null
  // }

  function dfs(root: BinaryTree | null): BinaryTree | null {
    if (!root) return null
    root.left = dfs(root.left)
    root.right = dfs(root.right)
    if (!root.left && !root.right && root.val === target) return null
    return root
  }
}

console.log(removeLeafNodes(deserializeNode([1, 2, 3, 2, null, 2, 4]), 2))

// root = null删不掉的原因:
// const obj1 = { a: 11 }
// let p: any = obj1
// p = null
// console.log(obj1)
// 只是指针置为null,并不能删除对象
