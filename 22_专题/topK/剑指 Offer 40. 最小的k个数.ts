/**
 * @param {number[]} arr  0 <= arr[i] <= 10000  暗示可以计数排序
 * @param {number} k  0 <= k <= arr.length <= 10000
 * @return {number[]}
 *
 */
var getLeastNumbers = function (arr: number[], k: number): number[] {
  if (!arr.length || !k) return []
  const size = Math.max.apply(null, arr)
  const countArray = Array(size + 1).fill(0)
  arr.forEach(value => countArray[value]++)

  // 注意这种写法
  const res: number[] = Array(k).fill(0)
  let index = 0
  for (let num = 0; num < countArray.length; num++) {
    while (countArray[num]-- && index < k) {
      res[index++] = num
    }
    if (index >= k) break
  }

  return res
}

console.log(getLeastNumbers([3, 2, 1], 2))
