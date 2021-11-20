import { useUnionFindArray } from '../../14_并查集/推荐使用并查集精简版'

function countComponents(n: number, edges: number[][]): number {
  const uf = useUnionFindArray(n)

  for (const [u, v] of edges) {
    uf.union(u, v)
  }

  return uf.getCount()
}

console.log(
  countComponents(5, [
    [0, 1],
    [1, 2],
    [2, 3],
    [3, 4],
  ])
)
