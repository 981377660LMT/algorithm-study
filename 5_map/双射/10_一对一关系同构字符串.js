/**
 * @param {string} s
 * @param {string} t
 * @return {boolean}
 */
function isIsomorphic(s, t) {
  if (s.length !== t.length) return false

  const sMap = new Map()
  const tMap = new Map()

  for (let i = 0; i < s.length; i++) {
    const letterS = s[i]
    const letterT = t[i]
    if (
      (sMap.has(letterS) && sMap.get(letterS) !== letterT) ||
      (tMap.has(letterT) && tMap.get(letterT) !== letterS)
    )
      return false
    sMap.set(letterS, letterT)
    tMap.set(letterT, letterS)
  }

  return true
}

console.log(isIsomorphic('paper', 'title'))
