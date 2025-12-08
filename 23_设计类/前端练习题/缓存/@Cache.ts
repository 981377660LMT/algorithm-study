/* eslint-disable no-promise-executor-return */
/* eslint-disable no-inner-declarations */
/* eslint-disable no-lone-blocks */
/* eslint-disable @typescript-eslint/no-explicit-any */

/**
 * 缓存配置接口
 */
export interface CacheOptions {
  ttl?: number
  strategy?: 'lazy' | 'periodic'
  cleanupInterval?: number
  capacity?: number
  /** 自定义 Key 生成器，解决 JSON.stringify 性能或循环引用问题 */
  keyGenerator?: (...args: any[]) => string
}

export type CachedMethod<T extends (...args: any[]) => any> = T & {
  clearCache: () => void
  disposeCache: () => void
}

type Timer = ReturnType<typeof setTimeout>

class CacheManager {
  private store = new Map<string, { value: any; expiry: number }>()
  private timer: Timer | null = null

  constructor(private options: Required<CacheOptions>) {
    if (options.strategy === 'periodic') {
      this.startCleanupTask()
    }
  }

  get(key: string): any | undefined {
    const item = this.store.get(key)
    if (!item) return undefined
    if (Date.now() > item.expiry) {
      this.store.delete(key)
      return undefined
    }
    this.store.delete(key)
    this.store.set(key, item)
    return item.value
  }

  set(key: string, value: any) {
    if (this.store.has(key)) {
      this.store.delete(key)
    } else if (this.store.size >= this.options.capacity) {
      const firstKey = this.store.keys().next().value
      if (firstKey) this.store.delete(firstKey)
    }

    this.store.set(key, {
      value,
      expiry: Date.now() + this.options.ttl
    })
  }

  delete(key: string) {
    this.store.delete(key)
  }

  clear() {
    this.store.clear()
  }

  dispose() {
    this.clear()
    this.stopCleanupTask()
  }

  private startCleanupTask() {
    this.stopCleanupTask()
    this.timer = setInterval(() => {
      const now = Date.now()
      for (const [key, item] of this.store.entries()) {
        if (now > item.expiry) {
          this.store.delete(key)
        }
      }
    }, this.options.cleanupInterval)

    // 兼容性处理：Node 环境下不阻塞进程退出
    if (this.timer && typeof (this.timer as any).unref === 'function') {
      ;(this.timer as any).unref()
    }
  }

  private stopCleanupTask() {
    if (this.timer) {
      clearInterval(this.timer)
      this.timer = null
    }
  }
}

export function Cache(optionsOrTtl: number | CacheOptions = 60000) {
  const defaultOptions: Required<CacheOptions> = {
    ttl: 60000,
    strategy: 'lazy',
    cleanupInterval: 60000,
    capacity: 1000,
    keyGenerator: (...args: any[]) => {
      try {
        return JSON.stringify(args)
      } catch (e) {
        return null as any
      }
    }
  }

  let options: Required<CacheOptions>
  if (typeof optionsOrTtl === 'number') {
    options = { ...defaultOptions, ttl: optionsOrTtl }
  } else {
    options = { ...defaultOptions, ...optionsOrTtl }
  }

  return function <T extends (...args: any[]) => any>(
    target: object,
    propertyKey: string | symbol,
    descriptor: TypedPropertyDescriptor<T>
  ): TypedPropertyDescriptor<T> | void {
    const originalMethod = descriptor.value
    if (!originalMethod) throw new Error('@Cache can only be applied to methods.')

    // 注意：这个 manager 是类级别的（所有实例共享）
    const manager = new CacheManager(options)

    const wrappedMethod = function (this: any, ...args: Parameters<T>): ReturnType<T> {
      // 即使类名被压缩，只要方法名不重复（通常不会），结合 args 也能唯一确定
      const argsKey = options.keyGenerator(...args)
      // 如果 Key 生成失败（如循环引用），直接跳过缓存，防止碰撞
      if (argsKey === null) {
        console.warn(`@Cache: Key generation failed for ${String(propertyKey)}, skipping cache.`)
        return originalMethod.apply(this, args)
      }

      const cacheKey = `${String(propertyKey)}:${argsKey}`
      const cachedValue = manager.get(cacheKey)
      if (cachedValue !== undefined) {
        return cachedValue
      }

      const result = originalMethod.apply(this, args)

      // 并发处理：立即缓存 Promise
      if (result && result instanceof Promise) {
        const cachedPromise = result.catch((err: any) => {
          manager.delete(cacheKey)
          throw err
        })
        manager.set(cacheKey, cachedPromise)
        return cachedPromise as ReturnType<T>
      }

      manager.set(cacheKey, result)
      return result as ReturnType<T>
    }

    const augmentedMethod = wrappedMethod as CachedMethod<T>
    augmentedMethod.clearCache = () => manager.clear()
    augmentedMethod.disposeCache = () => manager.dispose()

    descriptor.value = augmentedMethod
    return descriptor
  }
}

{
  const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

  class UserService {
    // 缓存 2 秒
    @Cache(2000)
    async getUserInfo(id: number) {
      console.log(`[${new Date().toISOString()}] 正在从数据库查询 User ${id}... (模拟耗时)`)
      await delay(500)
      return { id, name: `User_${id}`, role: 'admin' }
    }
  }

  async function testAsync() {
    const service = new UserService()

    console.log('--- 1. 测试并发请求合并 ---')
    // 同时发起两个请求，理论上只会打印一次 "正在从数据库查询"
    const p1 = service.getUserInfo(1)
    const p2 = service.getUserInfo(1)

    const [r1, r2] = await Promise.all([p1, p2])
    console.log('结果是否相同:', r1 === r2) // true

    console.log('\n--- 2. 测试缓存命中 ---')
    await service.getUserInfo(1) // 此时还在 2秒 缓存期内，不会打印查询日志

    console.log('\n--- 3. 测试过期重新查询 ---')
    await delay(2100) // 等待过期
    await service.getUserInfo(1) // 应该再次打印查询日志
  }

  testAsync()
}
