let j = 0

const a = () => {
  for (const i of [1, 2, 3, 4, 5]) {
    if (i === 3) return
    j += i
  }
}
a()

console.log(j)

// 注意 return 必须用于函数
// 注意return 只是中止当前调用栈中函数 并不能中止所有调用栈上函数
