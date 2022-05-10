import { deserializeNode } from '../../6_tree/力扣加加/构建类/297.二叉树的序列化与反序列化'

class TreeNode {
  val: number
  left: TreeNode | null
  right: TreeNode | null
  constructor(val?: number, left?: TreeNode | null, right?: TreeNode | null) {
    this.val = val === undefined ? 0 : val
    this.left = left === undefined ? null : left
    this.right = right === undefined ? null : right
  }
}

// 勇者`最多`可以击破一个节点泡泡
// 2 <= 树中节点个数 <= 10^5
// -10000 <= 树中节点的值 <= 10000
function getMaxLayerSum(root: TreeNode | null): number {
  // dfs序处理子树
  let id = 1
  const start = new Map<TreeNode, number>() // 子树最开始的结点序号
  const IdByNode = new Map<TreeNode, number>() // 本身最后映射到几(id)
  // start[node] = 1, IdByNode[node] = 6，表示编号为 node 的子树映射到的区间是 [1, 6]，本身映射到 6

  const depthById = new Map<number, number>() // 每个节点的深度
  const sumByDepth = new Map<number, number>() // 每个深度的和
  const valueById = new Map<number, number>()
  const adjMap = new Map<number, number[]>()
  const treePresum = new Map<number, number[]>() // 每个子树层值得前缀和
  const preSum = [0]

  const dfs1 = (cur: TreeNode | null): void => {
    if (!cur) return
    start.set(cur, id)
    dfs1(cur.left)
    dfs1(cur.right)
    IdByNode.set(cur, id)
    id++
  }
  dfs1(root)

  // console.log(IdByNode)

  const dfs2 = (cur: TreeNode | null, depth: number, parent: number): void => {
    if (!cur) return
    const curId = IdByNode.get(cur)!
    valueById.set(curId, cur.val)

    depthById.set(curId, depth)
    sumByDepth.set(depth, (sumByDepth.get(depth) ?? 0) + cur.val)

    !treePresum.has(curId) && treePresum.set(curId, [0])

    if (parent !== -1) {
      !adjMap.has(parent) && adjMap.set(parent, [])
      adjMap.get(parent)!.push(curId)
    }

    dfs2(cur.left, depth + 1, curId)
    dfs2(cur.right, depth + 1, curId)
  }
  dfs2(root, 0, -1)
  //   depthById Map(3) { 0 => [ 4 ], 1 => [ 2, 3 ], 2 => [ 1 ] }
  // valueById Map(4) { 4 => 6, 2 => 0, 1 => 8, 3 => 3 }
  // adjMap Map(2) { 4 => [ 2, 3 ], 2 => [ 1 ] }
  // console.log('depthById', depthById)
  // console.log('valueById', valueById)
  // console.log('adjMap', adjMap)
  // for (let i = 1; i <= id; i++) {
  //   preSum.push(preSum[i - 1] + (valueById.get(i) ?? 0))
  // }

  // 不破泡
  let res = Math.max(...sumByDepth.values())
  const cands = [...adjMap.keys()].filter(cur => adjMap.get(cur)!.length === 1)
  for (const cur of cands) {
    update(cur)
  }

  return res

  // 删除root结点，看当前
  function update(root: number): void {
    const curDepth = depthById.get(root)!
    const curSum = sumByDepth.get(curDepth)!

    let childSum = 0
    for (const next of adjMap.get(root)!) {
      childSum += valueById.get(next)!
    }

    const diff = childSum - valueById.get(root)!
    res = Math.max(res, curSum + diff)
  }
}

const root1 = deserializeNode([6, 0, 3, null, 8])
const root2 = deserializeNode([5, 6, 2, 4, null, null, 1, 3, 5])
const root3 = deserializeNode([-5, 1, 7])
console.log(getMaxLayerSum(root1))
// console.log(getMaxLayerSum(root2))
// console.log(getMaxLayerSum(root3))

export {}
