import { SparseTableSqrt } from './SparseTableSqrt'

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

function treeQueries(root: TreeNode | null, queries: number[]): number[] {
  const n = 1e5
  const ins = new Uint32Array(n + 10) // 子树中最小的结点序号,const left = ins[removeRoot]
  const outs = new Uint32Array(n + 10) // 子树中最大的结点序号,即自己的id,const right = outs[removeRoot]
  const depth = new Uint32Array(n + 10) // 深度 根节点深度为0
  let dfsId = 1
  dfsOrder(root, 0)
  const st = new SparseTableSqrt(depth, () => 0, Math.max)

  const res = []
  for (let i = 0; i < queries.length; i++) {
    const removeRoot = queries[i]
    const left = ins[removeRoot]
    const right = outs[removeRoot]
    const max1 = st.query(0, left)
    const max2 = st.query(right + 1, n + 1)
    res.push(Math.max(max1, max2))
  }
  return res

  function dfsOrder(curRoot: TreeNode | null, dep: number): void {
    if (!curRoot) return
    ins[curRoot.val] = dfsId
    dfsOrder(curRoot.left, dep + 1)
    dfsOrder(curRoot.right, dep + 1)
    outs[curRoot.val] = dfsId
    depth[dfsId] = dep
    dfsId++
  }
}
