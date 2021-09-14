/**
 * @param {number[]} flowerbed
 * @param {number} n
 * @return {boolean}
 * 花不能种植在相邻的地块上
 * 能否在不打破种植规则的情况下种入 n 朵花
 */
var canPlaceFlowers = function (flowerbed, n) {
  flowerbed.unshift(0)
  flowerbed.push(0)
  let res = 0
  for (let i = 1; i < flowerbed.length - 1; i++) {
    const canPlace = flowerbed[i - 1] === 0 && flowerbed[i] === 0 && flowerbed[i + 1] === 0
    if (canPlace) {
      flowerbed[i] = 1
      res++
    }
  }
  return res >= n
}

console.log(canPlaceFlowers([1, 0, 0, 0, 1], 1))
console.log(canPlaceFlowers([1, 0, 0, 0, 1], 2))
