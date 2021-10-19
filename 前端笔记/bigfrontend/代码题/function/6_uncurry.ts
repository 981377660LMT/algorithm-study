export {}
const add = (x: number) => (y: number) => (z: number) => x + y + z
const uncurriedAdd = uncurry(add)
console.log(uncurriedAdd(1, 2, 3)) // 6
type Function = (...args: any[]) => any

/**
 * Uncurries a function up to depth n.
 * Return a `variadic function.`
 */
function uncurry(curriedFunc: Function): Function {
  return (...args: any[]) => {
    const next =
      (curriedFunc: Function) =>
      (...args: any[]) =>
        args.reduce((curriedFunc, arg) => curriedFunc(arg), curriedFunc)

    return next(curriedFunc)(...args)
  }
}
