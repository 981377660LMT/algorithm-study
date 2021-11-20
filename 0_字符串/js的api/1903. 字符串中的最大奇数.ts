function largestOddNumber(num: string): string {
  let lastIndex = -1

  for (let i = num.length - 1; ~i; i--) {
    if ((Number(num[i]) & 1) === 1) {
      lastIndex = i
      break
    }
  }

  return num.slice(0, lastIndex + 1)
}

// console.log(largestOddNumber('52'))
console.log(largestOddNumber('35427'))
// // 输出："35427"
// // 解释："35427" 本身就是一个奇数。
