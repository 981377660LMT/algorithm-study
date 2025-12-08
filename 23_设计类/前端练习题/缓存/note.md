这是一个非常完善的 `@Cache` 实现。为了让你更好地理解如何使用它，我为你编写了几个典型的使用案例，涵盖了**异步并发**、**同步缓存**、**手动清理**以及**高级配置**。

你可以将以下代码保存为 `example.ts` 并直接运行。

### 1. 基础用法：异步请求与并发合并

这是最常见的场景。当多个请求同时发起时，装饰器会自动合并它们，只执行一次底层逻辑。

```typescript
import { Cache } from './Cache' // 假设你的文件名为 Cache.ts

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
```

### 2. 进阶用法：手动清除缓存 (TypeScript)

当你更新了数据，需要强制刷新缓存时使用。注意 TypeScript 类型断言的用法。

```typescript
import { Cache, CachedMethod } from './Cache'

class ProductService {
  @Cache(60000) // 缓存 1 分钟
  getProduct(id: number) {
    console.log(`查询商品 ${id}...`)
    return { id, price: 100 }
  }

  updateProduct(id: number, newPrice: number) {
    console.log(`更新商品 ${id} 价格为 ${newPrice}`)
    // ... 更新数据库逻辑 ...

    // 【关键点】手动清除缓存
    // 需要断言为 CachedMethod 类型才能看到 clearCache 方法
    const cachedMethod = this.getProduct as CachedMethod<typeof this.getProduct>
    cachedMethod.clearCache()
    console.log('缓存已清除')
  }
}

function testManualClear() {
  const service = new ProductService()

  service.getProduct(101) // 第一次：执行逻辑
  service.getProduct(101) // 第二次：走缓存

  service.updateProduct(101, 200) // 更新并清除

  service.getProduct(101) // 第三次：重新执行逻辑
}

testManualClear()
```

### 3. 高级配置：LRU 容量限制与自定义 Key

适用于参数复杂（对象）或内存敏感的场景。

```typescript
import { Cache } from './Cache'

interface QueryParams {
  category: string
  page: number
}

class ReportService {
  @Cache({
    ttl: 5000,
    capacity: 2, // 最多只存 2 个结果，多余的会被挤掉 (LRU)
    // 自定义 Key 生成器：只根据 category 缓存，忽略 page (仅作演示)
    keyGenerator: (params: QueryParams) => `cat_${params.category}`
  })
  getReport(params: QueryParams) {
    console.log(`生成报表: ${params.category} - Page ${params.page}`)
    return `Report Data`
  }
}

function testAdvanced() {
  const service = new ReportService()

  console.log('--- 测试 LRU 容量 ---')
  service.getReport({ category: 'A', page: 1 }) // 缓存 A (Size: 1)
  service.getReport({ category: 'B', page: 1 }) // 缓存 B (Size: 2)
  service.getReport({ category: 'C', page: 1 }) // 缓存 C (Size: 2), A 被挤掉!

  service.getReport({ category: 'B', page: 1 }) // 命中 B
  service.getReport({ category: 'A', page: 1 }) // A 之前被挤掉了，重新执行
}

testAdvanced()
```

### 4. 资源释放：定期清理策略

如果你使用了

strategy: 'periodic'

，在组件销毁时最好手动停止定时器。

```typescript
import { Cache, CachedMethod } from './Cache'

class MonitorService {
  @Cache({
    strategy: 'periodic', // 启用后台定时清理
    cleanupInterval: 1000, // 每秒检查一次
    ttl: 2000
  })
  getMetrics() {
    return 'metrics'
  }

  destroy() {
    // 停止后台定时器，防止内存泄漏
    ;(this.getMetrics as CachedMethod<any>).disposeCache()
    console.log('MonitorService 销毁，缓存定时器已停止')
  }
}
```

### 总结

1.  **默认情况**：直接用 `@Cache(毫秒数)`，最简单，性能最好（懒删除）。
2.  **并发安全**：对于返回 Promise 的方法，装饰器自动处理了并发请求，不会导致重复查询。
3.  **类型提示**：调用

clearCache

时，记得使用

as CachedMethod<typeof ...>

来获取正确的类型提示。
