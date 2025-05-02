# 并发工具

- [Mutex](Mutex2.ts)
- [ReadWriteLock](ReadWriteLock.ts)
- [ReentrantLock](ReentrantLock.ts)
- [StampedLock](StampedLock.ts)

- [Semaphore](Semaphore.ts)
- [Cond](Cond2.ts)
- [Barrier](Barrier.ts)
- [CountDownLatch](CountDownLatch.ts)

- [BlockingQueue](BlockingQueue.ts)

- [Once](Once2.ts)

- [Pool](Pool.ts)

---

# Java 并发工具类的 TypeScript 实现解析

以下是对列出的各种并发工具的精简讲解：

## 1. Mutex (互斥锁)

**核心功能**：确保同一时间只有一个线程可以访问共享资源。

```typescript
// 基本使用
const mutex = new Mutex()
await mutex.lock()
try {
  // 临界区代码
} finally {
  mutex.unlock()
}

// 更优雅的使用
await mutex.withLock(async () => {
  // 临界区代码
})
```

**关键特性**：

- 异步锁定/解锁操作
- 支持超时锁定
- FIFO (先进先出) 等待队列
- 提供 withLock 便捷方法

## 2. ReentrantLock (可重入锁)

**核心功能**：允许同一线程多次获取同一个锁，无死锁风险。

```typescript
const lock = new ReentrantLock()
await lock.lock()
try {
  // 临界区代码
  await lock.lock() // 同一线程可以再次获取锁
  try {
    // 嵌套临界区
  } finally {
    lock.unlock()
  }
} finally {
  lock.unlock()
}
```

**关键特性**：

- 计数功能，跟踪同一线程获取锁的次数
- 线程身份跟踪，确保只有获取锁的线程才能释放锁
- 支持嵌套锁定场景

## 3. Semaphore (信号量)

**核心功能**：限制同时访问某一资源的线程数量。

```typescript
// 最多允许3个并发访问
const semaphore = new Semaphore(3)

await semaphore.acquire()
try {
  // 访问受限资源
} finally {
  semaphore.release()
}

// 批量获取释放
await semaphore.acquireMultiple(2) // 获取多个许可
semaphore.releaseMultiple(2) // 释放多个许可
```

**关键特性**：

- 可配置的许可证数量
- 支持获取/释放多个许可
- 支持尝试获取(非阻塞)和带超时的获取
- 适用于资源池和并发限制

## 4. Cond (条件变量)

**核心功能**：允许线程等待某个条件成立，并在条件满足时被其他线程唤醒。

```typescript
const mutex = new Mutex()
const condition = new Cond(mutex)

// 消费者
await mutex.lock()
try {
  while (!dataReady) {
    await condition.wait() // 自动释放锁并等待
  }
  // 处理数据
} finally {
  mutex.unlock()
}

// 生产者
await mutex.lock()
try {
  // 准备数据
  dataReady = true
  condition.signal() // 唤醒一个等待的线程
  // 或 condition.broadcast(); // 唤醒所有等待的线程
} finally {
  mutex.unlock()
}
```

**关键特性**：

- 与互斥锁协同工作
- 提供 wait/signal/broadcast 操作
- 支持带超时等待
- 自动处理锁的释放和重获取

## 5. Barrier (屏障)

**核心功能**：让一组线程在某个点同步等待，直到所有线程都到达该点后才继续执行。

```typescript
// 需要3个线程到达才能继续
const barrier = new Barrier(3)

async function worker(id) {
  console.log(`线程 ${id} 开始工作`)
  // 处理第一阶段

  await barrier.wait() // 等待所有线程完成第一阶段

  console.log(`线程 ${id} 继续执行`)
  // 处理第二阶段
}
```

**关键特性**：

- 可配置等待线程数
- 自动重置，可重复使用
- 支持超时等待
- 可选的触发回调函数
- 提供损坏状态处理机制

## 6. CountDownLatch (倒计时锁存器)

**核心功能**：允许一个或多个线程等待直到指定数量的操作完成。

```typescript
// 等待3个操作完成
const latch = new CountDownLatch(3)

// 等待线程
async function waiter() {
  await latch.await()
  console.log('所有操作完成，继续执行')
}

// 工作线程
function worker(id) {
  // 执行任务
  console.log(`任务 ${id} 完成`)
  latch.countDown() // 减少计数
}
```

**关键特性**：

- 一次性使用（计数到零后不能重置）
- 支持等待超时
- 提供计数批量减少方法
- 适用于启动信号和完成信号场景

## 7. BlockingQueue (阻塞队列)

**核心功能**：提供线程安全的队列操作，支持在队列空/满时阻塞。

```typescript
const queue = new BlockingQueue<Task>(10) // 容量为10的队列

// 生产者
async function producer() {
  while (true) {
    const task = generateTask()
    await queue.put(task) // 队列满时会阻塞
  }
}

// 消费者
async function consumer() {
  while (true) {
    const task = await queue.take() // 队列空时会阻塞
    processTask(task)
  }
}
```

**关键特性**：

- 可配置容量限制
- 提供阻塞操作(put/take)和非阻塞操作(offer/poll)
- 支持操作超时
- 优化的直通模式(当生产者和消费者同时等待时)
- 线程安全的队列管理

## 8. Once (一次性执行)

**核心功能**：确保某个函数只被执行一次，即使在多线程环境下也是如此。

```typescript
const once = new Once()

async function initialize() {
  await once.do(() => {
    // 只会执行一次的初始化代码
    console.log('Initializing...')
    return 'Initialized'
  })
  console.log('Already initialized')
}

// 多次调用，初始化代码只会执行一次
initialize()
initialize()
initialize()
```

**关键特性**：

- 保证回调函数只执行一次
- 所有调用共享同一个结果
- 错误传播机制
- 支持异步操作

## 9. Pool (资源池)

**核心功能**：管理和复用有限的资源，如数据库连接。

```typescript
// 创建一个连接池，最多5个连接
const pool = new Pool<Connection>(
  5,
  // 创建连接的工厂函数
  () => createConnection(),
  // 可选的销毁函数
  conn => conn.close()
)

async function query(sql) {
  // 从池中获取连接
  const connection = await pool.acquire()
  try {
    // 使用连接
    return connection.execute(sql)
  } finally {
    // 归还连接到池中
    pool.release(connection)
  }
}
```

**关键特性**：

- 可配置最大资源数量
- 自动创建资源
- 支持资源验证和清理
- 提供资源借用和归还机制
- 支持池耗尽时的等待队列
- 适用于连接池、线程池等场景

## 总结

这些并发工具提供了一套完整的线程同步和协作机制，使多线程程序能够正确、高效地运行：

1. **互斥锁和可重入锁**：提供基本的排他性访问控制
2. **信号量**：控制并发访问数量
3. **条件变量**：支持基于条件的等待和唤醒
4. **栅栏和倒计时锁存器**：提供线程协作和同步点
5. **阻塞队列**：实现线程间安全的数据传递
6. **一次性执行**：确保初始化操作只执行一次
7. **资源池**：管理有限资源的分配和回收

这些工具在TypeScript中实现了Java并发工具包的核心功能，适合在Node.js或浏览器环境中构建高效可靠的并发应用程序。
