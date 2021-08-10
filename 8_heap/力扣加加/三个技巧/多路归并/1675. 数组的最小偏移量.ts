import { smallestRange } from './632. 最小区间'

/**
 * @param {number[]} nums
 * @return {number}
 * @description 如果元素是 偶数 ，可以除以 2 如果元素是 奇数 ，可以乘上 2
 * 求任意两个元素间最大差值的最小值
 * 这题和632最小区间等价
 */
const minimumDeviation = function (nums: number[]): number {
  const select = Array.from<number, number[]>({ length: nums.length }, () => [])
  nums.forEach((num, index) => {
    if (num % 2 === 1) {
      select[index].push(num, num * 2)
    } else {
      let start = num
      for (; start % 2 === 0; start /= 2) {
        select[index].push(start)
      }
      select[index].push(start)
    }
    select[index].sort((a, b) => a - b)
  })
  const res = smallestRange(select)
  return res[1] - res[0]
}

console.log(minimumDeviation([1, 2, 3, 4]))
