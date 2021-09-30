// 灯的颜色要想 变成蓝色 就必须同时满足下面两个条件
// 灯处于打开状态。
// 排在它之前（左侧）的所有灯也都处于打开状态。
// 请返回能够让 所有开着的 灯都 变成蓝色 的时刻 数目 。

// 如果当前所有灯都是蓝色，那么亮灯的编号和必定等于1+2+...+亮灯数目，维护亮灯数目和亮灯编号和即可
function numTimesAllBlue(light: number[]): number {
  let res = 0
  let sum = 0
  for (let i = 0; i < light.length; i++) {
    sum += light[i]
    if (sum === ((i + 1) * (i + 2)) / 2) res++
  }
  return res
}
console.log(numTimesAllBlue([2, 1, 3, 5, 4]))
// 解释：所有开着的灯都变蓝的时刻分别是 1，2 和 4 。
