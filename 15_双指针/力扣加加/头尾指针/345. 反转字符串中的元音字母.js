/**
 * @param {string} s
 * @return {string}
 */
var reverseVowels = function (s) {
  if (s.length <= 1) return s

  const res = s.split('')
  const vowels = new Set(['a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U'])
  let i = 0
  let j = s.length - 1

  while (i < j) {
    while (i < j && !vowels.has(res[i])) i++
    while (i < j && !vowels.has(res[j])) j--
    ;[res[i], res[j]] = [res[j], res[i]]
    i++
    j--
  }

  return res.join('')
}

console.log(reverseVowels('leetcode'))
