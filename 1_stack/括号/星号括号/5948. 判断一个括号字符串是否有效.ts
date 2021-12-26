// 678. 有效的括号字符串
function canBeValid(s: string, locked: string): boolean {
  const n = s.length
  if (n & 1) return false
  const sb: string[] = []
  for (let i = 0; i < locked.length; i++) {
    if (locked[i] === '0') sb.push('*')
    else sb.push(s[i])
  }

  const sWithStar = sb.join('')
  return check(sWithStar) && check(sWithStar, true)

  function check(str: string, reversed = false) {
    reversed && (str = str.split('').reverse().join(''))
    let count = 0
    let score = reversed ? -1 : 1

    for (const char of str) {
      switch (char) {
        case '(':
          count += score
          break
        case ')':
          count += -score
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

console.log(canBeValid('))()))', '010100'))
console.log(
  canBeValid('())()))()(()(((())(()()))))((((()())(())', '1011101100010001001011000000110010100101')
)
console.log(canBeValid(')', '0'))
console.log(canBeValid(')(', '00'))
