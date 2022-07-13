// 一个一个值添加到栈里面 然后控制栈里面的元素
// 带方向的使用栈

function asteroidCollision(asteroids: number[]): number[] {
  const stack: number[] = []

  for (const asteroid of asteroids) {
    stack.push(asteroid)

    // 会碰撞的条件
    while (stack.length >= 2 && stack.at(-2)! > 0 && stack.at(-1)! < 0) {
      const pre = stack.at(-2)!
      const cur = -stack.at(-1)!
      if (pre > cur) {
        stack.pop()
      } else if (pre < cur) {
        const top = stack.pop()!
        stack.pop()
        stack.push(top)
      } else {
        stack.pop()
        stack.pop()
      }
    }
  }

  return stack
}

console.log(asteroidCollision([5, 10, -5]))
console.log(asteroidCollision([8, -8]))
export {}
