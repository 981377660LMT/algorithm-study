/**
 * 根据parent遍历树，不需要parent有序。
 * @param parents parent为-1是根节点。
 * @param f 保证每个节点只调用一次，parent的调用一定早于children。
 * @see https://github.com/tdzl2003/algorithm_workspace
 */
function treeDpWithParents(parents: ArrayLike<number>, f: (v: number) => void) {
  const n = parents.length
  const visited = new Uint8Array(n)
  const stack = []
  for (let i = 0; i < n; i++) {
    for (let j = i; ~j && !visited[j]; j = parents[j]) {
      visited[j] = 1
      stack.push(j)
    }
    while (stack.length) {
      const v = stack.pop()!
      f(v)
    }
  }
}

export {}

if (require.main === module) {
  const parents = [-1, 0, 0, 1, 1, 2, 2, 3, 3, 4]
  const f = (v: number) => console.log(v)
  treeDpWithParents(parents, f)
}
