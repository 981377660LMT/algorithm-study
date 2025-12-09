/**
 * 缓存配置接口
 */
export interface CacheOptions {
  /**
   * 缓存存活时间 (毫秒). 默认为 60000ms (1分钟).
   * 设置为 -1 表示永久缓存.
   */
  ttl?: number

  /**
   * 自定义生成缓存 Key 的函数.
   * 默认使用 `ClassName:MethodName:JSON.stringify(args)`
   */
  keyGenerator?: (...args: any[]) => string
}

interface CacheEntry {
  value: any
  expiry: number
}

/**
 * 内存缓存存储 (简单的 Map 实现)
 * 注意：这是全局单例存储，所有使用该装饰器的类共享此 Map，但 Key 会区分。
 */
const cacheStore = new Map<string, CacheEntry>()

/**
 * @Cache 装饰器
 * 用于缓存类方法的返回值，支持异步 Promise 和同步结果。
 */
export function Cache(options: CacheOptions = {}) {
  const ttl = options.ttl ?? 60000 // 默认 1 分钟

  return function (target: any, propertyKey: string, descriptor: PropertyDescriptor) {
    const originalMethod = descriptor.value
    const className = target.constructor.name

    descriptor.value = function (...args: any[]) {
      // 1. 生成缓存 Key
      let cacheKey = ''
      if (options.keyGenerator) {
        cacheKey = options.keyGenerator(...args)
      } else {
        // 默认 Key: 类名:方法名:参数JSON
        // 注意：JSON.stringify 对于循环引用或特殊对象可能不适用，生产环境可换成 object-hash 等库
        try {
          cacheKey = `${className}:${propertyKey}:${JSON.stringify(args)}`
        } catch (e) {
          console.warn(`[Cache] Failed to generate key for ${propertyKey}, skipping cache.`)
          return originalMethod.apply(this, args)
        }
      }

      // 2. 检查缓存是否存在且未过期
      const cachedEntry = cacheStore.get(cacheKey)
      const now = Date.now()

      if (cachedEntry) {
        if (ttl === -1 || cachedEntry.expiry > now) {
          // console.debug(`[Cache] Hit: ${cacheKey}`);
          // 如果缓存的是 Promise (异步方法)，需要克隆一个新的 Promise 防止状态被消费
          if (cachedEntry.value instanceof Promise) {
            return cachedEntry.value.then((res: any) => JSON.parse(JSON.stringify(res)))
          }
          // 对于对象，返回深拷贝防止外部修改污染缓存
          if (typeof cachedEntry.value === 'object' && cachedEntry.value !== null) {
            return JSON.parse(JSON.stringify(cachedEntry.value))
          }
          return cachedEntry.value
        } else {
          // 过期删除
          cacheStore.delete(cacheKey)
        }
      }

      // 3. 执行原方法
      const result = originalMethod.apply(this, args)

      // 4. 处理返回值并存入缓存
      if (result instanceof Promise) {
        // 异步方法：缓存 Promise 对象，但在 catch 时移除缓存
        // 这里的 result 是原始 Promise
        const cachedPromise = result
          .then(data => {
            // 成功后，更新缓存值为实际数据（可选，或者一直存 Promise）
            // 这里选择保持 Promise 结构，但在 resolve 时深拷贝数据返回
            return data
          })
          .catch(err => {
            // 如果请求失败，立即移除缓存，以便下次重试
            cacheStore.delete(cacheKey)
            throw err
          })

        // 存入缓存
        cacheStore.set(cacheKey, {
          value: cachedPromise,
          expiry: ttl === -1 ? Number.MAX_SAFE_INTEGER : now + ttl
        })

        return cachedPromise
      } else {
        // 同步方法
        cacheStore.set(cacheKey, {
          value: result,
          expiry: ttl === -1 ? Number.MAX_SAFE_INTEGER : now + ttl
        })
        return result
      }
    }

    return descriptor
  }
}

/**
 * 清除指定 Key 前缀的缓存 (辅助工具)
 */
export function clearCache(prefix?: string) {
  if (!prefix) {
    cacheStore.clear()
    return
  }
  for (const key of cacheStore.keys()) {
    if (key.startsWith(prefix)) {
      cacheStore.delete(key)
    }
  }
}
