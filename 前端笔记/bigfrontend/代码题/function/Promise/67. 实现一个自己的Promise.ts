import type { Executor, Reject, Resolve, Status } from './typings'

class MyPromise<T = any> {
  private resolve: Resolve<T>
  private reject: Reject
  private status: Status
  private resolveCallbacks: (() => void)[]
  private rejectCallbacks: (() => void)[]
  private resolvedValue?: T
  private rejectedReason?: any

  constructor(executor: Executor<T>) {
    this.status = 'pending'
    this.resolveCallbacks = []
    this.rejectCallbacks = []

    this.resolve = (value: T) => {
      if (this.status === 'pending') {
        this.status = 'onfulfilled'
        this.resolvedValue = value
        this.resolveCallbacks.forEach(cb => cb()) //消费， resolve 消费后面then里的cb
      }
    }

    this.reject = (reason?: any) => {
      if (this.status === 'pending') {
        this.status = 'onrejected'
        this.rejectedReason = reason
        this.rejectCallbacks.forEach(cb => cb()) //消费， reject 消费后面then里的cb
      }
    }

    // 执行中途出错
    try {
      executor(this.resolve, this.reject)
    } catch (error: any) {
      console.log('error', 789)
      this.status = 'pending'
      this.reject(error)
      throw new Error('程序停止...')
    }
  }

  /**
   *
   * @param onfulfilled
   * @param onrejected
   * @returns
   * then 里会判断当前promise的状态
   */
  then(onfulfilled?: Resolve<T>, onrejected?: Reject): MyPromise<T> {
    return new MyPromise<T>((resolve, reject) => {
      let result: any // 这一个promise的结果

      if (this.status === 'onfulfilled') {
        if (onfulfilled) {
          result = onfulfilled(this.resolvedValue!)
          resolve(result)
        }
      } else if (this.status === 'onrejected') {
        if (onrejected) {
          result = onrejected(this.rejectedReason)
          reject(result)
        }
      } else {
        this.processAsyncTask(onfulfilled, onrejected, resolve, reject)
      }
    })
  }

  catch(onrejected?: Reject): MyPromise<T> {
    return this.then(undefined, onrejected)
  }

  /**
   *
   * @param onfulfilled
   * @param onrejected
   * @param nextResolve
   * @param nextReject
   * 1.推入队列
   * 2.获取当前promise返回结果
   * 3.如果是promise 则result.then(nextResolve,nextReject) 否则 nextResolve(result)
   */
  private processAsyncTask(
    onfulfilled: Resolve<T> | undefined,
    onrejected: Reject | undefined,
    nextResolve: Resolve<T> | undefined,
    nextReject: Reject | undefined
  ) {
    let result: any // 当前promise的返回值

    this.resolveCallbacks.push(() => {
      if (onfulfilled) {
        result = onfulfilled(this.resolvedValue!)

        if (isPromise(result)) {
          result.then(nextResolve, nextReject)
        } else {
          nextResolve && nextResolve(result)
        }
      }
    })

    this.rejectCallbacks.push(() => {
      onrejected && (result = onrejected(this.rejectedReason))
      nextReject && nextReject(result)
    })
  }
}

if (require.main === module) {
  const promise = new MyPromise<string>((resolve, reject) => {
    setTimeout(() => {
      resolve('ok')
    }, 1)
  })
    .then(
      data => {
        return new MyPromise<string>(resolve =>
          setTimeout(() => {
            resolve(data + '2')
          }, 5)
        )
      },
      err => console.log(err)
    )
    .then(data => console.log(data))
}

// new Promise<void>((resolve, reject) => {
//   // resolve()
//   setTimeout(() => {
//     resolve()
//   }, 10000)
//   reject('no')
// })
//   // @ts-ignore
//   .then(console.log(100))
//   .catch(err => console.log(err))

// 1.防止resolve中途出错
// 2.同步级联then
// 3. 异步级联then：then 生产 resolve时消费

function isFunction(input: unknown): input is Function {
  return typeof input === 'function'
}

function isObject(input: unknown): input is Record<any, unknown> {
  return input !== null && typeof input === 'object'
}

function isPromise(input: unknown): input is MyPromise {
  return isObject(input) && isFunction(input.then)
}
