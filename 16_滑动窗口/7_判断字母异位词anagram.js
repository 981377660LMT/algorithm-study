/**
 * @param {string} s
 * @param {string} t
 * @return {boolean}
 * @description lackNum+lackMap实现
 */
var isAnagram = function (s, t) {
  if (s.length !== t.length) return false
  let lackNum = s.length
  const lackMap = new Map()

  for (const letter of s) lackMap.set(letter, lackMap.get(letter) + 1 || 1)
  for (const letter of t) {
    if (lackMap.has(letter)) {
      const count = lackMap.get(letter)
      if (count > 0) lackNum--
      lackMap.set(letter, count - 1)
    }
  }

  return lackNum === 0
}

console.log(isAnagram('anagram', 'nagaram'))
