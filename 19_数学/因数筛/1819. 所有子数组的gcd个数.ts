import { logTrick } from '../../21_位运算/logTrick/logTrick'
import { gcd } from '../数论/扩展欧几里得/gcd'

/**
 * 统计数组所有子数组的 gcd 的不同个数，复杂度 O(n*log^2max)
 */
function countGcdOfAllSubarray(nums: number[]): number {
  return logTrick(nums, gcd).size
}

if (require.main === module) {
  console.log(countGcdOfAllSubarray([6, 10, 15]))
  console.log(countGcdOfAllSubarray([5, 5, 5]))
}
