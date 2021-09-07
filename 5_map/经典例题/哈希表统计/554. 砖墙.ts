type Prefix = number
type Count = number

/**
 * @param {number[][]} wall
 * @return {number}
 * @description
 * 你需要找出怎样画才能使这条线 穿过的砖块数量最少 ，并且返回 穿过的砖块数量 。
 * 穿过的砖块数量最少===经过的交界线最多(除了两边)
 */
const leastBricks = function (wall: number[][]): number {
  const map = new Map<Prefix, Count>()

  for (const row of wall) {
    let sum = 0
    // 每一行砖的最后一列不要计算进来, 否则会把最右侧的垂直线考虑进去
    for (let i = 0; i < row.length - 1; i++) {
      sum += row[i]
      map.set(sum, (map.get(sum) || 0) + 1)
    }
  }

  return wall.length - Math.max(...map.values(), 0)
}

console.log(
  leastBricks([
    [1, 2, 2, 1],
    [3, 1, 2],
    [1, 3, 2],
    [2, 4],
    [3, 1, 2],
    [1, 3, 1, 1],
  ])
)

// 2

export {}
