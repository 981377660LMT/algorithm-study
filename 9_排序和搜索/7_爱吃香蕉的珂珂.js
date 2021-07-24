/**
 * @param {number[]} piles
 * @param {number} h
 * @return {number}
 * @description 注意time(k)是单调的 可以用二分搜索
 * 不断收拢右侧边界
 */
var minEatingSpeed = function (piles, h) {
  const max = Math.max.apply(null, piles)

  const calTime = (piles, rate) => piles.reduce((pre, cur) => pre + Math.ceil(cur / rate), 0)

  let leftPoint = 1
  let rightPoint = max
  // 终止循环条件是l===r
  while (leftPoint < rightPoint) {
    const mid = Math.floor((leftPoint + rightPoint) / 2)
    const time = calTime(piles, mid)
    if (time <= h) {
      rightPoint = mid
    } else {
      leftPoint = mid + 1
    }
  }

  return rightPoint
}

console.log(minEatingSpeed([30, 11, 23, 4, 20], 5))
