import { DefaultDict } from '../../../5_map/DefaultDict'
import { ACAutoMatonLegacy } from './ACAutoMatonMapLegacy'

// !因为字符串总长度太大(1e6),MLE了
const INF = 2e15
function longestValidSubstring(word: string, forbidden: string[]): number {
  const acm = new ACAutoMatonLegacy()
  const minLen = new DefaultDict(() => INF)
  forbidden.forEach((pattern, id) => {
    acm.insert(id, pattern, pos => {
      minLen.set(pos, Math.min(minLen.get(pos), pattern.length))
    })
  })

  acm.build(false, (pre, cur) => {
    minLen.set(cur, Math.min(minLen.get(cur), minLen.get(pre)))
  })

  let res = 0
  let left = 0
  let state = 0
  for (let right = 0; right < word.length; right++) {
    state = acm.move(state, word[right])
    const start = right - minLen.get(state) + 1
    left = Math.max(left, start + 1)
    res = Math.max(res, right - left + 1)
  }
  return res
}

export {}
