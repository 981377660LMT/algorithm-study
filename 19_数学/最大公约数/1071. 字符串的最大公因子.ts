// 返回最长字符串 X，要求满足 X 能除尽 str1 且 X 能除尽 str2。
const gcdOfStrings = (str1: string, str2: string) => {
  // 不相等则间接说明了不存在字符串X
  if (str1 + str2 !== str2 + str1) {
    return ''
  }

  // 最大公约数计算公式
  const gcd = (num1: number, num2: number): number =>
    // 利用辗转相除法来计算最大公约数，即字符串X在字符串str1中截止的索引位置
    num2 === 0 ? num1 : gcd(num2, num1 % num2)

  // 截取匹配的字符串
  return str1.substring(0, gcd(str1.length, str2.length))
}

console.log(gcdOfStrings('ABCABC', 'ABC'))
