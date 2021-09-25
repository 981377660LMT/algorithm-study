import { ArrayDeque } from '../2_queue/Deque'

type Direction = 'L' | 'R' | 'U' | 'D'

class SnakeGame {
  private width: number
  private height: number
  private foodId: number
  private score: number
  private food: number[][]
  private snack: ArrayDeque<number>
  private visited: Set<number>

  constructor(width: number, height: number, food: number[][]) {
    this.width = width
    this.height = height
    this.food = food
    this.foodId = 0
    this.score = 0
    this.snack = new ArrayDeque(10 ** 8)
    this.visited = new Set()
    this.snack.push(0)
    this.visited.add(0)
  }

  move(direction: Direction): number {
    const head = this.snack.rear()!
    let row = ~~(head / this.width)
    let col = head % this.width

    if (direction === 'D') row++
    else if (direction === 'U') row--
    else if (direction == 'L') col--
    else col++

    const position = row * this.width + col
    // 1.第一种情况，是否超出边界
    if (row < 0 || row >= this.height || col < 0 || col >= this.width) return -1

    // 2. 第二种情况，吃到食物加头
    if (
      this.foodId < this.food.length &&
      row === this.food[this.foodId][0] &&
      col === this.food[this.foodId][1]
    ) {
      this.visited.add(position)
      this.snack.push(position)
      this.foodId++
      return ++this.score
    }

    // 3. 去尾
    this.visited.delete(this.snack.shift()!)

    // 4. 检查是否与自身相撞
    if (this.visited.has(position)) return -1
    else {
      this.visited.add(position)
      this.snack.push(position)
      return this.score
    }
  }
}
// 蛇在左上角的 (0, 0) 位置，身体长度为 1 个单位。
// 用一个Deque维护蛇身体，用一个哈希Set快速判断是否撞上身子

// 每次蛇的移动分为三种情况：

// 蛇超出边界，撞墙了，游戏结束
// 蛇吃到食物，增加头
// 蛇没吃到食物，去尾加头，这里去尾后如果蛇与自身相撞，游戏结束
