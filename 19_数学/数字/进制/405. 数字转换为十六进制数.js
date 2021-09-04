/**
 * @param {number} num
 * @return {string}
 * 对于负整数，我们通常使用 补码运算 方法
 */
var toHex = function (num) {
  const arr = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f']
  if (num === 0) return '0'
  if (num < 0) num += 2 ** 32
  const res = []
  while (num) {
    const digit = num % 16
    res.push(arr[digit])
    num = ~~(num / 16)
  }
  return res.reverse().join('')
}

console.log(toHex(26))
console.log(toHex(-1))
