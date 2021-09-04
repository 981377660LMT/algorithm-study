/**
 * @param {number} n
 * @return {number[]}
 * 格雷编码是一个二进制数字系统，在该系统中，
 * 两个连续的数值仅有一个位数的差异。
 * 只需要返回其中一种。
 * 格雷编码序列必须以 0 开头。
 */
var grayCode = function (n: number): number[] {
  return Array.from({ length: 2 ** n }, (_, i) => i ^ (i >> 1))
}

console.log(grayCode(2))
