import { BITArray } from './BIT'

// 2528. 最大化城市的最小供电站数目
function maxPower(stations: number[], r: number, k: number): number {
  const n = stations.length
  let left = 1
  let right = 2e15
  while (left <= right) {
    const mid = Math.floor((left + right) / 2)
    if (check(mid)) left = mid + 1
    else right = mid - 1
  }

  return right

  // !可以用滑窗优化
  function check(mid: number): boolean {
    const bit = new BITArray(stations)
    let curK = k
    for (let i = 0; i < n; i++) {
      const cur = bit.queryRange(Math.max(0, i - r), Math.min(i + r + 1, n))
      if (cur < mid) {
        const diff = mid - cur
        bit.add(Math.min(i + r, n - 1), diff)
        curK -= diff
        if (curK < 0) return false
      }
    }
    return true
  }
}

// stations = [1,2,4,5,0], r = 1, k = 2
console.log(maxPower([1, 2, 4, 5, 0], 1, 2)) // 5

export {}
