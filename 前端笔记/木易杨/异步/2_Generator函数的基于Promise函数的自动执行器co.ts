import { readFile } from './utils'
import type { Callback } from './utils'

// 因为剩余参数只能放最后 所以callback不得不写前面
// 虽然这与nodejs 的 回调形式的 api 不符 (
const promisify = <T>(
  func: (callback: Callback, ...args: T[]) => void
): ((...args: T[]) => Promise<any>) => {
  return function (this: any, ...args: T[]) {
    return new Promise((resolve, reject) => {
      func((err, data) => {
        if (err) return reject(err)
        resolve(data)
      }, ...args)
    })
  }
}

const promisifiedReadFile = promisify(readFile)

const gen = function* (): Generator<Promise<any>, any, any> {
  try {
    const file1 = yield promisifiedReadFile('a.txt')
    console.log(file1)
    const file2 = yield promisifiedReadFile('b.txt')
    console.log(file2)
    // ...
    const file3 = yield promisifiedReadFile('c.txt')
    console.log(file3)
  } catch (error) {
    console.log(error, 'oops')
  }
}

/**
 *
 * @param generator
 * @returns {Promise<any>} co 函数返回一个 Promise 对象，因此可以用 then 方法添加回调函数。
 * @description
 * 有点像 `flattenThunk` 函数里的 `inhancedCallback`
 */
const co = <T = any>(generator: () => Generator<Promise<T>, any, any>): Promise<any> => {
  const gen = generator()

  return new Promise((resolve, reject) => {
    onFulfilled(null)

    function onFulfilled(res: any) {
      let iteratorResult: IteratorResult<Promise<T>, any>

      try {
        iteratorResult = gen.next(res)
      } catch (error) {
        return reject(error)
      }

      next(iteratorResult)
    }

    function onReject(err: Error | null) {
      let iteratorResult: IteratorResult<Promise<T>, any>

      try {
        iteratorResult = gen.throw(err)
      } catch (error) {
        return reject(error)
      }

      next(iteratorResult)
    }

    function next(iteratorResult: IteratorResult<Promise<T>, any>) {
      if (iteratorResult.done) return resolve(iteratorResult.value)
      const promise = iteratorResult.value
      promise.then(onFulfilled).catch(onReject)
    }
  })
}

if (require.main === module) {
  co(gen).catch(console.log)
}

// 思路:
// run(gen)的思路
// 1.gen 是一个 ()=>Generator<Promise<any>, any, any> 的生成器函数
// 2.run 内部有一个 next(iteratorResult: IteratorResult<Promise<T>, any>): void
//   获取iteratorResult中的下一个promise 并执行 onFulfilled onReject

// 仔细看gen 函数 已经与 async await 非常相似了
// * 类比 async
// yield 类比 await
