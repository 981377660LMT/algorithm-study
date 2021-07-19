import { UnionFind } from './0_并查集'

const findCircleNum = (isConnected: number[][]) => {
  const uf = new UnionFind<number>()
  for (let i = 0; i < isConnected.length; i++) {
    uf.add(i)
    for (let j = 0; j < isConnected.length; j++) {
      if (isConnected[i][j] === 1) {
        uf.union(i, j)
      }
    }
  }
  console.log(uf)
  return uf.count
}

console.log(
  findCircleNum([
    [1, 1, 0],
    [1, 1, 0],
    [0, 0, 1],
  ])
)
console.log(
  findCircleNum([
    [1, 0, 0, 1],
    [0, 1, 1, 0],
    [0, 1, 1, 1],
    [1, 0, 1, 1],
  ])
)

export {}
