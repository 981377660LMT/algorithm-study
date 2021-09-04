class LazyMan {
  private name: string
  private tasks: ((...args: any[]) => Promise<unknown>)[]

  constructor(name: string) {
    this.name = name
    this.tasks = []
    // 同步执行完后才执行
    setTimeout(() => {
      this.trigger()
    }, 0)
  }

  eat(food: string) {
    const fn = () => Promise.resolve(console.log(`吃了${food}`))
    this.tasks.push(fn)
    return this
  }

  // sleep 和下面 sleepFirst 一样效果
  sleep(time: number) {
    const fn = () =>
      new Promise(resolve => setTimeout(resolve, time)).then(() => console.log(`睡了${time}`))
    this.tasks.push(fn)
    return this
  }

  sleepFirst(time: number) {
    const fn = () =>
      new Promise<void>(resolve =>
        setTimeout(() => {
          console.log(`睡了${time}`)
          resolve()
        }, time)
      )
    this.tasks.unshift(fn)
    return this
  }

  private async trigger() {
    this.tasks.reduce<Promise<unknown>>((pre, cur) => pre.then(cur), Promise.resolve())
    // for (const task of this.tasks) {
    //   await task()
    // }
  }
}

const lazyMan = (name: string) => new LazyMan(name)
// lazyMan('tony').eat('lunch').sleep(1000).eat('dinner')
lazyMan('tony').eat('lunch').eat('dinner').sleepFirst(4000).sleep(3000).eat('junkFood')

export default 1
