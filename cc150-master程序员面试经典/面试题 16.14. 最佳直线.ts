/**
 * @param {number[][]} points
 * @return {number[]}
 * 请找出一条直线，其通过的点的数目最多。
 * @link https://leetcode-cn.com/problems/best-line-lcci/solution/shu-xue-ti-zhi-xian-de-yi-ban-shi-fang-cheng-by-tu/
 * 直线的一般式方程(A,B,C)作为key 记得求最大公约数
 */
var bestLine = function (points: number[][]): number[] {
  const gcd = (...nums: number[]) => {
    const twoNumGcd = (a: number, b: number): number => {
      return b === 0 ? a : twoNumGcd(b, a % b)
    }
    return nums.reduce(twoNumGcd)
  }

  const len = points.length

  let max = 0
  let res: number[] = []

  // 若有多条直线穿过了相同数量的点，则选择S[0]值较小的直线返回，S[0]相同则选择S[1]值较小的直线返回。
  // 考虑通过每个点的所有直线（和这个点之后的所有点形成的线段斜率），使用最大公约数归一化斜率为分数：
  for (let i = 0; i < len - 1; i++) {
    const counter = new Map<string, number>()
    const initialTwo = new Map<string, number[]>()
    for (let j = i + 1; j < len; j++) {
      const x1 = points[i][0]
      const y1 = points[i][1]
      const x2 = points[j][0]
      const y2 = points[j][1]
      let A = y2 - y1
      let B = x1 - x2
      let C = x2 * y1 - x1 * y2
      const divide = gcd(A, B, C)
      A /= divide
      B /= divide
      C /= divide

      // 斜率作为key
      // let A = y2 - y1
      // let B = x2 - x1
      // const divide = gcd(A, B)
      // A /= divide
      // B /= divide
      // const key = `${A}#${B}`

      const key = `${A}#${B}#${C}`
      !initialTwo.has(key) && initialTwo.set(key, [i, j])
      const curCount = (counter.get(key) || 0) + 1
      counter.set(key, curCount)
      if (curCount > max) {
        max = curCount
        res = initialTwo.get(key)!
      }
    }
  }

  return res
}

console.log(
  bestLine([
    [26072, -12996],
    [-41195, -34139],
    [6491, 14145],
    [275, 4007],
    [14321, -15055],
    [-38983, -49757],
    [-28710, -15391],
    [-29300, 12859],
    [-34606, -25274],
    [-37657, 14795],
    [-32300, 1599],
    [-24623, -14921],
    [-35555, -43348],
    [-41350, -16826],
    [-38848, -6454],
    [5588, -6650],
    [-8414, -38364],
    [15409, 20374],
    [29264, 28598],
    [-9066, -388],
    [-32891, -25982],
    [4402, 6766],
    [-12017, -22656],
    [-35555, -12886],
    [-10073, -11487],
    [10118, -18119],
    [-11704, 11154],
    [-25250, 28060],
    [-36845, -27142],
    [-16539, -8717],
    [-9274, 23635],
    [-7038, -17688],
    [-4654, -3477],
    [-30050, 10044],
    [-31933, -42528],
    [-20460, -15066],
    [27274, -18550],
    [22048, 28678],
    [-35555, -17101],
    [-33957, 26896],
    [-8262, -11077],
    [15830, -9823],
    [-38355, -30257],
    [14949, 4445],
    [-24900, 21759],
    [-24800, 29749],
    [-35555, -18966],
    [10459, 17639],
    [-40180, -29454],
    [-24194, -17257],
    [-9540, -28060],
    [-8029, -4150],
  ])
)
