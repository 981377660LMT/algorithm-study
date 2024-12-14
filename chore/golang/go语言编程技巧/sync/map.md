下面将详细解析 Go 1.19+ 标准库 `sync` 包中的并发安全 `Map` 实现原理。这段代码在 Go 的 runtime 源码中有微调，但思路一致。这是 Go 中为了解决高并发访问下 **map 读多写少**场景而专门优化的数据结构。

**概述**：  
Go 的 `sync.Map` 是一种线程安全的字典结构，与传统的 `map + sync.Mutex` 或 `map + sync.RWMutex` 不同，它在大量读少量写场景下具有更高的性能。`sync.Map` 内部使用了两份数据结构（read-only map 和 dirty map）以及原子操作和延迟清理等手段减少锁争用。

**核心思想**：

- 使用两个 map 结构：
  1. **read**：只读的 map（`readOnly.m`），多读场景下可以无锁快速访问。如果 key 存在于只读 map 中，即可快速返回结果。
  2. **dirty**：脏 map，用来存放自上次只读快照以来新增或更新的条目。访问 dirty map 时需要上锁。
- 脏标记：当只读 map 不包含目标 key，而 `read.amended` 标识为 true（表示 dirty 中有新键未同步到 read 中），则需要上锁去 dirty map 检查。
- 定期重构：如果锁内访问 dirty map 次数（misses）超过一定阈值，就会将 dirty map 中的数据提升为新的只读 map，从而降低后续访问时的锁争用和 miss 次数。

通过这种机制，`sync.Map` 在频繁读取、很少（或稀少）更新的情况下，读操作几乎无锁直接从只读 map 快速返回；只有当有写操作导致需要更新只读快照时，才会在合适的时机进行 map 的复制或同步。

---

### 主要数据结构

```go
type Map struct {
    mu     Mutex
    read   atomic.Pointer[readOnly] // 原子指针，指向只读结构
    dirty  map[any]*entry
    misses int
}
```

- `mu`：在需要访问 dirty map 或修改全局结构时加锁。
- `read`：原子存储一个 `readOnly` 结构指针，`readOnly` 内包含：

  ```go
  type readOnly struct {
      m       map[any]*entry
      amended bool
  }
  ```

  `amended = false` 表示 `read.m` 已经完整同步了 dirty map 中的键值；为 true 表示 dirty 中有一些 `read.m` 不包含的条目。

- `dirty`：在只读 map 同步之后对新插入或更新的键值存放于 dirty map 中，等待合适时机再同步给 read。
- `misses`：记录访问了 dirty map 的次数，当 `misses` 达到一定阈值后，会将 dirty map 提升为新的只读 map。

```go
type entry struct {
    p atomic.Pointer[any]
}
```

`entry` 是存储单个键值对的结构，其中 `p` 是一个原子指针，指向存储的值。`p` 的值可能为：

- `nil`：表示该键已经被删除
- `expunged`：一个特殊标记，表示该 entry 被标记为已删除并从 dirty 中移除，需要特殊处理才能再次写入。
- 非 nil 且非 `expunged`：正常存储的值指针。

---

### Load(key) 流程

```go
func (m *Map) Load(key any) (value any, ok bool)
```

1. **无锁快读**：先从 `read` (只读 map) 中查找：
   - 如果 `read.m[key]` 存在且未被删除（entry 的 p 不为 nil 或 expunged），直接返回。
2. 如果 `read.m` 中找不到且 `read.amended` 为 false，说明 dirty 中也没有该 key，直接返回不存在。

3. 如果 `read.amended` 为 true，表示 dirty 中可能有该 key。此时：
   - 加锁 (`m.mu.Lock()`)，再次获取最新的 `read`（防止在获取锁期间 `read` 被更新），从头检查。
   - 如果此时仍在只读 map 中找不到，就在 dirty map 中查找 `m.dirty[key]`。
   - 若在 dirty 中也没有该 key，则记录一次 miss (`m.missLocked()`)，因为下次相同 key 的查询会再次需要上锁，超过一定次数就会整合 map。
   - 解锁返回结果。

---

### Store(key, value) 流程

```go
func (m *Map) Store(key, value any)
```

Store 内部调用 `Swap(key, value)`。简化逻辑为：

1. 从只读 map 中查找 key 的 entry。
   - 如果找到且非 expunged，可直接在 entry 中原子交换新值并返回旧值，无需加锁。
   - 如果 entry 是 expunged，需要加锁将其恢复，并加入 dirty。
2. 如果只读 map 中找不到该 key：
   - 加锁检查 dirty map 中是否有该 key。
   - 如果没有且 `read.amended == false`，表示这是插入新 key，需要将只读 map 标记为 amended，并初始化或更新 dirty map 来存储新 key。
3. 在加锁状态下往 dirty 中更新或插入新值。释放锁后完成。

---

### LoadOrStore(key, value)

```go
func (m *Map) LoadOrStore(key, value any) (actual any, loaded bool)
```

此操作要么返回已存在的值，要么存入并返回新值。逻辑类似 Load + Store 的组合，但做了优化：

1. 快速路径：从只读 map 查找，如果找到该 key 的 entry：

   - 尝试原子地加载值，如果该值不为 expunged 和 nil，直接返回已存在的值（loaded=true）。
   - 若为 nil，需要加锁处理（意味着 entry 可能被标记为删除、需要恢复等情况）。

2. 加锁后再检查 read/dirty 决定是更新已有 entry，还是添加新 key 到 dirty map 中。

---

### LoadAndDelete(key)

```go
func (m *Map) LoadAndDelete(key any) (value any, loaded bool)
```

1. 如果只读 map 中找到 key 且有效，尝试删除 entry 的值（将 p 从非 nil 的值原子 CAS 为 nil）。
2. 如果在只读 map 中找不到，就加锁去 dirty 中删除，删除后 miss 计数加一。
3. 返回被删除的旧值和是否存在的标志。

`Delete(key)` 是 `LoadAndDelete(key)` 的简化版本，不返回旧值。

---

### CompareAndSwap(key, old, new) / CompareAndDelete(key, old)

这两个操作在检查当前值是否等于 old 的情况下有条件更新或删除。其流程与前面类似，只是多一步 compare 判断值是否匹配 old。

---

### Range(func(key,value any)bool)

遍历整个 map：

1. 如果 `read.amended == false`，直接遍历 `read.m`。
2. 如果 `read.amended == true`，说明 dirty 中有新条目还未同步到 read。
   - 加锁后将 dirty 提升为新的 read（相当于一次全面同步），再遍历。

遍历过程中，不保证一致性快照（遍历时可能有并发写操作），但保证不会对同一键重复访问。

---

### miss 机制和提升脏表

`m.misses` 记录从只读 map 中找不到键而不得不查询 dirty map 的次数。当 `misses >= len(m.dirty)` 时，意味着脏表访问变频繁，维护现状不如整合。此时将 `m.dirty` 升格为新的只读 map：

- `m.read.Store(&readOnly{m: m.dirty, amended: false})`
- `m.dirty = nil`
- `m.misses = 0`

下一次有新插入键时，会再次创建一个全新的 dirty map。

这样可以达到一个动态平衡：当写入频繁导致大量 miss 时，系统会将 dirty 合并进 read 减少 miss；当之后再次有新增键时，又会创建一个新的 dirty 映射。

---

### Expunged 标记

`expunged` 是一个指针常量，用于标记一个 entry 已被删除并且从 dirty map 中剔除。要在此 entry 重设值之前，需先将它的状态从 `expunged` 恢复（unexpunge），并放回 dirty map 管理。通过 `expunged` 标记可以区分：

- `nil`：entry 已被删除（等待下次 dirty map 创建时转为 expunged）
- `expunged`：entry 确认被彻底删除，没有在 dirty 中管理
- 非空值：有效存储值的正常状态

---

### 性能优化手段小结

- **读多写少场景**：频繁读时，多数读取操作直接从只读 map 中查询，避开 mutex 加锁。
- **写入偶尔发生**：当有写操作时才在必要时对 dirty map 和只读 map 进行更新。
- **延迟同步**：不立即合并，只在 miss 次数达到 dirty 大小时才合并，从而减少拷贝成本。
- **无锁读取通道**：使用 `atomic.Pointer` 实现对 `read` 字段的无锁加载和存储。

---

**总结**：  
`sync.Map` 利用只读 map 与 dirty map 双重结构，结合原子操作和惰性同步策略，让在读多写少场景下读取操作几乎无锁，性能远超简单加锁的 map。写操作虽然仍需加锁，但通过延迟同步、合并 dirty map 来摊薄锁的成本，从整体上提供更好的并发性能。
