// eslint-disable-next-line @typescript-eslint/ban-types
export function once<T extends Function>(this: unknown, fn: T): T {
  // eslint-disable-next-line @typescript-eslint/no-this-alias
  const this_ = this
  let didCall = false
  let result: unknown

  // eslint-disable-next-line func-names
  return function () {
    if (didCall) {
      return result
    }

    didCall = true
    // eslint-disable-next-line prefer-rest-params
    result = fn.apply(this_, arguments)

    return result
  } as unknown as T
}

// TODO: 带条件的once
// 例如参数发生变化时，重新执行(mobx里的autorun设计)

/**
 * 只执行一次的函数.
 */
// eslint-disable-next-line @typescript-eslint/ban-types
function once2<T extends Function>(this: unknown, fn: T): T {
  let called = false
  let res: unknown

  // eslint-disable-next-line @typescript-eslint/ban-types
  const newFn: Function = (...args: unknown[]) => {
    if (called) return res
    called = true
    res = fn.apply(this, args)
    return res
  }

  return newFn as T
}
