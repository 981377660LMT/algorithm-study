// Promise.race(iterable) 方法返回一个 promise，
// 一旦迭代器中的某个promise解决或拒绝，
// 返回的 promise就会解决或拒绝。

function race(arr: any[]): Promise<any> {
  if (arr.length === 0) return Promise.resolve()

  // 1
  const promises = arr.map(item => (item instanceof Promise ? item : Promise.resolve(item)))

  // 2
  return new Promise((resolve, reject) => {
    // 3   注意 for of 是并发
    for (const [_, promise] of promises.entries()) {
      promise.then(resolve).catch(reject)
    }
  })
}
