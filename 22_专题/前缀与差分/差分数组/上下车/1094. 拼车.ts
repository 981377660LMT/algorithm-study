/**
 * @param {number[][]} trips 必须接送的乘客数量；乘客的上车地点；以及乘客的下车地点。
 * @param {number} capacity 车上最初有 capacity 个空座位可以用来载客
 * 判断你的车是否可以顺利完成接送所有乘客的任务
 */
function carPooling(trips: number[][], capacity: number): boolean {
  const diff = new Int32Array(1001)
  for (const [count, up, down] of trips) {
    diff[up] += count
    diff[down] -= count
  }
  diff.reduce((pre, _, index, arr) => (arr[index] += pre))
  return diff.every(num => num <= capacity)
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
