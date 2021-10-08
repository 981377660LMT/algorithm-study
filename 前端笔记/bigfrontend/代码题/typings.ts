type Func<T = unknown> = (...args: any[]) => T
type AsyncFunc<T = unknown> = (...args: any[]) => Promise<T>

export { Func, AsyncFunc }
