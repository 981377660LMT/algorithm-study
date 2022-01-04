// 一个一个值添加到栈里面 然后控制栈里面的元素
// 带方向的使用栈
const asteroidCollision = (asteroids: number[]) => {
  const stack: number[] = []

  for (let i = 0; i < asteroids.length; i++) {
    const asteroid = asteroids[i]
    stack.push(asteroid)

    // 会碰撞的条件
    while (stack.length >= 2 && stack[stack.length - 2] > 0 && stack[stack.length - 1] < 0) {
      const pre = stack[stack.length - 2]
      const cur = Math.abs(stack[stack.length - 1])
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
