function fib(n: number, a = 0, b = 1): number {
  if (n === 0) return a
  if (n === 1) return b
  return fib(n - 1, b, a + b)
}

console.log(fib(1000))
console.log(fib(1001))

// 尾递归为啥能优化？
// 完全等效于一个无栈的循环

// 尾递归就要把返回值放在函数参数里
