/**
 * @param {number} num
 * @return {number}
 * 给你一个 正 整数 num ，输出它的补数。补数是对该数的二进制表示取反
 * 即将num二进制各位由1变成0，0变成1，由此想到将各位与1做异或操作即可
 * @summary
 * 十进制整数的反码
 * 找到与num二进制有效位（没有前导零位）个数相同且都是1的数(num高一位减一即可)
 * @warning
 * `2147483647 时失效`
 *
 */
const findComplement = function (num: number): number {
  if (num === 0) return 1

  let mask = 1
  while (mask <= num) {
    mask <<= 1
  }

  return num ^ (mask - 1)
}

console.log(findComplement(5))
// 5 的二进制表示为 101（没有前导零位），其补数为 010。所以你需要输出 2 。

export default 1
