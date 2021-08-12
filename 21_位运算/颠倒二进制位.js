// 输入是一个长度为 32 的二进制字符串
var reverseBits = function (n) {
  var result = 0
  var count = 32

  while (count--) {
    result *= 2
    result += n & 1
    n = n >> 1
  }
  return result
}

console.log('4' >> 1)
