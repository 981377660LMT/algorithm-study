import { readFile } from '../utils'
import type { Callback, Thunk } from '../utils'

// 因为剩余参数只能放最后 所以callback不得不写前面
// 虽然这与nodejs 的 回调形式的 api 不符 (
const thunkify = <T>(
  func: (callback: Callback, ...args: T[]) => void
): ((...args: T[]) => Thunk) => {
  return function (this: any, ...args: T[]) {
    return (callback: Callback) => {
      func.call(this, callback, ...args)
    }
  }
}

const thunkifiedReadFile = thunkify(readFile)

const gen = function* (): Generator<Thunk, any, any> {
  try {
    const file1 = yield thunkifiedReadFile('a.txt')
    console.log(file1)
    const file2 = yield thunkifiedReadFile('b.txt')
    console.log(file2)
    // ...
    const file3 = yield thunkifiedReadFile('c.txt')
    console.log(file3)
  } catch (error) {
    console.log(error, 'oops')
  }
}

/**
 *
 * @param generator
 * @description
 * 有点像 `flattenThunk` 函数里的 `inhancedCallback`
 */
const run = (generator: () => Generator<Thunk, any, any>): void => {
  const gen = generator()

  function next(err: Error | null, data: any) {
    if (err) gen.throw(err)
    const iteratorResult = gen.next(data)
    if (iteratorResult.done) return
    iteratorResult.value(next)
  }

  next(null, null)
}

if (require.main === module) {
  run(gen)
}

// 思路:
// run(gen)的思路
// 1.gen 是一个 ()=>Generator<Thunk, any, any> 的生成器函数
// 2.run 内部有一个 next(err: Error | null, data: any): void
//   最关键的 next 函数，它会反复调用自身
//   这个next 是传给各个thunk的 内部递归
// 3.thunk 是一个接收callback的函数
// 4.thunkify 将异步函数全部改为 (...args: T[]) => Thunk的形式

// 仔细看gen 函数 已经与 async await 非常相似了
// * 类比 async
// yield 类比 await
