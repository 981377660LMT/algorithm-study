/**
 Do not return anything, modify s in-place instead.
 使用 O(1) 额外空间复杂度的原地解法
 */
function reverseWords(strArr: string[]): void {
  // 翻转
  reverse(strArr, 0, strArr.length - 1)

  let start = 0

  // 注意最后一个i===strArr.length
  for (let i = 0; i <= strArr.length; i++) {
    if (strArr[i] === ' ' || i === strArr.length) {
      // 翻转单词
      reverse(strArr, start, i - 1)
      start = i + 1
    }
  }

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

reverseWords(['t', 'h', 'e', ' ', 's', 'k', 'y', ' ', 'i', 's', ' ', 'b', 'l', 'u', 'e'])
