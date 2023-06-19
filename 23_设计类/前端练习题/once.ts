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
