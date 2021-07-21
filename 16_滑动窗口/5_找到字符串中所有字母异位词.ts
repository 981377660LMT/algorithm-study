// 注意lackNum的条件
// 关键是lackMap与lackNum记录
const findAnagrams = (s: string, p: string): number[] => {
  if (s.length < p.length) return []
  const res: number[] = []

  const lackMap = new Map<string, number>()
  let lackNum = p.length
  let l = 0
  let r = p.length - 1

  // 初始化需要的字符串与对应个数
  for (const letter of p) {
    const count = lackMap.get(letter) || 0
    lackMap.set(letter, count + 1)
  }

  // 初始化开始滑动窗口的缺失值，缺失值为负表示多了
  for (let i = l; i <= r; i++) {
    const letter = s[i]
    if (lackMap.has(letter)) {
      const count = lackMap.get(letter)!
      lackMap.set(letter, count - 1)
      // 注意lackNum的条件 count>0 因为会有重复的字符
      if (count > 0) lackNum--
    }
  }

  if (lackNum === 0) res.push(0)
  console.log(l, r, lackNum, lackMap)
  // 开始移动滑动窗口
  while (r < s.length - 1) {
    l++
    r++
    const preL = s[l - 1]
    const curR = s[r]
    if (lackMap.has(preL)) {
      const count = lackMap.get(preL)!

      if (count >= 0) lackNum++
      lackMap.set(preL, count + 1)
    }

    if (lackMap.has(curR)) {
      const count = lackMap.get(curR)!
      if (count > 0) lackNum--
      lackMap.set(curR, count - 1)
    }
    console.log(l, r, lackMap, lackNum)
    lackNum === 0 && res.push(l)
  }

  // while (r < s.length) {}
  // console.log(lackMap, lackNum)
  return res
}

// console.log(findAnagrams('cbaebabacd', 'abc'))
// // [0,6]
// console.log(findAnagrams('baa', 'aa'))
console.log(findAnagrams('aaaaaaaaaaaaaaabaaaaaaaaaaaa', 'aaaa'))

export {}
