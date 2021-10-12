import { Observable } from './57. 实现Observable'

/**
 * @param {Array | ArrayLike | Promise | Iterable | Observable} input
 * @return {Observable}
 * 注意
   1.Observable是现成的，可以直接使用，不用再实现一遍。
   2.本问题中Observable-like意味着 Observable instance，
   虽然现实世界中你需要检查Symbol.observable，但是本题目中你可以简单处理。
 */
function from(
  input: Array<any> | ArrayLike<any> | Promise<any> | Iterable<any> | Observable
): Observable {
  if (Array.isArray(input) || 'length' in input) {
    return new Observable(sub => {
      for (let i = 0; i < input.length; i++) {
        sub.next(input[i])
      }
      sub.complete()
    })
  }

  if (Object.prototype.toString.call(input) === '[object Promise]') {
    return new Observable(async sub => {
      try {
        const res = await input
        sub.next(res)
        sub.complete()
      } catch (e: any) {
        sub.error(e)
      }
    })
  }

  // @ts-ignore
  // isIterable
  if (input != null && typeof input[Symbol.iterator] === 'function') {
    return new Observable(sub => {
      try {
        // @ts-ignore
        for (const item of input) {
          sub.next(item)
        }
        sub.complete()
      } catch (e: any) {
        sub.error(e)
      }
    })
  }

  if (input instanceof Observable) return input

  throw new Error('invalid input')
}

if (require.main === module) {
  from([1, 2, 3]).subscribe(console.log)
}

// from 函数 构造obersvable对象 ，根据参数类型决定 setup 种类

export { from }
