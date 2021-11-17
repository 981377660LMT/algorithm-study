/**
 * @param {string} s
 * @return {string}
 * 请尝试使用 O(1) 额外空间复杂度的原地解法。
 * 将每个单词反转+
   将整个字符串反转
   先整体反转再局部反转
 */
const reverseWords = (s: string): string => {
  // 字符串转数组
  const strArr = Array.from(s)

  // 删除多余空格
  // removeExtraSpaces(strArr)
  // 翻转
  reverse(strArr, 0, strArr.length - 1)

  let start = 0

  for (let i = 0; i <= strArr.length; i++) {
    if (strArr[i] === ' ' || i === strArr.length) {
      // 翻转单词
      reverse(strArr, start, i - 1)
      start = i + 1
    }
  }

  return strArr.join('')

  // 翻转从 start 到 end 的字符
  function reverse(strArr: string[], left: number, right: number) {
    while (left < right) {
      // 交换
      ;[strArr[left], strArr[right]] = [strArr[right], strArr[left]]
      left++
      right--
    }
  }
}

console.log(reverseWords('the sky is blue'))
// "blue is sky the"
export default 1
