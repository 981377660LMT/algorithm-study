export {}

function maxTargetNodes(edges1: number[][], edges2: number[][]): number[] {
  const n = edges1.length + 1
  const m = edges2.length + 1
  const tree1: number[][] = Array.from({ length: n }, () => [])
  const tree2: number[][] = Array.from({ length: m }, () => [])
  edges1.forEach(([u, v]) => {
    tree1[u].push(v)
    tree1[v].push(u)
  })
  edges2.forEach(([u, v]) => {
    tree2[u].push(v)
    tree2[v].push(u)
  })

  const getDepths = (tree: number[][], size: number) => {
    const depths = new Int32Array(size).fill(-1)
    depths[0] = 0
    let even = 0
    let odd = 0
    const stack = new Uint32Array(size)
    let stackPtr = 0
    stack[stackPtr++] = 0
    while (stackPtr > 0) {
      const u = stack[--stackPtr]
      if (depths[u] % 2 === 0) even++
      else odd++
      for (const v of tree[u]) {
        if (depths[v] === -1) {
          depths[v] = depths[u] + 1
          stack[stackPtr++] = v
        }
      }
    }
    return { depths, even, odd }
  }

  const { depths: d1, even: c1Even, odd: c1Odd } = getDepths(tree1, n)
  const { even: c2Even, odd: c2Odd } = getDepths(tree2, m)
  const c2Max = Math.max(c2Even, c2Odd)
  const res: number[] = []
  for (let i = 0; i < n; i++) {
    const target = d1[i] % 2 === 0 ? c1Even : c1Odd
    res.push(target + c2Max)
  }
  return res
}
