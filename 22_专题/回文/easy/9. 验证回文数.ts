/**
 * @param {number} x
 * @return {boolean}
 * @description 倒序计算
 */
const isPalindrome = function (x: number): boolean {
  if (x < 0) return false
  let res = 0
  let i = x
  while (i >= 1) {
    res = res * 10 + (i % 10)
    i = ~~(i / 10)
  }

  return res === x
}

console.log(isPalindrome(121))

export default 1
