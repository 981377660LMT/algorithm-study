// 前面我们已经得到了一个在前端限制调用频率的方案。
// 但是，即使我们已经在前端限制了调用频率，依然可能遇到错误：
//
// 前端的流控无法完全满足后端的流控限制。后端可能会对所有用户的调用之和做一个整体限制。
// 比如所有用户的调用频率不能超过每秒一万次，前端流控无法对齐这种限制。
// 非流控错误。比如后端服务或网络不稳定，造成的短暂不可用。
// !因此，面对这些前端不可避免的错误，需要通过重试来得到结果。
//
// 这里提供一个重试工具函数wrapRetry，它的好处是：
// 使用简单、对调用者透明：与前面的流控工具函数一样，只需要包装一下你原本的异步函数，即可得到自动重试的函数，它与原本的异步函数使用方式相同。
// !为了增加扩展性，将控制权全部交给调用者。

import Limit from './1_限流'

function withRetry<A extends readonly unknown[], R>(
  fn: (...args: A) => Promise<R>,
  onError: (retry: () => void, fail: () => void, retries: number) => void
): (...args: A) => Promise<R> {
  const innerRetry =
    (retryCount: number) =>
    (...args: A): Promise<R> =>
      fn(...args).catch(
        err =>
          new Promise((resolve, reject) => {
            const retry = () => resolve(innerRetry(retryCount + 1)(...args))
            const fail = () => reject(err)
            onError(retry, fail, retryCount)
          })
      )

  return innerRetry(0)
}

export default {
  withRetry
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const f = (a: number) =>
    new Promise<string>((resolve, reject) => {
      if (Math.random() > 0.9) {
        resolve(`success: ${a}`)
      } else {
        reject(new Error('fail'))
      }
    })

  let g = Limit.withLimit(f, 3)
  g = withRetry(g, (retry, fail, retries) => {
    console.log('重试次数：', retries)
    retries < 15 ? retry() : fail()
  })

  g(1).then(v => {
    console.log('最终结果：', v)
  })
}
