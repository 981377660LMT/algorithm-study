// 有时候你会有一个长时间运行的任务，并且你可能希望在它完成之前取消它。
// 为了实现这个目标，请你编写一个名为 cancellable 的函数，它接收一个生成器对象，
// 并返回一个包含两个值的数组：一个 取消函数 和一个 promise 对象。

// 你可以假设生成器函数只会生成 promise 对象。你的函数负责将 promise 对象解析的值传回生成器。
// 如果 promise 被拒绝，你的函数应将该错误抛回给生成器。

// !如果在生成器完成之前调用了取消回调函数，则你的函数应该将错误抛回给生成器。
// !该错误应该是字符串 "Cancelled"（而不是一个 Error 对象）。
// 如果错误被捕获，则返回的 promise 应该解析为下一个生成或返回的值。
// 否则，promise 应该被拒绝并抛出该错误。不应执行任何其他代码。

// 当生成器完成时，您的函数返回的 promise 应该解析为生成器返回的值。
// 但是，如果生成器抛出错误，则返回的 promise 应该拒绝并抛出该错误。

// TODO https://leetcode.cn/problems/design-cancellable-function/
function cancellable<T>(generator: Generator<Promise<any>, T, unknown>): [() => void, Promise<T>] {
  let isDone = false

  const cancel = () => {
    if (isDone) return
    generator.throw('Cancelled')
  }

  const promise = new Promise<T>((resolve, reject) => {
    const next = (args?: any) => {
      if (isDone) return
      try {
        const { value: nextTask, done } = generator.next(args)
        if (done) {
          isDone = true
          resolve(nextTask)
        } else {
          ;(nextTask as Promise<T>).then(next).catch(reject)
        }
      } catch (error) {
        isDone = true
        reject(error)
      }
    }
    next()
  })

  return [cancel, promise]
}

/**
 * function* tasks() {
 *   const val = yield new Promise(resolve => resolve(2 + 2));
 *   yield new Promise(resolve => setTimeout(resolve, 100));
 *   return val + 1;
 * }
 * const [cancel, promise] = cancellable(tasks());
 * setTimeout(cancel, 50);
 * promise.catch(console.log); // logs "Cancelled" at t=50ms
 */

export {}

// function*() { return 42; }
// {"cancelledAt":100}
if (require.main === module) {
  const g = function* () {
    yield new Promise(res => setTimeout(res, 200))
    return 'Success'
  }
  const [cancel, promise] = cancellable(g())
  setTimeout(cancel, 100)
  promise.then(console.log).catch(console.log)
}

// 提示1：
// This question tests understanding of two-way communication between generator functions
// and the code that evaluates the generator.
// !It is a powerful technique which is used in libraries such as `redux-saga`.
// 提示2：
// You can pass a value value to a generator function X by calling generator.next(X).
// Then in the generator function,
// you can access this value by calling let X = yield "val to pass into generator.next()";
// 提示3：
// You can throw an error back to a generator function by calling generator.throw(err).
// If this error isn't caught in the generator function, that will throw an error.
