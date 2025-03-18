// 重试指的是当加载出错时，有能力重新发起加载组件的请求。在加载组件的过程中，发生错误的情况非常常见，尤其是在网络不稳定
// 的情况下。因此，提供开箱即用的重试机制，会提升用户的开发体验。

// !当捕获到错误时，我们有两种选择：要么抛出错误，要么返回一个新的Promise 实例，
// 并把该实例的 resolve 和 reject 方法暴露给用户，让用户来决定下一步应该怎么做
// !将新的 Promise 实例的 resolve 和 reject 分别封装为 retry 函数和 fail 函数，并
// 将它们作为 onError 回调函数的参数。这样，用户就可以在错误发生时主动选择重试或直接抛出错误。

let retries = 0

function withRetry<T extends Promise<any>>(
  fn: () => T,
  onError?: (retry: () => void, fail: () => void, retries: number) => void
): T {
  return fn().catch(err => {
    if (!onError) {
      throw err
    }

    // 如果用户指定了 onError 回调，则将控制权交给用户
    return new Promise((resolve, reject) => {
      const retry = () => {
        resolve(withRetry(fn, onError))
        retries++
      }
      const fail = () => reject(err)
      onError(retry, fail, retries) // 作为 onError 回调函数的参数，让用户来决定下一步怎么做
    })
  }) as T
}

export {}

if (require.main === module) {
  const f = () =>
    new Promise<string>((resolve, reject) => {
      if (Math.random() > 0.8) {
        resolve('success')
      } else {
        reject('fail')
      }
    })

  withRetry(f, (retry, fail, retries) => {
    console.log('重试次数：', retries)
    retries < 30 ? retry() : fail()
  }).then(v => {
    console.log('最终结果：', v)
  })
}
