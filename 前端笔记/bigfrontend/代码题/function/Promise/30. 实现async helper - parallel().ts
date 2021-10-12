import { promisify } from './promisify与callbackify'
import type { AsyncFunc, Callback, PromiseFunc } from './typings'

/**
 * @param {AsyncFunc[]} funcs
 * @return {(callback: Callback) => void}
 * 请实现一个async helper - race() 。 race()有点类似Promise.race()
 * 能否使用Promise完成题目？能否不使用Promise完成该题目？parallel()
 * race()会在任何一个function结束或者产生error的时候调用最终的callback。
 */
function parallel(funcs: AsyncFunc[]): AsyncFunc {
  return (arg, callback) => {
    Promise.all(funcs.map(func => promisify(func)(arg)))
      .then(data => callback(null, data))
      .catch(err => callback(err, null))
  }
}

if (require.main === module) {
  const async1 = (arg: any, callback: Callback) => {
    callback(null, 1)
  }

  const async2 = (arg: any, callback: Callback) => {
    callback(null, 2)
  }

  const async3 = (arg: any, callback: Callback) => {
    callback(null, 3)
  }

  const all = parallel([async1, async2, async3])

  all(null, (error, data) => {
    console.log(data) // [1, 2, 3]
  })
}
