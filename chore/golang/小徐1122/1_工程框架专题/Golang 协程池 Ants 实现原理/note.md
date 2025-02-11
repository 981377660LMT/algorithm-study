# Golang 协程池 Ants 实现原理

https://www.bilibili.com/video/BV1wk4y1h7k7
https://mp.weixin.qq.com/s/Uctu_uKHk5oY0EtSZGUvsA

## 1. 前置知识点

1. sync.Locker
   在 ants 中，作者不希望使用 Mutex 这种重锁，而是自定义实现了一种轻量级的自旋锁 `spinLock`。
   ![alt text](image.png)

   加锁通过 for 循环 + cas 操作实现自旋，无需操作系统介入执行 park 操作；通过变量 backoff 反映抢锁激烈度，每次抢锁失败，执行 backoff 次让 cpu 时间片动作；backoff 随失败次数逐渐升级，封顶 16.

   ```go
   type spinLock uint32
   const maxBackoff = 16  // 最大退避次数

   func (sl *spinLock) Lock() {
       backoff := 1
       for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
           // 调用 runtime.Gosched() 主动让出 CPU，避免忙等待浪费资源
           for i := 0; i < backoff; i++ {
               runtime.Gosched()
           }
           if backoff < maxBackoff {
               backoff <<= 1
           }
       }
   }

   func (sl *spinLock) Unlock() {
       atomic.StoreUint32((*uint32)(sl), 0)
   }
   ```

2. sync.Cond
   ![alt text](image-1.png)

   ```go
   type Cond struct {
       noCopy noCopy

       // L is held while observing or changing the condition
       L Locker

       notify  notifyList
       checker copyChecker
   }
   ```

   • 成员变量 noCopy + checker 是一套组合拳，保证 Cond 在第一次使用后不允许被复制；

   • 核心变量 L，一把锁，用于实现阻塞操作；

   • 核心变量 notify，阻塞链表，分别存储了调用 Cond.Wait() 方法的次数、goroutine 被唤醒的次数、一把系统运行时的互斥锁以及链表的头尾节点.

   ```go
   type notifyList struct {
       wait   uint32
       notify uint32
       lock   uintptr // key field of the mutex
       head   unsafe.Pointer
       tail   unsafe.Pointer
   }
   ```

   - Cond.Wait
   - Cond.Signal
   - Cond.Broadcast
