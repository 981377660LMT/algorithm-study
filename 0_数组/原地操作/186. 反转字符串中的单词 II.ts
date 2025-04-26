// https://leetcode.cn/problems/reverse-words-in-a-string-ii/description/
// 给你一个字符数组 s ，反转其中 单词 的顺序。
// 单词 的定义为：单词是一个由非空格字符组成的序列。s 中的单词将会由单个空格分隔。
// 必须设计并实现 原地 解法来解决此问题，即不分配额外的空间。
// !先整体反转， 再反转所有空格之间的单词

/**
 Do not return anything, modify s in-place instead.
 */
function reverseWords(s: string[]): void {
  const n = s.length
  _reverse(s, 0, n - 1)

  let start = 0
  for (let i = 0; i < n; i++) {
    if (s[i] === ' ') {
      _reverse(s, start, i - 1)
      start = i + 1
    }
  }
  _reverse(s, start, n - 1)
}

function _reverse(arr: any[], left: number, right: number): void {
  while (left < right) {
    const tmp = arr[left]
    arr[left] = arr[right]
    arr[right] = tmp
    left++
    right--
  }
}

export {}
