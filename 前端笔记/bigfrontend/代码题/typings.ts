type Func<T = any> = (...args: any[]) => T
type PromiseFunc<T = any> = (...args: any[]) => Promise<T>
interface Class<T = any> extends Function {
  new (...args: any[]): T
}

export { Func, PromiseFunc, Class }
