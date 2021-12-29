import { nextPermutation } from '../../../../12_贪心算法/经典题/排列/31. 下一个排列'

/**
 * @param {number} n
 * @return {number}
   下一个排列
 */
const nextGreaterElement = function (n: number): number {
  return Number(nextPermutation(n.toString().split('').map(Number)).join(''))
}

console.log(nextGreaterElement(12))
// 21
console.log(nextGreaterElement(101))
// -1
console.log(nextGreaterElement(2147483486))
