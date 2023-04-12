// 柯里化,收集参数的过程.
function curry(fn: Function): Function {
  const n = fn.length
  const args: unknown[] = []
  return function curried(...newArgs: unknown[]) {
    args.push(...newArgs)
    if (args.length < n) {
      return curried
    }
    return fn(...args)
  }
}

/**
 * function sum(a, b) { return a + b; }
 * const csum = curry(sum);
 * csum(1)(2) // 3
 */

export {}
