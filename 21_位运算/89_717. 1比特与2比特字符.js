/**
 * @param {number[]} bits
 * @return {boolean}
 * 第一种字符可以用一比特0来表示。第二种字符可以用两比特(10 或 11)来表示
 * 现给一个由若干比特组成的字符串。
 * 问最后一个字符是否必定为一个一比特字符。给定的字符串总是由0结束。
 * @summary
 * 霍夫曼编码，见1直接跳两位,见0跳一位
 * 当遇到1时，这个1一定会把下一个0或1吃掉，
 * 因此这时需要跳过下一个。
 * 如果能遍历到最后一个0，就说明成功了
 */
function isOneBitCharacter(bits) {
  // return /^(10|11|0)*0$/.test(bits.join(''))
  let i = 0

  while (i < bits.length - 1) {
    if (bits[i] === 1) i += 2
    else i += 1
  }

  return bits[i] === 0
}

console.log(isOneBitCharacter([1, 0, 0]))
console.log(isOneBitCharacter([1, 1, 1, 0]))
