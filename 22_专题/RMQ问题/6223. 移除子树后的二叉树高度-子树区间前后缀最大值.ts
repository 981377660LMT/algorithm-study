import { SparseTable } from './SparseTable'
import { deserializeNode } from '../../6_tree/重构json/297.二叉树的序列化与反序列化'

// !看到子树 直接弄到区间上
function treeQueries(root: TreeNode | null, queries: number[]): number[] {
  const n = 1e5
  const ins = new Uint32Array(n + 10) // 子树中最小的结点序号,const left = ins[removeRoot]
  const outs = new Uint32Array(n + 10) // 子树中最大的结点序号,即自己的id,const right = outs[removeRoot]
  const depth = new Uint32Array(n + 10) // 深度 根节点深度为0
  let dfsId = 1
  dfsOrder(root, 0)

  // const st = new SparseTable(depth, Math.max)
  // !前后缀最大值
  const preMax = depth.slice()
  for (let i = 1; i < n; i++) preMax[i] = Math.max(preMax[i], preMax[i - 1])
  const sufMax = depth.slice()
  for (let i = n - 2; i >= 0; i--) sufMax[i] = Math.max(sufMax[i], sufMax[i + 1])

  const res: number[] = []
  for (let i = 0; i < queries.length; i++) {
    const removeRoot = queries[i]
    const left = ins[removeRoot]
    const right = outs[removeRoot]
    // const max1 = left - 1 >= 1 ? st.query(1, left) : 0
    // const max2 = right + 1 <= n ? st.query(right + 1, n + 1) : 0
    const max1 = preMax[left - 1]
    const max2 = sufMax[right + 1]
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

// root = [5,8,9,2,1,3,7,4,6], queries = [3,2,4,8]
if (require.main === module) {
  const node = deserializeNode([5, 8, 9, 2, 1, 3, 7, 4, 6])
  console.log(treeQueries(node, [3, 2, 4, 8]))
}

export {}
