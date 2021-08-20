/**
 * @param {string} s
 * @return {string}
 */
const replaceSpace = function (s: string): string {
  const strArr = Array.from(s)
  let count = 0
  // 计算空格数量
  for (let i = 0; i < strArr.length; i++) {
    if (strArr[i] === ' ') {
      count++
    }
  }

  // # 将空格改成%20 使得字符串总长增长 2n，n为原本空格数量。
  // # 所以记录空格数量就可以得到目标字符串的长度
  let left = strArr.length - 1
  let right = strArr.length + count * 2 - 1
  while (left >= 0) {
    if (strArr[left] === ' ') {
      strArr[right--] = '0'
      strArr[right--] = '2'
      strArr[right--] = '%'
      left--
    } else {
      strArr[right--] = strArr[left--]
    }
  }

  return strArr.join('')
}

console.log(replaceSpace('We are happy.'))
// 输出："We%20are%20happy."
