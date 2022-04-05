// 为什么ts比python慢

import { memo } from '../../5_map/memo'

// js过不了 字符串存储开销太大了?
const cache = new Map<string, any>()

function numberOfWays(s: string): number {
  let dfs = (index: number, pre: string, count: number): number => {
    if (count === 3) return 1
    if (index === s.length) return 0

    const key = `${index}#${pre}#${count}`
    if (cache.has(key)) return cache.get(key)!

    let res = 0

    if (pre === '0' && s[index] === '1') {
      res += dfs(index + 1, s[index], count + 1)
    } else if (pre === '1' && s[index] === '0') {
      res += dfs(index + 1, s[index], count + 1)
    }

    res += dfs(index + 1, pre, count)

    cache.set(key, res)
    return res
  }

  let res = 0
  for (let i = 0; i < s.length; i++) {
    res += dfs(i, s[i], 1)
  }

  cache.clear()
  return res
}

export {}
