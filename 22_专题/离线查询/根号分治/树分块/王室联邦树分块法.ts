// !- 王室联邦分块法
// !分块方式:满足每块大小在 [B,3B]，块内每个点到核心点路径上的所有点都在块内.
// !但是不保证每个块都是连通的.
//
// dfs，并创建一个栈，dfs一个点时先记录初始栈顶高度，
// 每dfs完当前节点的一棵子树就判断栈内（相对于刚开始dfs时）新增节点的数量是否≥B，
// 是则将栈内所有新增点分为同一块，核心点为当前dfs的点，
// 当前节点结束dfs时将当前节点入栈，
// 整个dfs结束后将栈内所有剩余节点归入已经分好的最后一个块。

/**
 * 王室联邦树分块法.
 * @param n 树的节点数.0~n-1.
 * @param tree 树的邻接表.
 * @param blockSize 每个块的大小.
 * @returns [belong, blockRoot]
 * - belong: 每个节点所属的块的编号.
 * - blockRoot: 每个块的根节点(关键点).
 */
function treeDecomposition(
  n: number,
  tree: ArrayLike<ArrayLike<number>>,
  blockSize = 1 + (Math.sqrt(n) | 0)
): [belong: Uint16Array, blockRoot: number[]] {
  const blockRoot: number[] = []
  const belong = new Uint16Array(n)
  const stack = new Uint32Array(n)
  let stackTop = 0

  const dfs = (cur: number, pre: number): void => {
    const size = stackTop
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; ++i) {
      const next = nexts[i]
      if (next === pre) continue
      dfs(next, cur)
      if (stackTop - size >= blockSize) {
        blockRoot.push(cur)
        while (stackTop > size) {
          belong[stack[--stackTop]] = blockRoot.length - 1
        }
      }
    }
    stack[stackTop++] = cur
  }

  dfs(0, -1)
  if (!blockRoot.length) blockRoot.push(0)
  while (stackTop) {
    belong[stack[--stackTop]] = blockRoot.length - 1
  }

  return [belong, blockRoot]
}

export { treeDecomposition }

if (require.main === module) {
  const n = 5
  const tree = [[1], [0, 2, 3], [1, 4], [1], [2]]
  const blockSize = 2
  console.log(treeDecomposition(n, tree, blockSize))
}
