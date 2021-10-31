import { useUnionFindArray } from './推荐使用并查集精简版'

const findCircleNum = (isConnected: number[][]): number => {
  const uf = useUnionFindArray(isConnected.length)
  for (let i = 0; i < isConnected.length; i++) {
    for (let j = i; j < isConnected.length; j++) {
      if (isConnected[i][j] === 1) {
        uf.union(i, j)
      }
    }
  }

  return uf.getCount()
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
