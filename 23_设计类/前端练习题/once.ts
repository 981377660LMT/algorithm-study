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

/**
 * 方法装饰器 `@Once`，使被装饰的方法最多执行一次.
 */
export function Once<T extends (...args: any[]) => any>(
  _target: any,
  _propertyKey: string,
  descriptor: TypedPropertyDescriptor<T>
): TypedPropertyDescriptor<T> | void {
  const originalMethod = descriptor.value!
  let res: ReturnType<T>
  let hasRun = false

  const newMethod = function (this: any, ...args: Parameters<T>): ReturnType<T> {
    if (!hasRun) {
      hasRun = true
      res = originalMethod.apply(this, args)
    }
    return res
  }

  descriptor.value = newMethod as T
  return descriptor
}

if (require.main === module) {
  class Test {
    name = 'test'
    @Once
    test() {
      console.log(`test ${this.name}`)
    }

    @Once
    async add(a: number, b: number) {
      console.log(`add ${this.name}`)
      return a + b
    }
  }

  const test = new Test()
  test.test()
  test.test()

  test.add(1, 2).then(console.log)
  test.add(1, 2).then(console.log)
}
