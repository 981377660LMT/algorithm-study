以下是一个使用 `Cond` 类的示例，演示如何通过条件变量协调两个异步任务：

```typescript
import { Mutex } from 'async-mutex' // 假设使用 async-mutex 库的 Mutex

async function example() {
  const mutex = new Mutex()
  const cond = new Cond(mutex)
  let sharedState = false

  // 等待条件的异步任务
  async function waiter() {
    await mutex.lock()
    try {
      console.log('[Waiter] 获取锁，检查条件')
      while (!sharedState) {
        console.log('[Waiter] 条件未满足，进入等待...')
        await cond.wait() // 释放锁并等待
      }
      console.log('[Waiter] 条件满足，继续执行')
    } finally {
      mutex.unlock()
    }
  }

  // 修改条件并通知的异步任务
  async function notifier() {
    await mutex.lock()
    try {
      console.log('[Notifier] 获取锁，修改条件')
      sharedState = true
      cond.signal() // 通知一个等待者
      // cond.broadcast(); // 通知所有等待者
      console.log('[Notifier] 已发送通知')
    } finally {
      mutex.unlock()
    }
  }

  // 启动两个任务
  waiter()
  setTimeout(notifier, 2000) // 2秒后触发通知
}

example()
```

### 关键点说明：

1. **锁的获取顺序**：

   - `waiter` 和 `notifier` 在操作共享状态前必须先获取锁
   - `cond.wait()` 会自动释放锁，唤醒时自动重新获取锁

2. **条件检查模式**：

   ```typescript
   while (!sharedState) {
     await cond.wait()
   }
   ```

   这是标准条件变量用法，用循环防止虚假唤醒

3. **执行流程**：
   ```
   [Waiter] 获取锁，检查条件
   [Waiter] 条件未满足，进入等待...
   [Notifier] 获取锁，修改条件
   [Notifier] 已发送通知
   [Waiter] 条件满足，继续执行
   ```

### 典型使用场景：

- 生产者/消费者队列
- 资源池管理
- 异步任务协调

### 注意事项：

1. 必须在持有锁时调用 `signal()`/`broadcast()`
2. 推荐将 `cond.wait()` 放在 `while` 循环中检查条件
3. 使用 `finally` 块确保锁被释放
4. 如果使用 `broadcast()`，所有等待者会同时竞争锁，只有一个能立即获得锁
