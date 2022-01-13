import { useUnionFindArray } from './推荐使用并查集精简版'

/**
 * @param {number[][]} stones
 * @return {number}
 * 可以将石子全部建立并查集的联系，并计算联通子图的个数。答案就是 n - 联通子图的个数
 */
const removeStones = function (stones: number[][]): number {
  const uf = useUnionFindArray(stones.length)
  for (let i = 0; i < stones.length; i++) {
    for (let j = i + 1; j < stones.length; j++) {
      if (stones[i][0] === stones[j][0] || stones[i][1] === stones[j][1]) {
        uf.union(i, j)
      }
    }
  }

  return stones.length - uf.getCount()
}

console.log(
  removeStones([
    [0, 0],
    [0, 1],
    [1, 0],
    [1, 2],
    [2, 1],
    [2, 2],
  ])
)
console.log(
  removeStones([
    [0, 0],
    [0, 2],
    [1, 1],
    [2, 0],
    [2, 2],
  ])
)

export {}
