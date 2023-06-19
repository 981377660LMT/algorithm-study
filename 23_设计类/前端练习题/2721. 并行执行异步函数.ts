// 2721. 并行执行异步函数
// 手写promiseAll 函数，实现并行执行异步函数的功能。
// https://leetcode.cn/problems/execute-asynchronous-functions-in-parallel/
async function promiseAll<T>(functions: (() => Promise<T>)[]): Promise<T[]> {
  return new Promise((resolve, reject) => {
    const res: T[] = Array(functions.length)
    let ok = 0
    for (let i = 0; i < functions.length; i++) {
      const task = functions[i]
      task()
        // eslint-disable-next-line no-loop-func
        .then(data => {
          res[i] = data
          ok++
          if (ok === functions.length) {
            resolve(res)
          }
        })
        .catch(reject)
    }
  })
}

/**
 * const promise = promiseAll([() => new Promise(res => res(42))])
 * promise.then(console.log); // [42]
 */
