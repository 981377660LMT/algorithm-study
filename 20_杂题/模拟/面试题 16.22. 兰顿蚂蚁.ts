const enum Color {
  White = '_',
  Black = 'X',
}

const enum Direction {
  Up = 'U',
  Down = 'D',
  Right = 'R',
  Left = 'L',
}

type Position = string

/**
 * @param {number} K  K <= 100000
 * @return {string[]}
 * 一只蚂蚁坐在由白色和黑色方格构成的无限网格上。
 * 开始时，网格全白，蚂蚁面向右侧。
 * (1) 如果在白色方格上，则翻转方格的颜色，向右(顺时针)转 90 度，并向前移动一个单位。
   (2) 如果在黑色方格上，则翻转方格的颜色，向左(逆时针方向)转 90 度，并向前移动一个单位。
   编写程序来模拟蚂蚁执行的前 K 个动作，并返回最终的网格。
   状态模式

 */
const printKMoves = function (K: number): string[] {
  let [curX, curY] = [0, 0]
  let curDirection: Direction = Direction.Right
  const trace = new Map<Position, Color>()

  const moveUp = () => {
    curY--
    curDirection = Direction.Up
  }
  const moveDown = () => {
    curY++
    curDirection = Direction.Down
  }
  const moveRight = () => {
    curX++
    curDirection = Direction.Right
  }
  const moveLeft = () => {
    curX--
    curDirection = Direction.Left
  }
  const moveOnWhite = (position: string) => {
    // 改变当前格子的颜色再走
    trace.set(position, Color.Black)
    switch (curDirection) {
      case Direction.Left:
        moveUp()
        break
      case Direction.Up:
        moveRight()
        break
      case Direction.Right:
        moveDown()
        break
      case Direction.Down:
        moveLeft()
        break
      default:
        break
    }
  }
  const moveOnBlack = (position: string) => {
    trace.set(position, Color.White)
    switch (curDirection) {
      case Direction.Left:
        moveDown()
        break
      case Direction.Right:
        moveUp()
        break
      case Direction.Up:
        moveLeft()
        break
      case Direction.Down:
        moveRight()
        break
      default:
        break
    }
  }

  for (let i = 0; i < K; i++) {
    const position = `${curX}#${curY}`
    !trace.has(position) && trace.set(position, Color.White)
    if (trace.get(position) === Color.White) moveOnWhite(position)
    else if (trace.get(position) === Color.Black) moveOnBlack(position)
  }

  //  把走过的格子拿出来作为规范化的依据
  //  主要是求出最小的格子，注意可能是负数
  //  还要求出返回的矩形长宽
  const pos = [...trace.keys()].map(p => p.split('#'))
  const x = pos.map(p => parseInt(p[0]))
  const y = pos.map(p => parseInt(p[1]))
  x.push(curX)
  y.push(curY)
  console.log(trace, x, y)
  const xMax = Math.max.apply(null, x)
  const yMax = Math.max.apply(null, y)
  const xMin = Math.min.apply(null, x)
  const yMin = Math.min.apply(null, y)
  const xLen = xMax - xMin + 1
  const yLen = yMax - yMin + 1

  // 创建默认值均为白色的二维列表，因为没走过的格子都是白色
  // 注意行是y长度 列是x长度
  const res = Array.from<unknown, string[]>({ length: yLen }, () => Array(xLen).fill(Color.White))
  // 把字典中记录的格子颜色取出
  for (const [position, color] of trace.entries()) {
    const pos = position.split('#')
    const x = parseInt(pos[0])
    const y = parseInt(pos[1])
    res[y - yMin][x - xMin] = color
  }
  console.log(res, curX, curY, xMin, yMin)
  res[curY - yMin][curX - xMin] = curDirection

  // console.table(res)
  return res.map(row => row.join(''))
}

// console.log(printKMoves(2))
// 输出:
// [
//   "_X",
//   "LX"
// ]
// console.log(printKMoves(0))
// 输出: ["R"]
console.log(printKMoves(1))

export {}
