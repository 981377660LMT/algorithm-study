// 返回 s 中涵盖 t 所有字符的最小子串。如果 s 中不存在涵盖 t 所有字符的子串，则返回空字符串 "" 。
// 需要 o(n) 时间内解决
// 如果 s 中存在这样的子串，我们保证它是唯一的答案。
// 关键是lackMap与lackNum记录
const minWindow = (s: string, t: string): string => {
  if (s.length < t.length) return ''
  // 缺少的对应关系,负数代表多了
  const lackMap = new Map<string, number>()
  // 缺少的数量,非负数
  let lackNum = t.length
  let l = 0
  let r = t.length - 1
  let minLength = Infinity
  const res: [number, number] = [Infinity, Infinity]

  // 初始化lackMap与lackNum
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
    // 不符合条件，扩张右边
    while (lackNum > 0) {
      if (r > s.length - 1) break
      r++
      const cur = s[r]
      if (lackMap.has(cur)) {
        const count = lackMap.get(cur)!
        if (count > 0) lackNum--
        lackMap.set(cur, count - 1)
      }
    }

    // 符合条件，更新答案，开始收缩左边
    while (lackNum === 0) {
      if (r - l + 1 < minLength) {
        minLength = r - l + 1
        res[0] = l
        res[1] = r + 1
      }
      l++
      const pre = s[l - 1]
      if (lackMap.has(pre)) {
        const count = lackMap.get(pre)!
        if (count >= 0) lackNum++
        lackMap.set(pre, count + 1)
      }
    }
  }

  return minLength === Infinity ? '' : s.slice(res[0], res[1])
}

console.log(minWindow('ADOBECODEBANC', 'ABC'))
console.log(minWindow('ab', 'A'))
// console.log(minWindow('ABC', 'ABC'))
// console.log(minWindow('ABDCABC', 'ABC'))

export {}
