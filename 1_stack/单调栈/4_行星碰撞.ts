// 一个一个值添加到栈里面 然后控制栈里面的元素
// 带方向的使用栈
const asteroidCollison = (asteroids: number[]) => {
  let left = 0
  let right = 1
  const stack: number[] = []

  for (let i = 0; i < asteroids.length; i++) {
    const asteroid = asteroids[i]
    stack.push(asteroid)

    while (stack.length >= 2 && stack[stack.length - 2] > 0 && stack[stack.length - 1] < 0) {
      const bottomValue = stack[stack.length - 2]
      const topValue = Math.abs(stack[stack.length - 1])
      if (bottomValue > topValue) {
        stack.pop()
      } else if (bottomValue < topValue) {
        const top = stack.pop()
        stack.pop()
        stack.push(top!)
      } else {
        stack.pop()
        stack.pop()
      }
    }
  }

  return stack
}

console.log(asteroidCollison([5, 10, -5]))
console.log(asteroidCollison([8, -8]))
export {}
