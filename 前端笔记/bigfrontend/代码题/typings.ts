type Func<T = any> = (...args: any[]) => T
type PromiseFunc<T = unknown> = (...args: any[]) => Promise<T>

export { Func, PromiseFunc }
