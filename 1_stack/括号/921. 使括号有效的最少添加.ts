/**
 * @param {string} s
 * @return {number}
 * 我们需要添加最少的括号（ '(' 或是 ')'，可以在任何位置），
 * 以使得到的括号字符串有效。
 */
function minAddToMakeValid1(s: string): number {
  const stack: string[] = []
  let res = 0

  for (const char of s) {
    if (char === '(') {
      stack.push(char)
    } else if (char === ')' && stack.length > 0) {
      stack.pop()
    } else {
      res++
    }
  }

  return res + stack.length
}

console.log(minAddToMakeValid1('((('))

// 输出：3
// 碰到左括号 无条件入栈
// 碰到右括号进判断，如果空栈，计数器+1；否则pop栈顶
// 最后判断，如果读完全部的S栈内不是空的，则计数器加上栈的长度。。。
let minAddToMakeValid2 = function (s: string): number {
  // 记录多余的个数
  let open = 0
  let close = 0
  for (const char of s) {
    if (char === '(') open++
    else if (!open) close++
    else open--
  }

  return open + close
}
