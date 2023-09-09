// !promise 本质是函数调用，无法取消副作用，`promise 一旦创建，就不能取消，`

// !其实取消promise执行和取消请求是一样的，并不是真的终止了代码的执行，而是对结果不再处理

// 1. 利用Promise.race
function makeCancelablePromise1<T>(rawPromise: Promise<T>) {
  let cancel!: () => void
  const cancelPromise = new Promise<any>((_, reject) => {
    cancel = () => reject(new Error('promise is canceled'))
  })
  const wrappedPromise = Promise.race<Promise<T>[]>([cancelPromise, rawPromise])

  return {
    cancel,
    promise: wrappedPromise
  }
}

const p = makeCancelablePromise1(
  new Promise<string>(resolve => {
    setTimeout(() => {
      resolve('data')
    }, 1000)
  })
)
p.cancel()
p.promise.then(console.log).catch(console.error)

// 2. 利用标志位
type CancellablePromise<T> = Promise<T> & { cancel: () => void }
function useCancellablePromise<T>(p: Promise<T>): CancellablePromise<T> {
  let isCanceled = false

  const cancel = () => {
    isCanceled = true
  }

  const wrapped = new Promise<T>((resolve, reject) => {
    p.then(data => {
      isCanceled ? reject(new Error('promise is canceled')) : resolve(data)
    }).catch(err => {
      isCanceled ? reject(new Error('promise is canceled')) : reject(err)
    })
  })

  wrapped.cancel = cancel

  return wrapped
}

const cp = useCancellablePromise(
  new Promise<string>(resolve => {
    setTimeout(() => {
      resolve('data')
    }, 1000)
  })
)
cp.cancel()
cp.then(console.log).catch(console.error)

// !3.  AbortController 取消fetch
// AbortController接口表示一个控制器对象，允许你根据需要中止一个或多个 Web请求。s
// AbortController.abort()
// 中止一个尚未完成的Web(网络)请求。这能够中止fetch 请求，任何响应Body的消费者和流。
// let controller = new AbortController()

// let task = new Promise((resolve, reject) => {
//    some logic ...(请求后的处理)
//   controller.signal.addEventListener('abort', () => reject('oops'))
// })

// controller.abort() // task is now in rejected state
export {}
