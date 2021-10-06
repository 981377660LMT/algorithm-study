type Excutor = (resolve: (value: unknown) => void, reject: (reason?: any) => void) => void
type Function = (...args: any[]) => any

class MyPromise {
  private executor: Excutor
  private fullfilled: boolean
  private rejected: boolean
  private pending: boolean
  private handlers: Function[]
  private errorHandlers: Function[]

  constructor(executor: Excutor) {
    this.executor = executor
    this.fullfilled = false
    this.rejected = false
    this.pending = true
    this.handlers = []
    this.errorHandlers = []
    this.executor(this.resolve.bind(this), this.reject.bind(this))
  }

  static race(promises: MyPromise[]) {
    return new MyPromise((resolve, reject) => {
      for (const promise of promises) {
        // promise数组只要有任何一个promise 状态变更  就可以返回
        Promise.resolve(promise).then(resolve, reject)
      }
    })
  }

  static all(promises: MyPromise[]) {
    return new MyPromise((resolve, reject) => {
      let count = 0
      const allData: unknown[] = []
      for (const [i, promise] of promises.entries()) {
        //这里用 MyPromise.resolve包装一下 防止不是Promise类型传进来
        Promise.resolve(promise).then(res => {
          allData[i] = res
          count++
          if (count === promises.length) resolve(allData)
        }, reject)
      }
    })
  }

  then(onfulfilled?: Function, onrejected?: Function) {
    onfulfilled && this.handlers.push(onfulfilled)
    onrejected && this.errorHandlers.push(onrejected)
    return this
  }

  catch(onrejected?: Function) {
    onrejected && this.errorHandlers.push(onrejected)
    return this
  }

  private resolve(...args: any[]) {
    this.handlers.forEach(handler => handler(...args))
  }

  private reject(...args: any[]) {
    this.errorHandlers.forEach(handler => handler(...args))
  }
}

export {}
// test
const p1 = new MyPromise(resolve => setTimeout(() => resolve('ok'), 3000))
p1.then(console.log).then((...args) => console.log('second', ...args))

const p2 = new MyPromise((resolve, reject) => setTimeout(() => reject('rejected'), 3000))
p2.then(console.log).catch((...args) => console.log('fail', ...args))
