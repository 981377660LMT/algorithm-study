// 请通过连续相同字符的计数来压缩一个字符串。个数为1时不用包含数字。

/**
 * @param {string} str
 * @return {string}
 */
function compress(str: string): string {
  // your code here
  let count = 0
  const sb: string[] = []
  for (let i = 0; i < str.length; i++) {
    count++
    if (str[i] !== str[i + 1]) {
      sb.push(str[i])
      count > 1 && sb.push(count.toString())
      count = 0
    }
  }

  return sb.join('')
}

if (require.main === module) {
  console.log(compress('a')) // 'a'
  console.log(compress('aa')) // 'a2'
  console.log(compress('aaa')) // 'a3'
  console.log(compress('aaab')) // 'a3b'
  console.log(compress('aaabb')) // 'a3b2'
  console.log(compress('aaabba')) // 'a3b2a'
}
