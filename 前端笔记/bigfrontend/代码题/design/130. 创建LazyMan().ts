interface Laziness {
  sleep: (time: number) => Laziness
  sleepFirst: (time: number) => Laziness
  eat: (food: string) => Laziness
}

// 1. Promise.resolve() 和new Promise都会立即执行
// 2.如果需要等待同步都执行完，需要采用微/宏任务
// 3.task的签名是一个返回Promise的函数
class _LazyMan implements Laziness {
  private task: ((...args: any[]) => Promise<any>)[]

  constructor(private name: string, private logFn: (log: string) => void) {
    this.task = [() => Promise.resolve(logFn(`Hi,I am ${name}`))]

    // 注意这里 同步执行完才调用
    setTimeout(() => {
      this.init()
    }, 0)
  }

  sleep(time: number): Laziness {
    const fn = () =>
      new Promise<void>(resolve =>
        setTimeout(() => {
          resolve()
        }, time * 1000)
      )
    this.task.push(fn)
    return this
  }

  sleepFirst(time: number): Laziness {
    const fn = () =>
      new Promise<void>(resolve => {
        setTimeout(() => {
          resolve()
        }, time * 1000)
      })
    this.task.unshift(fn)
    return this
  }

  eat(food: string): Laziness {
    const fn = () => Promise.resolve(this.logFn(`Eat ${food}`))
    this.task.push(fn)
    return this
  }

  private async init() {
    for (const t of this.task) {
      await t()
    }
  }
}

/**
 * @param {string} name
 * @param {(log: string) => void} logFn
 * @returns {Laziness}
 */
function LazyMan(name: string, logFn: (log: string) => void): Laziness {
  return new _LazyMan(name, logFn)
}
// LazyMan('Jack', console.log).eat('banana').eat('apple')
// LazyMan('Jack', console.log).eat('banana').sleep(10).eat('apple').sleep(1)
LazyMan('Jack', console.log).eat('banana').sleepFirst(10).eat('apple').sleep(1)
