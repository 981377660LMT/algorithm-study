/**
 * @param {Array<any>} promises - notice input might have non-Promises
 * @return {Promise<any[]>}
 * Promise.all(iterable) 方法返回一个 Promise 实例，
 * 此实例在 iterable 参数内所有的 promise 都“完成（resolved）”或参数中不包含 promise 时回调完成（resolve）；
 * 如果参数中 promise 有一个失败（rejected），此实例回调失败（reject），失败的原因是第一个失败 promise 的结果。
 */
async function all(promises: any[]): Promise<any[]> {
  // your code here\
  const results: any[] = []

  // 这样写的缺点：阻塞
  // 多个Promise是同步串行的，而不是异步并发的。在循环中使用 await 将使每个迭代停止，直到每个Promise被解析。
  for (const promise of promises) {
    results.push(await promise)
  }

  return results
}

function all2(arr: any[]): Promise<any[]> {
  if (arr.length === 0) return Promise.resolve([])

  // 1.
  const promises = arr.map(item => (item instanceof Promise ? item : Promise.resolve(item)))

  // 2.
  return new Promise((resolve, reject) => {
    const res: any[] = []
    let fulfilledCount = 0
    let isErrored = false

    // 3.
    for (const [index, promise] of promises.entries()) {
      promise.then(
        value => {
          if (isErrored) return
          res[index] = value
          fulfilledCount++
          if (fulfilledCount === promises.length) resolve(res)
        },
        reason => {
          if (isErrored) return
          isErrored = true
          reject(reason)
        }
      )
    }
  })
}

all2([
  Promise.resolve(1),
  Promise.reject('foo'),
  Promise.resolve('no'),
  1,
  Promise.resolve(2),
  null,
])
  .then(console.log)
  .catch(console.error)
