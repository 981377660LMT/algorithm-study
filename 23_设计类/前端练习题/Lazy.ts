/* eslint-disable prefer-destructuring */
/* eslint-disable no-inner-declarations */

/**
 * 惰性求值.
 * vscode中的实现.
 */
class Lazy<V> {
  private readonly _executor: () => V
  private _didRun = false
  private _value?: V
  private _error: Error | undefined

  constructor(executor: () => V) {
    this._executor = executor
  }

  /**
   * True if the lazy value has been resolved.
   */
  get hasValue(): boolean {
    return this._didRun
  }

  /**
   * Get the wrapped value.
   *
   * This will force evaluation of the lazy value if it has not been resolved yet.
   * Lazy values are only resolved once.
   * `getValue` will re-throw exceptions that are hit while resolving the value.
   */
  get value(): V {
    if (!this._didRun) {
      try {
        this._value = this._executor()
      } catch (error) {
        this._error = error as Error
      } finally {
        this._didRun = true
      }
    }

    if (this._error) {
      throw this._error
    }
    return this._value!
  }

  /**
   * Get the wrapped value without forcing evaluation.
   */
  get rawValue(): V | undefined {
    return this._value
  }
}

export {}

if (require.main === module) {
  // 定义一个需要延迟加载的函数
  function expensiveCalculation(): number {
    console.log('Executing expensiveCalculation')
    let result = 0
    for (let i = 0; i < 1000000000; i++) {
      result += i
    }
    return result
  }

  // 创建一个 Lazy 对象
  const lazyValue = new Lazy(() => expensiveCalculation())

  // 访问 Lazy 对象的值
  console.log(lazyValue.hasValue) // false

  const value1 = lazyValue.value // 这里会触发 expensiveCalculation 的执行
  console.log(value1) // 输出 expensiveCalculation 的返回值
  console.log(lazyValue.hasValue) // true

  const value2 = lazyValue.value // 这里不会再次执行 expensiveCalculation
  console.log(value2) // 输出 expensiveCalculation 的返回值

  const rawValue = lazyValue.rawValue // 可以直接访问缓存的计算结果
  console.log(rawValue) // 输出 expensiveCalculation 的返回值
}
