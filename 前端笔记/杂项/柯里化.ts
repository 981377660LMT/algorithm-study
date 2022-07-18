const add = (x: number, y: number) => x + y

const curriedAdd = curry(add)

console.log(curriedAdd(1, 2))
console.log(curriedAdd(1)(2))

function curry<T>(this: T, func: (...args: any[]) => any) {
  const ctx = this
  const need = func.length

  return function inner(...args: any[]): any {
    // !传的满了,call
    if (args.length >= need) return func.call(ctx, ...args)
    // !传的没有满,bind
    return inner.bind(ctx, ...args)
  }
}
