/**
 * @param {string} s
 * @return {boolean}
 * 给定一个只包含三种字符的字符串：（ ，） 和 *，写一个函数来检验这个字符串是否为有效字符串。
 * * 可以被视为单个右括号 ) ，或单个左括号 ( ，或一个空字符串。
 * @summary 某些位置的括号可以任意转换的问题
 */
function checkValidString(s: string): boolean {
  return check(s) && check(s, true)
  function check(str: string, reversed = false): boolean {
    reversed && (str = str.split('').reverse().join(''))
    let count = 0
    let flag = reversed ? -1 : 1

    for (const char of str) {
      switch (char) {
        case '(':
          count += flag
          break
        case ')':
          count -= flag
          break
        case '*':
          count++
          break
        default:
          break
      }

      if (count < 0) return false
    }

    return true
  }
}

console.log(checkValidString('(*))'))
console.log(checkValidString('(*)'))

export {}

// 遍历两次，第一次顺序，第二次逆序。
// 第一次遇到左括号加一，右括号减一，星号加一，最后保证cnt >= 0,也就是可以保证产生的左括号足够
// 第二次遇到右括号加一，左括号减一，星号加一，最后保证cnt >= 0,也就是可以保证产生的右括号足够
