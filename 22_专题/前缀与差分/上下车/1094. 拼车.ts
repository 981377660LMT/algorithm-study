/**
 * @param {number[][]} trips 必须接送的乘客数量；乘客的上车地点；以及乘客的下车地点。
 * @param {number} capacity 车上最初有 capacity 个空座位可以用来载客
 * 判断你的车是否可以顺利完成接送所有乘客的任务
 */
const carPooling = function (trips: number[][], capacity: number): boolean {
  const res = Array<number>(1001).fill(0)
  for (const [count, up, down] of trips) {
    res[up] += count
    res[down] -= count
  }
  res.reduce((pre, _, index, arr) => (arr[index] += pre))
  return res.every(c => c <= capacity)
}

console.log(
  carPooling(
    [
      [2, 1, 5],
      [3, 3, 7],
    ],
    4
  )
)
console.log(
  carPooling(
    [
      [2, 1, 5],
      [3, 5, 7],
    ],
    3
  )
)
// true
console.log(
  carPooling(
    [
      [9, 0, 1],
      [3, 3, 7],
    ],
    4
  )
)
// false
export default 1
