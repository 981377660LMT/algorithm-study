import { promisify } from './promisify与callbackify'
import type { AsyncFunc, Callback } from './typings'

/**
 * @param {AsyncFunc[]} funcs
 * @return {(callback: Callback) => void}
 * 请实现一个async helper - parallel()。parallel() 有点类似Promise.all()
 * 能否使用Promise完成题目？能否不使用Promise完成该题目？
 */
function race(funcs: AsyncFunc[]): AsyncFunc {
  return (arg, callback) => {
    Promise.race(funcs.map(func => promisify(func)(arg)))
      .then(data => callback(null, data))
      .catch(err => callback(err, null))
  }
}

if (require.main === module) {
  const async1 = (arg: any, callback: Callback) => {
    setTimeout(() => callback(null, 1), 300)
  }

  const async2 = (arg: any, callback: Callback) => {
    setTimeout(() => callback(null, 2), 100)
  }

  const async3 = (arg: any, callback: Callback) => {
    setTimeout(() => callback(null, 3), 200)
  }

  const first = race([async1, async2, async3])

  first(null, (error, data) => {
    console.log(data) // 2, 因为2是第一个成功执行的结果
  })
}
