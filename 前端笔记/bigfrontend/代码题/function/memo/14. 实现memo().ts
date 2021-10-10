// Memoization 是应用广泛的性能优化的手段，如果你开发过React应用，你一定不会对React.memo感到陌生。
// Memoization在算法题目中也经常用到，如果你可以用递归解决某个问题，那么很多时候加上Memoization可以得到更好的解法，甚至最终引导到动态规划的解法。
// 那么，请实现你自己的memo() 函数。传入相同的参数的时候，直接返回上一次的结果而不经过计算。
// 这是一种空间换时间的优化，在实际面试中，请仔细分析时间空间复杂度。

import { Func } from '../../typings'

/**
 * @param {Function} func
 * @param {(args:[]) => string }  [resolver] - cache key generator 决定缓存key的参数
 */
function memo(func: Func, resolver = (...args: any[]) => args.join('_')) {
  // your code here
  const cache = new Map<string, any>()

  return function (this: any, ...args: any[]) {
    const cacheKey = resolver(...args)
    if (cache.has(cacheKey)) return cache.get(cacheKey)
    const value = func.apply(this, args)
    cache.set(cacheKey, value)
    return value
  }
}

if (require.main === module) {
  const func = (arg1: number, arg2: number) => {
    return arg1 + arg2
  }
  const memoed = memo(func)
  memoed(1, 2)
  // 3， func 被调用
  memoed(1, 2)
  // 3，func 未被调用
  memoed(1, 3)
  // 4，新参数，func 被调用

  const memoed2 = memo(func, () => 'samekey')
  memoed2(1, 2)
  // 3，func被调用，缓存key是 'samekey'
  memoed2(1, 2)
  // 3，因为key是一样的，3被直接返回，func未调用
  memoed2(1, 3)
  // 3，因为key是一样的，3被直接返回，func未调用
}
