import { maxJump } from './maxJump'

// 花园里总共有 n + 1 个水龙头，分别位于 [0, 1, ..., n]
// 请你返回可以灌溉整个花园的 最少水龙头数目
// 如果打开点 i 处的水龙头，可以灌溉的区域为 [i -  ranges[i], i + ranges[i]]
// 请你返回可以灌溉整个花园的 最少水龙头数目 。如果花园始终存在无法灌溉到的地方，请你返回 -1
function minTaps(n: number, ranges: number[]): number {
  const jumps = Array<number>(n + 1).fill(0)
  for (const [index, range] of ranges.entries()) {
    const start = index - range
    const dist = range * 2
    jumps[Math.max(0, start)] = Math.max(jumps[Math.max(0, start)], start + dist)
  }

  const res = maxJump(jumps, n)
  return res === -1 ? -1 : res
}

console.log(minTaps(5, [3, 4, 1, 1, 0, 0]))
console.log(minTaps(3, [0, 0, 0, 0]))
