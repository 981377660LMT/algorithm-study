/**
 * @param {string} s
 * @return {number}
 */
const romanToInt = function (s: string): number {
  const mapper = {
    I: 1,
    V: 5,
    X: 10,
    L: 50,
    C: 100,
    D: 500,
    M: 1000,
  } as Record<string, number>

  let res = 0

  for (let i = 0; i < s.length; i++) {
    mapper[s[i]] < mapper[s[i + 1]] ? (res -= mapper[s[i]]) : (res += mapper[s[i]])
  }

  return res
}
