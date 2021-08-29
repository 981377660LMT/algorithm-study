/**
 * @param {string} s
 * @return {boolean}
 * @description 只考虑字母和数字字符，可以忽略字母的大小写。
 * 如果两个指针的元素不相同，则直接返回 false,
   如果两个指针的元素相同，我们同时更新头尾指针，循环。 直到头尾指针相遇。
 */
var isPalindrome = function (s: string): boolean {
  const isLetterOrNumber = (code: number) => {
    return (
      (code >= 48 && code <= 57) || // numbers
      (code >= 65 && code <= 90) || // uppercase
      (code >= 97 && code <= 122) // lowercase
    )
  }

  const toLowercase = (code: number) => {
    if (code >= 65 && code <= 90) return code + 32
    else return code
  }

  let l = 0
  let r = s.length - 1
  while (l < r) {
    const lCode = s.codePointAt(l)!
    const rCode = s.codePointAt(r)!

    if (!isLetterOrNumber(lCode)) {
      l++
      continue
    }

    if (!isLetterOrNumber(rCode)) {
      r--
      continue
    }

    if (toLowercase(lCode) !== toLowercase(rCode)) return false

    l++
    r--
  }

  return true
}

console.log(isPalindrome('A man, a plan, a canal: Panama'))
