// 给你一个二叉树的根结点，返回其结点按 垂直方向（从上到下，逐列）遍历的结果。
// 如果两个结点在同一行和列，那么顺序则为 从左到右。
// 这个和987题：二叉树的垂序遍历仅有一字之差，有9.9分相似，一点差别在于，那个题在同一行和列，需要按值的大小升序，这个题是从左到右。
import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

function verticalOrder(root: BinaryTree | null): number[][] {
  if (!root) return []

  const nodes: [x: number, y: number, value: number][] = []
  const dfs = (root: BinaryTree, x: number, y: number) => {
    nodes.push([x, y, root.val])
    root.left && dfs(root.left, x - 1, y + 1)
    root.right && dfs(root.right, x + 1, y + 1)
  }
  dfs(root, 0, 0)
  nodes.sort((a, b) => a[0] - b[0] || a[1] - b[1])

  // x值为key的map
  const map = new Map<number, number[]>()
  for (const item of nodes) {
    const key = item[0]
    const val = item[2]
    !map.has(key) && map.set(key, [])
    map.get(key)!.push(val)
  }

  return [...map.values()]
}

export {}

console.log(
  verticalOrder(
    deserializeNode([1, 2, 3, 4, 5, 6, null, null, 7, 8, null, null, 9, null, 10, null, 11, 10])
  )
)
// [[4],[2,7,8],[1,5,6,10,11,10],[3,9]]
