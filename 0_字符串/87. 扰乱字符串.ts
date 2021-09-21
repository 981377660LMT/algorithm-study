/**
 * @param {string} s1
 * @param {string} s2
 * @return {boolean}
 */
var isScramble = function (s1: string, s2: string): boolean {
  const memo = new Map<string, boolean>()
  const inner = (s1: string, s2: string): boolean => {
    if (s1 === s2) return true

    const key = `${s1}#${s2}`
    if (memo.has(key)) return memo.get(key)!

    const sortedS1 = s1.split('').sort().join('')
    const sortedS2 = s2.split('').sort().join('')
    if (sortedS1 !== sortedS2) {
      memo.set(key, false)
      return false
    }

    for (let i = 1; i < s1.length; i++) {
      const pre1 = s1.slice(0, i)
      const pre2 = s2.slice(0, i)
      // 取后面长为i的
      const pre3 = s2.slice(s2.length - i)

      const suffix1 = s1.slice(i)
      const suffix2 = s2.slice(i)
      // 取前面长为s2.length - i的
      const suffix3 = s2.slice(0, s2.length - i)

      if (inner(pre1, pre2) && inner(suffix1, suffix2)) {
        memo.set(key, true)
        return true
      }

      // 反转比较,取后面i个
      if (inner(pre1, pre3) && inner(suffix1, suffix3)) {
        memo.set(key, true)
        return true
      }
    }

    memo.set(key, false)
    return false
  }
  return inner(s1, s2)
}

console.log(isScramble('great', 'rgeat'))
