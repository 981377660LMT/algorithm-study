/* eslint-disable no-await-in-loop */
/* eslint-disable no-loop-func */
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

// https://leetcode.cn/problems/design-cancellable-function/solution/an-ti-mu-miao-shu-zhu-ju-wu-nao-shi-xian-d8fe/
function cancellable<T>(generator: Generator<Promise<any>, T, unknown>): [() => void, Promise<T>] {
  // !Promise.withResolvers() 提案
  let promise: Promise<T>
  let resolve: (value: T) => void
  let reject: (reason?: any) => void
  promise = new Promise<T>((resolve_, reject_) => {
    resolve = resolve_
    reject = reject_
  })

  const run = async () => {
    try {
      let p = generator.next()
      while (!p.done) {
        try {
          // [2]生成器函数只会生成 promise 对象
          const value = await p.value
          // [3]将 promise 对象解析的值传回生成器
          p = generator.next(value)
        } catch (error) {
          // [4]如果 promise 被拒绝...应将该错误抛回给生成器
          p = generator.throw(error)
        }
      }

      // [9]当生成器完成时...应该解析为生成器返回的值
      resolve(await p.value)
    } catch (error) {
      // [10]如果生成器抛出错误...应该拒绝并抛出该错误
      reject(error)
    }
  }

  run()

  const cancel = async (): Promise<void> => {
    // [5]调用了取消回调函数
    try {
      // [6]应该将错误抛回给生成器。该错误应该是字符串 "Cancelled"
      let p = generator.throw('Cancelled')
      // [7]如果错误被捕获...应该解析为下一个生成或返回的值
      resolve(await p.value)
    } catch (error) {
      // [8]否则...应该被拒绝并抛出该错误
      reject(error)
    }
  }

  return [cancel, promise]
}

export {}
