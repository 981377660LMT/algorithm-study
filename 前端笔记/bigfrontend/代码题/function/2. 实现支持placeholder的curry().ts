interface Curry {
  (fn: (...args: any[]) => any): (...args: any[]) => any
  placeholder: Symbol
}

const curry: Curry = fn => {
  const ctx = this

  return function inner(this: any, ...args: any[]): any {
    const isComplete =
      args.length >= fn.length && !args.slice(0, fn.length).includes(curry.placeholder)

    if (isComplete) return fn.call(ctx, ...args)

    return (...newArgs: any[]) => {
      const res = args.map(arg =>
        arg === curry.placeholder && newArgs.length ? newArgs.shift()! : arg
      )
      console.log(newArgs, res)
      return inner(...res, ...newArgs)
    }
  }
}

curry.placeholder = Symbol()

const join = (a: string, b: string, c: string) => {
  return `${a}_${b}_${c}`
}

const curriedJoin = curry(join)
const _ = curry.placeholder

console.log(curriedJoin(_, _, _)(1)(_, 3)(2)) // '1_2_3'

export {}
