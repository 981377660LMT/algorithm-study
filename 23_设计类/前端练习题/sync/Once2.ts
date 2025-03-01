// TypeScript Once 实现
//
// Once 是一个同步原语，用于确保某个函数或代码块只会被执行一次，无论它被调用多少次。这在初始化、配置加载和单例模式等场景非常有用。
//
// 这个 Once 实现提供了以下功能：
//
// 1. 确保函数只执行一次，后续调用返回相同结果
// 2. 支持异步函数
// 3. 处理执行错误并在后续调用中重新抛出
// 4. 提供重置功能
// 5. 线程安全的执行（在 JavaScript 的单线程上下文中）
// 6. 提供便捷的装饰器，用于类方法
//
// Once 在许多场景中非常有用，例如：
// - 应用程序配置加载
// - 数据库或资源连接初始化
// - 单例模式实现
// - 昂贵计算的缓存

/**
 * Once 类 - 确保函数只执行一次
 */
class Once<T> {
  private done = false
  private result: T | undefined
  private error: Error | undefined
  private inProgress = false
  private waiters: Array<{
    resolve: (value: T) => void
    reject: (reason: Error) => void
  }> = []

  /**
   * 创建一个Once实例
   * @param fn 要执行的函数
   */
  // eslint-disable-next-line no-useless-constructor
  constructor(private fn: () => T | Promise<T>) {}

  /**
   * 执行函数（如果尚未执行过）并返回结果
   */
  async do(): Promise<T> {
    // 如果已经执行过，直接返回结果或抛出已保存的错误
    if (this.done) {
      if (this.error) {
        throw this.error
      }
      return this.result as T
    }

    // 如果正在执行中，等待执行完成
    if (this.inProgress) {
      return new Promise<T>((resolve, reject) => {
        this.waiters.push({ resolve, reject })
      })
    }

    // 标记为正在执行
    this.inProgress = true

    try {
      // 执行函数
      const result = await Promise.resolve(this.fn())
      this.result = result
      this.done = true

      // 通知所有等待者
      this.waiters.forEach(waiter => waiter.resolve(result))
      this.waiters = []

      return result
    } catch (err) {
      this.error = err as Error
      this.done = true

      // 通知所有等待者出错
      this.waiters.forEach(waiter => waiter.reject(err as Error))
      this.waiters = []

      throw err
    } finally {
      this.inProgress = false
    }
  }

  /**
   * 检查函数是否已被执行
   */
  isDone(): boolean {
    return this.done
  }

  /**
   * 重置Once实例，允许函数再次执行
   * @throws Error 如果函数正在执行中
   */
  reset(): void {
    if (this.inProgress) {
      throw new Error('Cannot reset while execution is in progress')
    }
    this.done = false
    this.result = undefined
    this.error = undefined
  }
}

/**
 * 创建一个只会执行一次的函数
 * @param fn 要执行的函数
 * @returns 包装后的函数
 */
function once<T>(fn: () => T | Promise<T>): () => Promise<T> {
  const o = new Once(fn)
  return () => o.do()
}

/**
 * 方法装饰器，确保方法只执行一次
 */
function OnceDecorator() {
  return function (_target: any, _propertyKey: string, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value
    const onceInstances = new WeakMap<any, Once<any>>()

    descriptor.value = function (...args: any[]) {
      // 为每个实例创建独立的Once
      if (!onceInstances.has(this)) {
        onceInstances.set(this, new Once(() => originalMethod.apply(this, args)))
      }
      return onceInstances.get(this)!.do()
    }

    return descriptor
  }
}

// 使用示例
async function example() {
  // 基本用法
  const initDatabase = once(async () => {
    console.log('数据库初始化...')
    await new Promise(resolve => setTimeout(resolve, 1000))
    console.log('数据库初始化完成')
    return { connection: 'db_connection_object' }
  })

  // 多次调用只会执行一次
  const connection1 = await initDatabase()
  const connection2 = await initDatabase()
  console.log('连接相同:', connection1 === connection2) // true

  // 使用Once类
  class ConfigLoader {
    private configLoader = new Once(async () => {
      console.log('加载配置...')
      await new Promise(resolve => setTimeout(resolve, 500))
      return { apiKey: 'secret_key' }
    })

    async getConfig() {
      return this.configLoader.do()
    }

    isConfigLoaded() {
      return this.configLoader.isDone()
    }

    resetConfigLoader() {
      this.configLoader.reset()
    }
  }

  // 使用装饰器
  class Service {
    @OnceDecorator()
    async initialize() {
      console.log('服务初始化...')
      await new Promise(resolve => setTimeout(resolve, 300))
      return 'service ready'
    }
  }
}

example()

export { once }
