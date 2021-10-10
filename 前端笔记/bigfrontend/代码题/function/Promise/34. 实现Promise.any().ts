// Promise.any() 接收一个Promise可迭代对象，
// 只要其中的一个 promise 成功，就返回那个已经成功的 promise 。
// 如果可迭代对象中没有一个 promise 成功（即所有的 promises 都失败/拒绝），
// 就返回一个失败的 promise 和AggregateError类型的实例，
// 它是 Error 的一个子类，用于把单一的错误集合在一起。本质上，这个方法和Promise.all()是相反的。

// AggregateError 暂时还没有被Chrome支持。
// 但是你仍然可以使用它因为我们在judge你的code时候添加了AggregateError。
// console.log(new AggregateError([Error('1'), Error('2')], 'msg'))

/**
 * @param {any[]} promises
 * @return {Promise}
 */
function any(arr: any[]): Promise<any> {
  // your code here
  if (arr.length === 0) return Promise.resolve([])

  const promises = arr.map(item => (item instanceof Promise ? item : Promise.resolve(item)))

  return new Promise((resolve, reject) => {
    const errors: Error[] = []
    let errorCount = 0

    for (const [index, promise] of promises.entries()) {
      promise.then(resolve).catch(err => {
        errors[index] = err
        errorCount++
        if (errorCount === promises.length) {
          reject(new AggregateError(errors, 'fail'))
        }
      })
    }
  })
}
