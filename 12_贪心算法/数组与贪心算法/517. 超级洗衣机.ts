/**
 * @param {number[]} machines
 * @return {number}
 * 你可以选择任意 m （1 ≤ m ≤ n） 台洗衣机，与此同时将每台洗衣机的一件衣服送到相邻的一台洗衣机
 * 请给出能让所有洗衣机中剩下的衣物的数量相等的最少的操作步数。
 * 如果不能使每台洗衣机中衣物的数量相等，则返回 -1。
 * @link
 * https://leetcode-cn.com/problems/super-washing-machines/comments/29839
 */
var findMinMoves = function (machines: number[]): number {
  const sum = machines.reduce((pre, cur) => pre + cur)
  const target = sum / machines.length
  if (target % 1) return -1

  let needToMoveForOne = 0
  let needToMoveTotal = 0
  let max = 0
  machines.forEach(machine => {
    needToMoveForOne = machine - target
    needToMoveTotal += needToMoveForOne
    max = Math.max(max, needToMoveForOne, Math.abs(needToMoveTotal))
  })

  return max
}

console.log(findMinMoves([1, 0, 5]))
