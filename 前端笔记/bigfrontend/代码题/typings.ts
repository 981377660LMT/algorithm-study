type Func<T = any> = (...args: any[]) => T
type PromiseFunc<T = any> = (...args: any[]) => Promise<T>

interface ClassConstructor<InstanceType = any> extends Function {
  new (...args: any[]): InstanceType
  readonly prototype: ClassConstructor<InstanceType>
}

export { Func, PromiseFunc, ClassConstructor }
