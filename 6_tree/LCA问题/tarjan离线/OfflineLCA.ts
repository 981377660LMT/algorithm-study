// https://ei1333.github.io/library/graph/tree/offline-lca.hpp
// O(n*α(n))，LCA离线，一般用于树上莫队

/**
 * 离线求LCA.
 */
function offLineLca(
  tree: number[][],
  queries: [start: number, end: number][] | number[][],
  root = 0
): number[] {
  const n = tree.length
  const data = new Int32Array(n)
  for (let i = 0; i < n; ++i) data[i] = -1
  const stack = new Uint32Array(n)
  const mark = new Int32Array(n)
  const ptr = new Uint32Array(n)
  const res = Array(queries.length)
  for (let i = 0; i < queries.length; ++i) res[i] = -1

  let top = 0
  stack[top] = root
  queries.forEach(([start, end]) => {
    mark[start]++
    mark[end]++
  })

  const q: [next: number, ei: number][][] = Array(n)
  for (let i = 0; i < n; ++i) {
    q[i] = []
    mark[i] = -1
    ptr[i] = tree[i].length
  }
  queries.forEach(([start, end], i) => {
    q[start].push([end, i])
    q[end].push([start, i])
  })

  while (top !== -1) {
    const u = stack[top]
    const cache = tree[u]
    if (mark[u] === -1) {
      mark[u] = u
    } else {
      union(u, cache[ptr[u]])
      mark[find(u)] = u
    }

    if (!run(u)) {
      q[u].forEach(([v, i]) => {
        if (mark[v] !== -1 && res[i] === -1) {
          res[i] = mark[find(v)]
        }
      })
      --top
    }
  }

  return res

  function union(key1: number, key2: number): boolean {
    let root1 = find(key1)
    let root2 = find(key2)
    if (root1 === root2) return false
    if (data[root1] > data[root2]) {
      root1 ^= root2
      root2 ^= root1
      root1 ^= root2
    }
    data[root1] += data[root2]
    data[root2] = root1
    return true
  }

  function find(key: number): number {
    if (data[key] < 0) return key
    data[key] = find(data[key])
    return data[key]
  }

  function run(u: number): boolean {
    const cache = tree[u]
    while (ptr[u]) {
      const v = cache[--ptr[u]]
      if (mark[v] === -1) {
        stack[++top] = v
        return true
      }
    }
    return false
  }
}

export { offLineLca }

if (require.main === module) {
  const n = 5
  const edges = [
    [0, 1],
    [0, 2],
    [0, 3],
    [2, 4]
  ]
  const tree: number[][] = Array(n)
  for (let i = 0; i < n; ++i) tree[i] = []
  edges.forEach(([u, v]) => {
    tree[u].push(v)
    tree[v].push(u)
  })

  const queries = [
    [1, 3],
    [2, 4],
    [0, 4]
  ]
  const res = offLineLca(tree, queries)
  console.log(res)
}
