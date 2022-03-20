import { callbackify, promisify } from './promisify与callbackify'
import type { AsyncFunc, Callback, PromiseFunc } from './typings'

/**
 * @param {AsyncFunc[]} funcs
 * @return {(callback: Callback) => void}
 * 请实现一个async helper - sequence()。sequence()像pipe() 那样将异步函数串联在一起。
 * 能否使用Promise完成题目？能否不使用Promise完成该题目？
 */
function sequence(funcs: AsyncFunc[]): AsyncFunc {
  const promiseFuncs = funcs.map(promisify)

  const mergeTwo = (promiseFunc1: PromiseFunc, promiseFunc2: PromiseFunc): PromiseFunc => {
    return arg => promiseFunc1(arg).then(promiseFunc2)
  }

  return callbackify(promiseFuncs.reduce(mergeTwo))
}

if (require.main === module) {
  const asyncTimes2 = (num: number, callback: Callback) => {
    setTimeout(() => callback(null, num * 2), 100)
  }

  const asyncTimes4 = sequence([asyncTimes2, asyncTimes2])

  asyncTimes4(1, (error, res) => {
    console.log(res) // 4
  })
}
