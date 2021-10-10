import { PromiseFunc } from './typings'

/**
 * @param {() => Promise<any>} func
 * @param {number} max
 * @return {Promise}
 * 假设你需要调用100个API获取数据，并且需要越快越好。
   如果使用Promise.all()，100个请求会同时到达你的服务器，如果你的服务器性能很低的话，这就会是个负担。
 * 请 节流API请求，使得任何时刻最多只有5个请求正在进行中。
   @description
   使用iter的精妙解法
 */
function throttlePromises(funcs: PromiseFunc[], max: number): Promise<unknown> {
  const results: unknown[] = []
  const iter = funcs.entries()

  // 注意他们共享一个iter 而iter是会被耗尽的 所以他们执行的是不同的promise
  // 即：使用fill来共用iter的引用
  const workers = Array(max).fill(iter).map(work)

  return Promise.all(workers).then(() => results)

  async function work(entries: IterableIterator<[number, PromiseFunc]>) {
    for (const [index, promiseFunc] of entries) {
      const result = await promiseFunc(1)
      results[index] = result
    }
  }
}

if (require.main === module) {
  const callApi = () =>
    new Promise(resolve => {
      setTimeout(() => {
        console.log(666)
        resolve('data')
      }, 1000)
    })

  throttlePromises([callApi, callApi, callApi], 2)
    .then(data => console.log(data))
    .catch(err => console.error(err))
}
