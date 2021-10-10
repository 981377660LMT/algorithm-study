// 在题目14. 实现memo()中，你实现了一个memo函数但并不考虑空间成本。
// 现实项目中，cache如果无限量使用的话可能会导致内存不足，所以最好加上一些限制。
// 比如memoize-one 做的就是，仅缓存上一次的结果。

import { Func } from '../../typings'

function defaultIsEqual(args: any[], newArgs: any[]): boolean {
  if (args.length !== newArgs.length) return false
  return args.every((item, index) => item === newArgs[index])
}

/**
 * @param {Function} func
 * @param {(args: any[], newArgs: any[]) => boolean} [isEqual]  判断当前和上次的调用参数是否equal的函数
 * 默认的equal 判断函数的话，用===直接对于数组元素进行比较即可
 * @returns {any}
 */
function memoizeOne(func: Func, isEqual = defaultIsEqual): Func {
  let preArgs: any[] = []
  let preThis: any = null
  let preRes: any = null
  let isCalled = false

  return function (this: any, ...args: any[]) {
    if (isCalled && preThis === this && isEqual(preArgs, args)) return preRes
    preArgs = args
    preThis = this
    preRes = func.call(this, ...args)
    isCalled = true

    return preRes
  }
}
