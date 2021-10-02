type Time = number
// 所有车的初始位置各不相同
function carFleet(target: number, position: number[], speed: number[]): number {
  const stack: Time[] = [] // 消耗时间单减
  const cars = position.map((pos, index) => [pos, speed[index]]).sort((a, b) => a[0] - b[0])

  for (const [pos, spd] of cars) {
    const time = (target - pos) / spd
    while (stack.length && time >= stack[stack.length - 1]) {
      stack.pop()
    }
    stack.push(time)
  }
  return stack.length
}

console.log(carFleet(12, [10, 8, 0, 5, 3], [2, 4, 1, 1, 3]))
// 会有多少车队到达目的地?
// 输入：target = 12, position = [10,8,0,5,3], speed = [2,4,1,1,3]
// 输出：3
// 解释：
// 从 10 和 8 开始的车会组成一个车队，它们在 12 处相遇。
// 从 0 处开始的车无法追上其它车，所以它自己就是一个车队。
// 从 5 和 3 开始的车会组成一个车队，它们在 6 处相遇。
// 请注意，在到达目的地之前没有其它车会遇到这些车队，所以答案是 3。
