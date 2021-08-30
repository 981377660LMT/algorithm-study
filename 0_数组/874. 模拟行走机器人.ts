/**
 * @param {number[]} commands
 * @param {number[][]} obstacles
 * @return {number}
 * 题目只有两种方式改变朝向，一种是左转（-2），另一种是右转（-1）。
 * 求的是机器人在运动过程中距离原点的最大值
 */
var robotSim = function (commands: number[], obstacles: number[][]): number {
  const ob = new Set<string>()
  obstacles.forEach(([x, y]) => ob.add(`${x}#${y}`))
  let x = 0
  let y = 0
  let max = 0
  let direction = 0
  /* 0 = north
   * 1 = east
   * 2 = south
   * 3 = west
   */
  for (let i = 0; i < commands.length; i++) {
    if (commands[i] === -1) {
      direction = (direction + 1) % 4 // updated
    } else if (commands[i] == -2) {
      // 注意这里的加4
      direction = (direction - 1 + 4) % 4
    } else {
      // 这样写更简洁
      while (commands[i]--) {
        // commands[i]--
        let previousX = x
        let previousY = y
        if (direction === 0) y++
        if (direction === 1) x++
        if (direction === 2) y--
        if (direction === 3) x--
        if (ob.has(`${x}#${y}`)) {
          {
            // 回退
            ;[x, y] = [previousX, previousY]
            break
          }
        }
      }
    }

    max = Math.max(max, x ** 2 + y ** 2)
  }

  return max
}

console.log(robotSim([4, -1, 4, -2, 4], [[2, 4]]))

export {}
