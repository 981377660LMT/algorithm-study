import type { AsyncFunc, PromiseFunc } from './typings'

// 化为fs.promises.readdir('path.txt')
function promisify(asyncFunc: AsyncFunc): PromiseFunc {
  return (arg: unknown) => {
    return new Promise((resolve, reject) => {
      asyncFunc(arg, (err, res) => {
        if (err) return reject(err)
        resolve(res)
      })
    })
  }
}

// 化为fs.readdir(path, options, callback)
function callbackify(promiseFunc: PromiseFunc): AsyncFunc {
  return (arg, callback) =>
    promiseFunc(arg)
      .then(data => callback(null, data))
      .catch(err => callback(err, null))
}

export { promisify, callbackify }
