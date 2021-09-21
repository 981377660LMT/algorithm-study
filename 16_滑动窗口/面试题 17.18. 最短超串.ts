// 找到长数组中包含短数组所有的元素的最短子数组，其出现顺序无关紧要。
// 返回最短子数组的左端点和右端点，如有多个满足条件的子数组，返回左端点最小的一个。若不存在，返回空数组。
// 6_最小覆盖子串
function shortestSeq(big: number[], small: number[]): number[] {
  if (big.length < small.length) return []
  // 缺少的对应关系,负数代表多了
  const lackMap = new Map<number, number>()
  // 缺少的数量,非负数
  let lackNum = small.length
  let l = 0
  let r = small.length - 1
  let minLength = Infinity
  const res: [number, number] = [Infinity, Infinity]

  // 初始化lackMap与lackNum
  for (const char of small) {
    lackMap.set(char, lackMap.get(char) || 0 + 1)
  }

  for (let i = 0; i < small.length; i++) {
    const char = big[i]
    if (lackMap.has(char)) {
      const count = lackMap.get(char)!
      lackMap.set(char, count - 1)
      // 注意lackNum的条件 count>0 因为会有重复的字符
      if (count > 0) lackNum--
    }
  }

  if (lackNum === 0) return [l, r]

  while (r < big.length - 1) {
    // 不符合条件，扩张右边
    while (lackNum > 0) {
      if (r >= big.length - 1) break
      r++
      const cur = big[r]
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
        res[1] = r
      }
      l++
      const pre = big[l - 1]
      if (lackMap.has(pre)) {
        const count = lackMap.get(pre)!
        if (count >= 0) lackNum++
        lackMap.set(pre, count + 1)
      }
    }
  }

  return minLength === Infinity ? [] : [res[0], res[1]]
}

console.log(shortestSeq([7, 5, 9, 0, 2, 1, 3, 5, 7, 9, 1, 1, 5, 8, 8, 9, 7], [1, 5, 9]))
