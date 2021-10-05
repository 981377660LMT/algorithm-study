/**
 * @param {string} s
 * @return {number}
 * 记住每个元素第一次出现的下标即可
 */
const maxLengthBetweenEqualCharacters = function (s: string): number {
  let res = -1
  const firstIndexRecord = new Map<string, number>()
  for (let i = 0; i < s.length; i++) {
    const cur = s[i]
    if (firstIndexRecord.has(cur)) res = Math.max(res, i - firstIndexRecord.get(cur)! - 1)
    else firstIndexRecord.set(cur, i)
  }

  return res
}

console.log(maxLengthBetweenEqualCharacters('abca'))
