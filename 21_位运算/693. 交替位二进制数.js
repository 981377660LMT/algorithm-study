/**
 * @param {number} n
 * @return {boolean}
 * 给定一个正整数，检查它的二进制表示是否总是 0、1 交替出现
 * 即检查二进制表示中相邻两位的数字永不相同
 */
function hasAlternatingBits(n) {
  const tmp = n ^ (n >> 1)
  return (tmp & (tmp + 1)) === 0 // 检查结果是否全为1
}
