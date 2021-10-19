chainAsync([
  (next: () => void) => {
    console.log('0 seconds')
    setTimeout(next, 1000)
  },
  (next: () => void) => {
    console.log('1 second')
    setTimeout(next, 1000)
  },
  () => {
    console.log('2 second')
  },
])

// 0 seconds
// 1 second
// 2 second
// 循环遍历包含异步事件的函数数组，并在每个异步事件完成时调用 next。
function chainAsync(funcs: Function[]) {
  let index = 0
  const last = funcs[funcs.length - 1]

  const next = () => {
    const fn = funcs[index++]
    fn === last ? fn() : fn(next)
  }

  next()
}
