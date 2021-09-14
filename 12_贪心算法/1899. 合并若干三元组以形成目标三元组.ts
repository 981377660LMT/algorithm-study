/**
 * @param {number[][]} triplets 1 <= triplets.length <= 105
 * @param {number[]} target
 * @return {boolean}
 * triplets[j] 为 [max(ai, aj), max(bi, bj), max(ci, cj)] 。
 * 如果通过以上操作我们可以使得目标 三元组 target 成为 triplets 的一个 元素 ，则返回 true
 * @summary
 * 1.操作与顺序无关 具有交换性
 * 2.每次操作是不减的
 * 3.triplets 中一定不能有比 target 对应位大
 */
const mergeTriplets = function (triplets: number[][], target: number[]): boolean {
  const [tx, ty, tz] = target
  let [curX, curY, curZ] = [0, 0, 0]

  for (const [a, b, c] of triplets) {
    if (a <= tx && b <= ty && c <= tz) {
      curX = Math.max(curX, a)
      curY = Math.max(curY, b)
      curZ = Math.max(curZ, c)
    }
  }

  return curX === tx && curY === ty && curZ === tz
}

console.log(
  mergeTriplets(
    [
      [2, 5, 3],
      [1, 8, 4],
      [1, 7, 5],
    ],
    [2, 7, 5]
  )
)
// 输出：true
// 解释：执行下述操作：
// - 选择第一个和最后一个三元组 [[2,5,3],[1,8,4],[1,7,5]] 。更新最后一个三元组为 [max(2,1), max(5,7), max(3,5)] = [2,7,5] 。triplets = [[2,5,3],[1,8,4],[2,7,5]]
// 目标三元组 [2,7,5] 现在是 triplets 的一个元素。

export {}

// 如果题目的 max 操作改成二进制位的或操作，那么会有什么不一样呢？
// 提示：或也具有单调递增性，本质和 max 操作差不多。
