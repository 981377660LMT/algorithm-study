/**
 * @param {number[]} weights
 * @param {number} days
 * @return {number}
 */
var shipWithinDays = function (weights, days) {
  const sum = weights.reduce((pre, cur) => pre + cur)
  const possible = mid => {
    let curWeight = 0
    let needDays = 1
    for (const w of weights) {
      if (w > mid) return false
      if (curWeight + w > mid) {
        curWeight = 0
        needDays++
      }
      curWeight += w
    }
    return needDays <= days
  }

  let l = 1
  let r = sum
  while (l <= r) {
    const mid = (l + r) >> 1
    if (possible(mid)) r = mid - 1
    else l = mid + 1
  }

  return l
}
