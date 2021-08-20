/**
 * @param {string} s
 * @param {number} n
 * @return {string}
 * 不能申请额外空间，只能在本串上操作
 * @summary 相当于循环移动字符串
 */
const reverseLeftWords = function (s: string, n: number): string {
  const reverse = (str: string, left: number, right: number) => {
    let strArr = str.split('')
    for (; left < right; left++, right--) {
      ;[strArr[left], strArr[right]] = [strArr[right], strArr[left]]
    }
    return strArr.join('')
  }

  s = reverse(s, 0, n - 1)
  s = reverse(s, n, s.length - 1)
  return reverse(s, 0, s.length - 1)
}

console.log(reverseLeftWords('abcdefg', 2))

export default 1
