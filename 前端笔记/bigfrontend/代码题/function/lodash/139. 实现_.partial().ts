// _.partial()类似于Function.prototype.bind() 但是不会固定this 。

/**
 * @param {Function} func
 * @param {any[]} args
 * @returns {Function}
 */
function partial(func: Function, ...args: any[]): Function {
  // your code here
  return function (this: any, ...restArgs: any[]) {
    const particalArgs = args.map(arg => (arg === partial.placeholder ? restArgs.shift() : arg)) // 遇到占位符直接从restArgs里取
    return func.call(this, ...particalArgs, ...restArgs)
  }
}

partial.placeholder = Symbol()

if (require.main === module) {
  const func = (...args: any[]) => args
  const _ = partial.placeholder
  const func1_3 = partial(func, 1, _, 3)

  console.log(func1_3(2, 4))
  // [1,2,3,4]
}
