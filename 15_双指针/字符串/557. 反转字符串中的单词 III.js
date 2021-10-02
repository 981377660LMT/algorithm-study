/**
 * @param {string} s
 * @return {string}
 * 请尝试使用 O(1) 额外空间复杂度的原地解法。
 */
const reverseWords = s => {
  return s
    .split(' ')
    .map(w => w.split('').reverse().join(''))
    .join(' ')
}

console.log(reverseWords('the sky is blue'))
// 输入："Let's take LeetCode contest"
// 输出："s'teL ekat edoCteeL tsetnoc"
