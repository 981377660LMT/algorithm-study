// _.once(func)可以用来缓存结果使得原函数至多被调用一次。
/**
 * @param {Function} func
 * @return {Function}
 */
function once(func: Function): Function {
  // your code here
  let preRes: unknown = null
  let isCalled = false

  return function (this: unknown, ...args: any[]) {
    if (isCalled) return preRes
    preRes = func.call(this, ...args)
    isCalled = true
    return preRes
  }
}

if (require.main === module) {
  function func<T>(num: T): T {
    return num
  }

  const onced = once(func)

  console.log(onced(1))
  // 1

  console.log(onced(2))
  // 1，因为已经调用过了，前一次的结果被直接返回
}
