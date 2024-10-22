// function Lockify(target: any, propertyKey: string, descriptor: PropertyDescriptor) {
//   const originalMethod = descriptor.value

//   let isLocked = false

//   descriptor.value = async function (...args: any[]) {
//     if (isLocked) {
//       console.log(`方法 ${propertyKey} 正在被调用，此次调用将被忽略。`)
//       return
//     }

//     isLocked = true
//     try {
//       const result = await originalMethod.apply(this, args)
//       return result
//     } finally {
//       isLocked = false
//     }
//   }
// }

type AsyncFunction<T extends any[], R> = (...args: T) => Promise<R>

/**
 * `@Lockify` 装饰器.
 *
 * 保证同一时刻只有一个异步操作可以执行, 调用期间的其他调用将被忽略.
 */
const Lockify = <T extends any[]>(
  _target: Object,
  propertyKey: string | symbol,
  descriptor: TypedPropertyDescriptor<AsyncFunction<T, void>>
): void => {
  const originalMethod = descriptor.value!
  let isLocked = false

  descriptor.value = async function (...args: T) {
    if (isLocked) {
      console.log(`Method ${String(propertyKey)} is being called. This call will be ignored.`)
      return
    }

    isLocked = true
    try {
      await originalMethod.apply(this, args)
    } finally {
      isLocked = false
    }
  }
}

/**
 * 保证同一时刻只有一个异步操作可以执行, 调用期间的其他调用将被忽略.
 * 注意 tt-mrn 不支持装饰器.
 */
function lockify<T extends any[], R>(fn: AsyncFunction<T, R>): AsyncFunction<T, R | void> {
  let isLocked = false

  return async function (...args: T): Promise<R | void> {
    if (isLocked) {
      console.log(`Function is being called. This call will be ignored.`)
      return
    }
    isLocked = true
    try {
      return await fn(...args)
    } finally {
      isLocked = false
    }
  }
}

export {}

if (require.main === module) {
  class Test {
    @Lockify
    async effect() {
      await new Promise(resolve => setTimeout(resolve, 2000))
      console.log('Request result')
    }
  }

  const test = new Test()
  test.effect()
  test.effect()
  test.effect()

  const f = lockify(async () => {
    await new Promise(resolve => setTimeout(resolve, 2000))
    console.log('Request result')
  })

  f()
  f()
  f()
}
