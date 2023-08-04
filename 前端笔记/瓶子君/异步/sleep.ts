// ;(async () => {
//   await new Promise((res: (...args: any[]) => void) => setTimeout(res, 5000))
//   console.log(1)
// })()

// new 的时候 promise会立马执行 excutor函数
// 需要包裹在函数里才能不立马执行

//Promise 调用后立刻计时器睡眠
const sleep = (time: number) =>
  new Promise<void>(resolve => {
    console.log('开始睡眠')
    setTimeout(() => {
      resolve()
      console.log('结束睡眠')
    }, time)
  })

// Generator
function* sleepGenerator(time: number): Generator<Promise<void>> {
  yield new Promise<void>(resolve => setTimeout(resolve, time))
}
sleep(1000)
sleepGenerator(1000)
  .next()
  .value.then(() => console.log(1))

export {}
