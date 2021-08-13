/**
 * @param {number[]} houses
 * @param {number[]} heaters
 * @return {number}
 * @description 在加热器的加热半径范围内的每个房屋都可以获得供暖。
 * 给出位于一条水平线上的房屋 houses 和供暖器 heaters 的位置，请你找出并返回可以覆盖所有房屋的最小加热半径。
 * @summary 能力检测二分/最左二分
 */
const findRadius = function (houses: number[], heaters: number[]): number {
  houses.sort((a, b) => a - b)
  heaters.sort((a, b) => a - b)

  const canCover = (radius: number): boolean => {
    let count = 0
    // 当前供暖器是否能覆盖到当前的房子

    for (let i = 0; i < heaters.length; i++) {
      const l = heaters[i] - radius
      const r = heaters[i] + radius
      while (count < houses.length && houses[count] <= r && houses[count] >= l) {
        count += 1
      }
      if (count === houses.length) return true
    }

    return count === houses.length
  }

  let l = 0
  let r = Math.max(houses[houses.length - 1] - heaters[0], heaters[heaters.length - 1] - houses[0])

  while (l <= r) {
    const mid = (l + r) >> 1
    if (canCover(mid)) {
      r = mid - 1
    } else {
      l = mid + 1
    }
  }

  return l
}

console.log(findRadius([1, 2, 3, 4, 5, 6, 7, 8, 9, 10], [2, 9]))

export {}
