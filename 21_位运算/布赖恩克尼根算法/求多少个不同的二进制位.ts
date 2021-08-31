/**
 * 求解二进制表示中有多少位不相同
 * @param {Number} a
 * @param {Number} b
 */
function getDiffBytes(a: number, b: number) {
  let count = 0,
    n = a ^ b

  while (n) {
    ++count
    n = n & (n - 1)
  }

  return count
}

/**
 * 测试代码
 */

console.log(getDiffBytes(1, 1))
console.log(getDiffBytes(3, 1))
