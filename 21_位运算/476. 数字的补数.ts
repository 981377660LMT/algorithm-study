/**
 * @param {number} num
 * @return {number}
 * 给你一个 正 整数 num ，输出它的补数。补数是对该数的二进制表示取反
 * 即将num二进制各位由1变成0，0变成1，由此想到将各位与1做异或操作即可
 * @summary
 * 十进制整数的反码
 * 找到与num二进制有效位（没有前导零位）个数相同且都是1的数(num高一位减一即可)
 */
const findComplement = function (num: number): number {
  // 计数前导零32
  // "clz32" 是CountLeadingZeroes32
  // 此函数用于获取数字的32位表示形式中出现的前导零位的数量
  const clz = Math.clz32(num)
  return (~num << clz) >>> clz
}

var a = 32776 // 00000000000000001000000000001000 (16个前导0)
Math.clz32(a) // 16

var b = ~32776 // 11111111111111110111111111110111 (对32776取反, 0个前导0)
Math.clz32(b) // 0 (相当于0个前导1)

console.log(findComplement(5))
// 101 =>10
export default 0

// ~是取补码(反码加1)
