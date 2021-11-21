function numDifferentIntegers(word: string): number {
  const digits = word.match(/\d+/g)
  if (digits == void 0) return 0
  return new Set(digits.map(str => lstrip(str, '0'))).size
}

// 去除前导
function lstrip(raw: string, remove: string) {
  const removeLeft = new RegExp(`^${remove}+`)
  return raw.replace(removeLeft, '')
}
// 输入：word = "a123bc34d8ef34"
// 输出：3
// 解释：不同的整数有 "123"、"34" 和 "8" 。注意，"34" 只计数一次。

console.log(numDifferentIntegers('a123bc34d8ef34'))
