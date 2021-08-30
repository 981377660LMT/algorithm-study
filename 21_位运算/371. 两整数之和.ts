/**
 * @param {number} a
 * @param {number} b
 * @return {number}
 * 不使用运算符 + 和 - ​​​​​​​，计算两整数 ​​​​​​​a 、b ​​​​​​​之和。
 * @summary
 * 1. 异或是一种不进位的加减法
 * 2. 求与之后左移一位来可以表示进位
 * 3. 把这两个数"相加"就得到结果
 */
const getSum = (a: number, b: number): number => {
  // (a & b) << 1是进位的结果，a ^ b是不考虑进位直接相加的结果
  // 最后如果进位运算结果等于0，说明没有进位了，可以直接返回不考虑进位直接相加的结果
  return a === 0 ? b : getSum((a & b) << 1, a ^ b)
}

console.log(getSum(1, 2))
