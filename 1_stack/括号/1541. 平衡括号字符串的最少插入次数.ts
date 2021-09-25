/**
 * @param {string} s
 * @return {number}
 * 一个括号字符串被称为平衡的当它满足：

任何左括号 '(' 必须对应两个连续的右括号 '))' 。
左括号 '(' 必须在对应的连续两个右括号 '))' 之前。

 */
const minInsertions = function (s: string): number {
  let res = 0,
    left = 0

  for (let i = 0; i < s.length; i++) {
    // 左括号入栈
    if (s[i] === '(') {
      left++
    } else {
      // 右括号判断:每次要凑齐两个
      if (i + 1 < s.length && s[i + 1] === ')') i++
      // 缺少第二个右括号就添加一个
      else res++

      // 中途栈为空：刚才的两个右括号抵消一个左括号
      if (left > 0) left--
      else res++ // 缺少左括号就添加一个
    }
  }

  // 结束非空栈
  res += left * 2
  return res
}

console.log(minInsertions('(()))'))
// 输出：1
// 解释：第二个左括号有与之匹配的两个右括号，
// 但是第一个左括号只有一个右括号。
// 我们需要在字符串结尾额外增加一个 ')' 使字符串变成平衡字符串 "(())))" 。

// 右括号不可能单独出现。
// 遇到单独出现的右括号就必须补齐一个使其成对，然后尝试去抵消一个左括号。
