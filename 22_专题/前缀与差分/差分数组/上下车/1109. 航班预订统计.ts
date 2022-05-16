/**
 * @param {number[][]} bookings 1 <= bookings.length <= 20000 1 <= bookings[i][0] <= bookings[i][1] <= n <= 20000
 * @param {number} n 这里有  n  个航班，它们分别从 1 到 n 进行编号。
 * @return {number[]}
 * 预订记录 bookings[i] = [firsti, lasti, seatsi] 意味着在从 firsti 到 lasti （包含 firsti 和 lasti ）的 每个航班 上预订了 seatsi 个座位。
 * @summary
 *  由题意 遍历O(n^2)的方式不可取
 *  [i, j, k] 其实代表的是 第 i 站上来了 k 个人， 一直到 第 j 站都在飞机上，到第 j + 1 就不在飞机上了。
 * 想先一下坐高铁：第i站上来人，第j+1站下去人，问每个站车上多少人 =>前缀和
 */
const corpFlightBookings = function (bookings: number[][], n: number): number[] {
  const res = Array<number>(n + 1).fill(0)
  for (const [up, down, count] of bookings) {
    res[up - 1] += count
    res[down] += -count
  }
  res.reduce((pre, _, index, arr) => (arr[index] += pre))
  return res.slice(0, -1)
}

console.log(
  corpFlightBookings(
    [
      [1, 2, 10],
      [2, 3, 20],
      [2, 5, 25],
    ],
    5
  )
)
// 输出：[10,55,45,25,25]
export default 1
