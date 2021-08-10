/**
 * @param {number[]} piles
 * @param {number} h
 * @return {number}
 * @description 返回她可以在 H 小时内吃掉所有香蕉的最小速度 K（K 为整数）。
 * @summary 最左二分
 */
const minEatingSpeed = function (piles: number[], h: number): number {
  // 能力检测
  const calTime = (rate: number) =>
    piles.reduce((pre: number, cur: number) => pre + Math.ceil(cur / rate), 0)

  // 解空间就是 [1,max(piles)]。
  let l = 1
  let r = Math.max.apply(null, piles)

  while (l <= r) {
    const mid = Math.floor((l + r) / 2)
    const time = calTime(mid)
    if (time <= h) {
      r = mid - 1
    } else {
      l = mid + 1
    }
  }

  return l
}

console.log(minEatingSpeed([3, 6, 7, 11], 8))

export {}
