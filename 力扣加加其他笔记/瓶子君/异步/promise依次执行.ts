const delay = (time: number) => () =>
  new Promise<void>(resolve =>
    setTimeout(() => {
      resolve()
      console.log('延迟结束')
    }, time)
  )
const task2 = () => Promise.resolve(console.log('ok2'))
const task3 = () => Promise.resolve(console.log('ok3'))
const task4 = () => Promise.resolve(console.log('ok4'))

const tasks = [delay(1000), task2, task3, task4]

// 下面三种效果一样
const run = async () => {
  // for await (const task of tasks) {
  //   await task()
  // }
  // for (const task of tasks) {
  //   await task()
  // }
  tasks.reduce<Promise<unknown>>((pre, cur) => pre.then(cur), Promise.resolve())
}
run()

export {}
