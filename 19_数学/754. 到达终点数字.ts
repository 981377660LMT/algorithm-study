/* eslint-disable no-param-reassign */

/**
 * 到达终点数字
 * 在一根无限长的数轴上，你站在0的位置。终点在target的位置。
 * 你可以做一些数量的移动 numMoves :
 * - 每次你可以选择向左或向右移动。
 * - 第 i 次移动（从  i == 1 开始，到 i == numMoves ），在选择的方向上走 i 步。
 * 给定整数 target ，返回 到达目标所需的 最小 移动次数(即最小 numMoves ) 。
 *
 * !转化为对1，2，3，4，5....i,添加正负号，使得其和等于target的最小的数目。
 * S-target=2N，S=(i*(i+1))/2
 * !因此就变为找找到最小的i,使得S>=target且S-target为偶数!!!
 * 必须线性查找,不能二分(没有单调性)
 *
 * 829. 连续整数求和
 * !等差数列求和/连续整数求和
 */
function reachNumber(target: number): number {
  if (target < 0) target = -target
  if (target === 0) return 0

  let i = 1
  while (true) {
    const sum = (i * (i + 1)) / 2
    if (sum >= target && (sum - target) % 2 === 0) return i
    i++
  }
}

console.log(reachNumber(3))

export {}
