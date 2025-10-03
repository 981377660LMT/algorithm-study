我们将模拟一个缓存客户端（如 Redis）和一个数据库客户端（如 MySQL），并展示不同策略下的读写操作。

### 1. 基础设置：模拟客户端和数据模型

首先，我们创建一些模拟的客户端和服务，以便于演示。

```typescript
// --- 数据模型 ---
interface Product {
  id: string
  name: string
  price: number
  version: number // 用于演示数据变化
}

// --- 模拟数据库客户端 ---
class MockDBClient {
  private storage = new Map<string, Product>()

  constructor() {
    // 初始化一条数据
    this.storage.set('prod-123', { id: 'prod-123', name: 'Laptop', price: 1000, version: 1 })
  }

  async get(id: string): Promise<Product | null> {
    console.log(`[DB] 正在查询 ID: ${id}`)
    await new Promise(resolve => setTimeout(resolve, 100)) // 模拟DB延迟
    const product = this.storage.get(id)
    return product ? { ...product } : null
  }

  async update(product: Product): Promise<Product> {
    console.log(`[DB] 正在更新 ID: ${product.id} 到版本 ${product.version}`)
    await new Promise(resolve => setTimeout(resolve, 150)) // 模拟DB延迟
    this.storage.set(product.id, product)
    return { ...product }
  }
}

// --- 模拟缓存客户端 ---
class MockCacheClient {
  private storage = new Map<string, string>()
  private readonly CACHE_TTL_SECONDS = 60

  async get(key: string): Promise<string | null> {
    console.log(`[Cache] 正在获取 Key: ${key}`)
    await new Promise(resolve => setTimeout(resolve, 20)) // 模拟缓存延迟
    return this.storage.get(key) || null
  }

  async set(key: string, value: string): Promise<void> {
    console.log(`[Cache] 正在设置 Key: ${key}`)
    await new Promise(resolve => setTimeout(resolve, 20))
    this.storage.set(key, value)
  }

  async del(key: string): Promise<void> {
    console.log(`[Cache] 正在删除 Key: ${key}`)
    await new Promise(resolve => setTimeout(resolve, 20))
    this.storage.delete(key)
  }
}

// --- 模拟消息队列客户端 ---
type MessageHandler = (message: any) => Promise<void>

class MockMQClient {
  private queue: any[] = []
  private subscribers: MessageHandler[] = []

  async publish(message: any): Promise<void> {
    console.log(`[MQ] 发布消息: ${JSON.stringify(message)}`)
    this.queue.push(message)
    // 立即异步触发消费
    this.processQueue()
  }

  subscribe(handler: MessageHandler): void {
    this.subscribers.push(handler)
  }

  private async processQueue(): Promise<void> {
    if (this.queue.length > 0) {
      const message = this.queue.shift()
      console.log(`[MQ] 消费者正在处理消息...`)
      for (const handler of this.subscribers) {
        await handler(message)
      }
    }
  }
}

const db = new MockDBClient()
const cache = new MockCacheClient()
const mq = new MockMQClient()
```

### 2. 方案一：Cache-Aside Pattern (先更新 DB，再删缓存)

这是最经典和常用的模式。

```typescript
// ...existing code...
// --- 方案一: Cache-Aside Pattern ---

class ProductServiceV1 {
  // 读操作
  async getProduct(id: string): Promise<Product | null> {
    // 1. 先读缓存
    const cachedData = await cache.get(id)
    if (cachedData) {
      console.log('✅ 缓存命中!')
      return JSON.parse(cachedData) as Product
    }

    console.log('❌ 缓存未命中!')
    // 2. 缓存未命中，读数据库
    const product = await db.get(id)

    if (product) {
      // 3. 将数据回写到缓存
      await cache.set(id, JSON.stringify(product))
    }

    return product
  }

  // 写操作
  async updateProduct(id: string, newPrice: number): Promise<Product | null> {
    let product = await db.get(id)
    if (!product) {
      console.error('产品不存在')
      return null
    }
    product.price = newPrice
    product.version += 1

    // 1. 先更新数据库
    const updatedProduct = await db.update(product)

    // 2. 再删除缓存
    // 增加重试机制来提高删除成功的概率
    try {
      await this.deleteCacheWithRetry(id, 3)
    } catch (error) {
      // 如果重试后仍然失败，需要记录日志并告警，进行人工干预
      console.error(`🚨 删除缓存失败 (ID: ${id})! 需要人工干预!`, error)
    }

    return updatedProduct
  }

  private async deleteCacheWithRetry(key: string, retries: number): Promise<void> {
    for (let i = 0; i < retries; i++) {
      try {
        await cache.del(key)
        console.log(`✅ 删除缓存成功 (Key: ${key})`)
        return
      } catch (e) {
        console.warn(`[重试 ${i + 1}/${retries}] 删除缓存失败...`)
        if (i === retries - 1) throw e
      }
    }
  }
}
```

### 3. 方案二：基于消息队列的最终一致性

这是大型系统中保证最终一致性的黄金方案。

```typescript
// ...existing code...
// --- 方案二: 基于消息队列的最终一致性 ---

class CacheDeletionConsumer {
  constructor() {
    // 订阅消息队列，处理缓存删除任务
    mq.subscribe(this.handleMessage.bind(this))
    console.log('缓存删除消费者已启动并订阅消息...')
  }

  private async handleMessage(message: { key: string }): Promise<void> {
    if (message && message.key) {
      console.log(`[消费者] 收到删除缓存任务, Key: ${message.key}`)
      // 同样可以加入重试逻辑
      try {
        await cache.del(message.key)
        console.log(`[消费者] ✅ 成功删除缓存, Key: ${message.key}`)
      } catch (error) {
        console.error(`[消费者] 🚨 删除缓存失败, Key: ${message.key}`, error)
        // 在实际的MQ中，如果处理失败，消息会根据策略重回队列或进入死信队列
      }
    }
  }
}

// 在应用启动时初始化消费者
const cacheConsumer = new CacheDeletionConsumer()

class ProductServiceV2 {
  // 读操作与 V1 相同
  async getProduct(id: string): Promise<Product | null> {
    // ... (代码同 ProductServiceV1.getProduct)
    const cachedData = await cache.get(id)
    if (cachedData) {
      return JSON.parse(cachedData) as Product
    }
    const product = await db.get(id)
    if (product) {
      await cache.set(id, JSON.stringify(product))
    }
    return product
  }

  // 写操作
  async updateProduct(id: string, newPrice: number): Promise<Product | null> {
    let product = await db.get(id)
    if (!product) return null

    product.price = newPrice
    product.version += 1

    // 1. 先更新数据库
    const updatedProduct = await db.update(product)

    // 2. 发送一个删除缓存的消息到MQ
    await mq.publish({ key: id })

    return updatedProduct
  }
}
```

### 4. 方案三：订阅数据库变更日志 (Canal) - 概念演示

这个方案在应用层代码最简洁，但架构复杂。我们用一个模拟的 `BinlogSubscriber` 来演示其概念。

```typescript
// ...existing code...
// --- 方案三: 订阅 Binlog (概念演示) ---

// 模拟 Canal 这类中间件，它会监听数据库变更并触发事件
class MockBinlogSubscriber {
  private handler: (tableName: string, data: any) => void = () => {}

  onUpdate(handler: (tableName: string, data: any) => void) {
    this.handler = handler
  }

  // 模拟当数据库发生变更时，Canal会调用这个方法
  triggerUpdate(tableName: string, data: any) {
    console.log(`[Binlog] 监听到表 '${tableName}' 的数据变更`)
    this.handler(tableName, data)
  }
}

const binlogSubscriber = new MockBinlogSubscriber()

// 应用服务订阅 Binlog 事件
binlogSubscriber.onUpdate(async (tableName, data) => {
  if (tableName === 'products' && data.id) {
    console.log(`[Binlog Handler] 收到产品表变更, ID: ${data.id}。准备删除缓存。`)
    await cache.del(data.id)
  }
})

class ProductServiceV3 {
  // 读操作与 V1/V2 相同
  async getProduct(id: string): Promise<Product | null> {
    // ... (代码同 ProductServiceV1.getProduct)
    const cachedData = await cache.get(id)
    if (cachedData) return JSON.parse(cachedData) as Product
    const product = await db.get(id)
    if (product) await cache.set(id, JSON.stringify(product))
    return product
  }

  // 写操作变得非常纯粹，只关心数据库
  async updateProduct(id: string, newPrice: number): Promise<Product | null> {
    let product = await db.get(id)
    if (!product) return null

    product.price = newPrice
    product.version += 1

    // 只更新数据库
    const updatedProduct = await db.update(product)

    // 模拟 Canal 捕获到这次变更
    binlogSubscriber.triggerUpdate('products', updatedProduct)

    return updatedProduct
  }
}
```
