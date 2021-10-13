type Func<T = any> = (...args: any[]) => T
type PromiseFunc<T = any> = (...args: any[]) => Promise<T>
interface Class<InstanceType = any> extends Function {
  new (...args: any[]): InstanceType
}

export { Func, PromiseFunc, Class }
