// Promise.allSettled()方法返回一个在所有给定的promise都已经fulfilled或rejected后的promise，
// 并带有一个对象数组，每个对象表示对应的promise结果。
// 和Promise.all()不同，Promise.allSettled() 会等待所有的promise直到fulfill或者reject。

interface AllSettledReturn {
  status: 'fulfilled' | 'rejected'
  value?: any
  reason?: any
}

/**
 * @param {Array<any>} promises - notice that input might contains non-promises
 * @return {Promise<Array<{status: 'fulfilled', value: any} | {status: 'rejected', reason: any}>>}
 */
function allSettled(arr: any[]): Promise<AllSettledReturn[]> {
  if (arr.length === 0) return Promise.resolve([])

  // 1.
  const promises = arr.map(item => (item instanceof Promise ? item : Promise.resolve(item)))

  // 2.
  return new Promise(resolve => {
    const res: AllSettledReturn[] = []
    let settledCount = 0

    // 3.
    for (const [index, promise] of promises.entries()) {
      promise
        .then(data => {
          res[index] = { status: 'fulfilled', value: data }
          settledCount++
          if (settledCount === promises.length) resolve(res)
        })
        .catch(err => {
          res[index] = { status: 'rejected', reason: err }
          settledCount++
          if (settledCount === promises.length) resolve(res)
        })
    }
  })
}

if (require.main === module) {
  allSettled([1, 2, 3, Promise.reject('err')]).then(console.log)
}
