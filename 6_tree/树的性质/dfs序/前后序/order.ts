//       0
//     / | \
//   1   2  3
//  / \     |
// 4   5    6
//
// !preOrder:
//  0 => [0, 7)
//  1 => [1, 4)
//  4 => [2, 3)
//  5 => [3, 4)
//  2 => [4, 5)
//  3 => [5, 7)
//  6 => [6, 7)
//
// !postOrder:
//  4 => [0, 1)
//  5 => [1, 2)
//  1 => [0, 3)
//  2 => [3, 4)
//  6 => [4, 5)
//  3 => [4, 6)
//  0 => [0, 7)

/**
 * 前序遍历dfs序.
 * !data[lid[i]] = values[i].
 */
function dfsPreOrder(tree: number[][], root = 0): [lid: Uint32Array, rid: Uint32Array] {
  const n = tree.length
  const lid = new Uint32Array(n)
  const rid = new Uint32Array(n)
  let dfn = 0

  const dfs = (cur: number, pre: number): void => {
    lid[cur] = dfn
    dfn++
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next !== pre) {
        dfs(next, cur)
      }
    }
    rid[cur] = dfn
  }
  dfs(root, -1)
  return [lid, rid]
}

/**
 * 后序遍历dfs序.
 * !data[rid[i]-1] = values[i].
 */
function dfsPostOrder(tree: number[][], root = 0): [lid: Uint32Array, rid: Uint32Array] {
  const n = tree.length
  const lid = new Uint32Array(n)
  const rid = new Uint32Array(n)
  let dfn = 0

  const dfs = (cur: number, pre: number): void => {
    lid[cur] = dfn
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i]
      if (next !== pre) {
        dfs(next, cur)
      }
    }
    dfn++
    rid[cur] = dfn
  }
  dfs(root, -1)
  return [lid, rid]
}

export { dfsPreOrder, dfsPostOrder }
