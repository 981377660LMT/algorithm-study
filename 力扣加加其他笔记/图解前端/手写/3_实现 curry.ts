// test
function test(a: number, b: number, c: number) {
  console.log(a, b, c)
}

const f1 = curry(test)(1)
const f2 = f1(2)
f2(3)

function curry<T>(this: T, fn: (...args: any[]) => any) {
  const ctx = this
  return function inner(...args: any[]): any {
    return args.length === fn.length ? fn.call(ctx, ...args) : inner.bind(ctx, ...args)
  }
}

export {}
