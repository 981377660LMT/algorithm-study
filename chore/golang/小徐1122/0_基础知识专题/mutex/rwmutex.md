以下是`sync.RWMutex`源代码的详细解读，分为核心结构、方法实现和设计思想三个部分：

---

### **一、核心结构**

```go
type RWMutex struct {
    w           Mutex         // 互斥锁，用于写操作的独占控制
    writerSem   uint32        // 写者信号量，写者等待读者完成
    readerSem   uint32        // 读者信号量，读者等待写者完成
    readerCount atomic.Int32  // 读者计数器（正：活跃读者数；负：有写者等待）
    readerWait  atomic.Int32  // 写者需要等待的剩余读者数
}
```

- **`w`**: 底层互斥锁，保证写操作的原子性。
- **`readerCount`**:
  - 正值：当前持有读锁的读者数量。
  - 负值：表示有写者持有或等待锁（通过减去`rwmutexMaxReaders`实现）。
- **`readerWait`**: 写者等待时，剩余需要释放的读者数量。
- **信号量**:
  - `writerSem`: 写者等待读者释放锁时休眠在此信号量。
  - `readerSem`: 读者等待写者释放锁时休眠在此信号量。

---

### **二、关键方法解析**

#### **1. RLock()：获取读锁**

```go
func (rw *RWMutex) RLock() {
    if rw.readerCount.Add(1) < 0 {
        runtime_SemacquireRWMutexR(&rw.readerSem) // 等待写者释放
    }
}
```

- **快速路径**：通过原子操作增加`readerCount`，若结果 ≥0，说明无写者等待，直接获得读锁。
- **慢速路径**：若`readerCount`变为负值（有写者等待），当前读者休眠在`readerSem`上，等待写者完成。

---

#### **2. RUnlock()：释放读锁**

```go
func (rw *RWMutex) RUnlock() {
    if r := rw.readerCount.Add(-1); r < 0 {
        rw.rUnlockSlow(r) // 处理写者等待的情况
    }
}

func (rw *RWMutex) rUnlockSlow(r int32) {
    if rw.readerWait.Add(-1) == 0 {
        runtime_Semrelease(&rw.writerSem) // 最后一个读者唤醒写者
    }
}
```

- 减少`readerCount`，若结果<0，说明有写者正在等待。
- `readerWait`递减，当减为 0 时（所有活跃读者已释放），通过`writerSem`唤醒写者。

---

#### **3. Lock()：获取写锁**

```go
func (rw *RWMutex) Lock() {
    rw.w.Lock() // 通过底层Mutex竞争写锁
    r := rw.readerCount.Add(-rwmutexMaxReaders) + rwmutexMaxReaders
    if r != 0 && rw.readerWait.Add(r) != 0 {
        runtime_SemacquireRWMutex(&rw.writerSem) // 等待现有读者完成
    }
}
```

1. 先通过`w.Lock()`阻止其他写者。
2. 将`readerCount`减去`rwmutexMaxReaders`，使其变为负值，阻止新读者。
3. 计算当前活跃读者数`r`，若存在活跃读者，写者休眠在`writerSem`上，直到`readerWait`归零。

---

#### **4. Unlock()：释放写锁**

```go
func (rw *RWMutex) Unlock() {
    r := rw.readerCount.Add(rwmutexMaxReaders) // 恢复readerCount为正
    for i := 0; i < int(r); i++ {
        runtime_Semrelease(&rw.readerSem) // 唤醒所有等待的读者
    }
    rw.w.Unlock() // 释放底层Mutex，允许其他写者竞争
}
```

1. 恢复`readerCount`为正，允许新读者获取锁。
2. 循环唤醒所有在`readerSem`上等待的读者。
3. 释放底层互斥锁`w`，允许其他写者竞争。

---

#### **5. TryRLock() 与 TryLock()**

- **非阻塞尝试**：通过原子操作（`CompareAndSwap`）检查当前状态，若条件允许立即获取锁，否则返回失败。
- **设计警示**：文档指出这些方法应谨慎使用，可能暗示并发设计问题。

---

### **三、设计思想**

1. **写者优先**:

   - 写者通过将`readerCount`置为负，阻止新读者，确保现有读者释放后，写者能及时获取锁，避免写者饥饿。

2. **原子操作与信号量**:

   - 使用`atomic`包保证计数器操作的原子性。
   - 通过信号量（`runtime_Semacquire`/`runtime_Semrelease`）实现高效的阻塞与唤醒。

3. **竞态检测优化**:

   - 通过`race.Enable()`和`race.Disable()`在关键路径关闭竞态检测，避免原子操作干扰分析器的正确性判断。

4. **递归锁限制**:
   - 明确禁止递归读锁，防止因写者等待导致死锁（如读锁未释放前尝试再次获取读锁）。

---

### **四、关键常量与辅助函数**

- **`rwmutexMaxReaders`**: 最大读者数（1<<30），通过`readerCount`的符号区分读者与写者。
- **`syscall_hasWaitingReaders`**: 内部接口，供`syscall`包检查是否有等待的读者。

---

### **五、总结**

`sync.RWMutex`通过精细的原子操作和信号量机制，在保证并发安全的前提下，实现了高效的读写锁：

- **读者间无竞争**：可并发读取。
- **写者独占**：通过底层 Mutex 和计数器状态确保原子性。
- **优先级平衡**：写者等待期间阻止新读者，避免饥饿。

代码中体现了 Go 语言对性能与正确性的极致追求，是并发原语设计的经典范例。
