const mask: Record<string, number> = {
  a: 1 << 0,
  e: 1 << 1,
  i: 1 << 2,
  o: 1 << 3,
  u: 1 << 4,
}
/**
 * @param {string} s
 * @return {number}
 * @summary 滑动窗口（这里是可变滑动窗口）我们需要扩张和收缩窗口大小，而这里不那么容易
 * 关注五个数的状态:每个数的状态为0/1 状态压缩即可
 */
function findTheLongestSubstring(s: string): number {
  let res = 0
  let state = 0b00000
  // 每种状态最早出现的索引
  const visited = new Map<number, number>([[state, -1]])

  for (let i = 0; i < s.length; i++) {
    if (s[i] in mask) state ^= mask[s[i]]
    if (!visited.has(state)) visited.set(state, i)
    else res = Math.max(res, i - visited.get(state)!)
  }

  return res
}

console.log(findTheLongestSubstring('eleetminicoworoep'))
export {}
