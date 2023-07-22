// 手写 promisify 函数，将 callback 类型的函数转换为 Promise 类型的函数

type CallbackFn<V, R> = (next: (data: R, error?: string) => void, ...args: V[]) => void
type Promisified<V, R> = (...args: V[]) => Promise<R>

function promisify<V, R>(fn: CallbackFn<V, R>): Promisified<V, R> {
  return async function (...args) {
    return new Promise<R>((resolve, reject) => {
      fn((data, error) => {
        if (error !== void 0) reject(error)
        else resolve(data)
      }, ...args)
    })
  }
}

/**
 * const asyncFunc = promisify(callback => callback(42));
 * asyncFunc().then(console.log); // 42
 */
export {}
