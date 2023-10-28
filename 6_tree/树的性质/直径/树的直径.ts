/**
 * 求带权树的(直径长度, 直径路径).
 */
function getDiameter(
  n: number,
  tree: [next: number, weight: number][][] | ArrayLike<ArrayLike<ArrayLike<number>>>,
  start = 0
): [diameter: number, path: number[]] {
  const dfs = (s: number): [endPoint: number, dist: Int32Array] => {
    const dist = new Int32Array(n).fill(-1)
    dist[s] = 0
    const stack = new Int32Array(n)
    let top = 0
    stack[top++] = s
    while (top) {
      const cur = stack[--top]
      const nexts = tree[cur]
      for (let i = 0; i < nexts.length; i++) {
        const { 0: next, 1: weight } = nexts[i]
        if (dist[next] !== -1) continue
        dist[next] = dist[cur] + weight
        stack[top++] = next
      }
    }

    let max = -1
    let endPoint = -1
    for (let i = 0; i < n; i++) {
      if (dist[i] > max) {
        max = dist[i]
        endPoint = i
      }
    }

    return [endPoint, dist]
  }

  const u = dfs(start)[0]
  let { 0: v, 1: dist } = dfs(u)
  const diameter = dist[v]
  const path = [v]
  while (u !== v) {
    const nexts = tree[v]
    for (let i = 0; i < nexts.length; i++) {
      const { 0: next, 1: weight } = nexts[i]
      if (dist[next] + weight === dist[v]) {
        path.push(next)
        v = next
        break
      }
    }
  }

  return [diameter, path]
}

export { getDiameter }

if (require.main === module) {
  const n = 5
  const tree = [
    [
      [1, 1],
      [2, 2]
    ],
    [
      [0, 1],
      [3, 3],
      [4, 4]
    ],
    [[0, 2]],
    [[1, 3]],
    [[1, 4]]
  ]
  console.log(getDiameter(n, tree))
}
