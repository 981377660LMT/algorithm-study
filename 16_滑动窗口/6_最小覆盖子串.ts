// 返回 s 中涵盖 t 所有字符的最小子串。如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 "" 。
// 需要 o(n) 时间内解决
// 如果 s 中存在这样的子串，我们保证它是唯一的答案。
// 关键是lackMap与lackNum记录
const minWindow = (s: string, t: string): string => {
  if (s.length < t.length) return ''
  const lackMap = new Map<string, number>()
  let lackNum = t.length
  let l = 0
  let r = t.length - 1
  let minLength = Infinity
  const memo: [number, number] = [Infinity, Infinity]

  // 初始化
  for (const letter of t) {
    const count = lackMap.get(letter) || 0
    lackMap.set(letter, count + 1)
  }

  for (let i = l; i <= r; i++) {
    const letter = s[i]
    if (lackMap.has(letter)) {
      const count = lackMap.get(letter)!
      lackMap.set(letter, count - 1)
      // 注意lackNum的条件 count>0 因为会有重复的字符
      if (count > 0) lackNum--
    }
  }

  if (lackNum === 0) return s.slice(l, r + 1)

  while (r < s.length - 1) {
    while (lackNum > 0) {
      if (r > s.length - 1) break
      r++
      const curR = s[r]
      if (lackMap.has(curR)) {
        const count = lackMap.get(curR)!
        if (count > 0) lackNum--
        lackMap.set(curR, count - 1)
      }
    }
    console.log(l, r, lackMap, lackNum, minLength)

    while (lackNum === 0) {
      if (r - l + 1 < minLength) {
        minLength = r - l + 1
        memo[0] = l
        memo[1] = r + 1
      }
      l++
      const preL = s[l - 1]
      if (lackMap.has(preL)) {
        const count = lackMap.get(preL)!
        if (count >= 0) lackNum++
        lackMap.set(preL, count + 1)
      }
    }
  }

  return minLength === Infinity ? '' : s.slice(memo[0], memo[1])
}

console.log(minWindow('ADOBECODEBANC', 'ABC'))
console.log(minWindow('ab', 'A'))
// console.log(minWindow('ABC', 'ABC'))
// console.log(minWindow('ABDCABC', 'ABC'))

export {}
