/**
 * @param {number} n
 * @return {boolean}
 * 二进制表示中相邻两位的数字永不相同
 */
var hasAlternatingBits = function (n) {
  const tmp = n ^ (n >> 1)
  return (tmp & (tmp + 1)) === 0 // 检查结果是否全为1
}
