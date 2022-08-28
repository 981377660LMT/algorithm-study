/* eslint-disable semi-style */

const VOWEL = new Set(['a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U'])

/**
 * @param {string} s
 * @return {string}
 */
function reverseVowels(s) {
  if (s.length <= 1) return s

  const res = s.split('')
  let i = 0
  let j = s.length - 1

  while (i < j) {
    while (i < j && !VOWEL.has(res[i])) i++
    while (i < j && !VOWEL.has(res[j])) j--
    ;[res[i], res[j]] = [res[j], res[i]]
    i++
    j--
  }

  return res.join('')
}

console.log(reverseVowels('leetcode'))
