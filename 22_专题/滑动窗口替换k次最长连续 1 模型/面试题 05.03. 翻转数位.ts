import { fix } from './最长连续 1 模型'

/**
 * @param {number} num
 * @return {number}
 * 给定一个32位整数 num，你可以将一个数位从0变为1。
 * 请编写一个程序，找出你能够获得的最长的一串1的长度。
 */
const reverseBits = function (num: number): number {
  let str = (num >>> 0).toString(2) // 负数时左边全补1 相当于padStart(32, '1')
  num > 0 && (str = str.padStart(32, '0'))
  // num < 0 && (str = str.padStart(32, '1'))
  console.log(str)
  return fix(str, '1', 1)
}

console.log(reverseBits(1775))
console.log(reverseBits(2147483647))
console.log(reverseBits(-1))

export {}
console.log(Number(2147483647).toString(2).length)

//
