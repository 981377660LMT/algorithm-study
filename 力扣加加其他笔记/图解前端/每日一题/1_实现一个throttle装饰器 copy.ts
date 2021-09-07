const throttle =
  (interval: number): MethodDecorator =>
  (target: Object, propertyKey: string | symbol, descriptor: PropertyDescriptor) => {
    const method = descriptor.value
    let timer: NodeJS.Timer | null = null

    descriptor.value = (...args: any[]) => {
      if (!timer) {
        timer = setTimeout(() => {
          timer = null
          method(...args)
        }, interval)
      }
    }

    return descriptor
  }

class Test {
  @throttle(1000)
  sayHi() {
    console.log('hi')
  }
}

const test = new Test()
test.sayHi()
test.sayHi()
test.sayHi()

export {}

function funcThrottle(fn: (...args: any[]) => any, interval: number) {
  let timer: NodeJS.Timer | null = null

  return function (this: any, ...args: any[]) {
    if (!timer)
      timer = setTimeout(() => {
        timer = null
        fn.call(this, ...args)
      }, interval)
  }
}

const testFunc = () => console.log('test')
const debouncedTestFunc = funcThrottle(testFunc, 1000)
debouncedTestFunc()
