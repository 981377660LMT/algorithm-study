const debounce =
  (delay: number): MethodDecorator =>
  (target: Object, propertyKey: string | symbol, descriptor: PropertyDescriptor) => {
    const method = descriptor.value
    let timer: ReturnType<typeof setTimeout> | null
    descriptor.value = (...args: any[]) => {
      if (timer) {
        clearTimeout(timer)
      }
      timer = setTimeout(() => method(...args), delay)
    }

    return descriptor
  }

class Test {
  @debounce(1000)
  sayHi() {
    // console.log(this.foo)
    console.log('hi')
  }
}

const test = new Test()
test.sayHi()
test.sayHi()
test.sayHi()
export {}

function funcDebounce(fn: (...args: any[]) => any, delay: number) {
  let timer: ReturnType<typeof setTimeout> | null

  // this的指向要跟原来函数一样
  return function (this: any, ...args: any[]) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      fn.call(this, ...args)
    }, delay)
  }
}

const testFunc = () => console.log('test')
const debouncedTestFunc = funcDebounce(testFunc, 1000)
debouncedTestFunc()

export { debounce }
