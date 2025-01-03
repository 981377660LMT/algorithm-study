下面给出对这段 **Dtrie**（动态、持久化哈希前缀树）实现的详细分析，并附带简要用例，帮助理解其原理与使用方式。

---

## 一、Dtrie 简介

- **Dtrie**：一种持久化（不可变）哈希前缀树（或称**哈希数组映射前缀树**，HAMT），支持动态扩展或收缩，以达到在不同维度上高效的内存使用。
- **持久化**：插入、删除操作不在原结构上就地修改，而是返回一个新的 `Dtrie` 实例，允许保留旧版本（历史快照）。
- **哈希前缀**：对 key 做哈希后，每次解析 `w=5` bit（本实现中 5 位）决定树节点的索引，逐层往下。最多 7 层 (`level` 0 到 6) 即可处理 32 位哈希（因为 \(6 \times 5 = 30\) bit，再加上第 7 层即余下 2 bit，合理地处理 32 位范围）。
- **冲突处理**：若哈希“过深”或出现同哈希而 key 不同的情况，会转换为一个“collisionNode”存放多个值。

---

## 二、主要结构解析

### 1. `Dtrie` 结构

```go
type Dtrie struct {
    root   *node
    hasher func(v interface{}) uint32
}
```

- **root**：指向一棵前缀树的根节点（类型 `*node`）。
- **hasher**：对 key 进行哈希的函数，若为 `nil` 则使用默认的 `defaultHasher`。

### 2. `node`

```go
type node struct {
    entries []Entry
    nodeMap Bitmap32
    dataMap Bitmap32
    level   uint8
}
```

- `entries`：大小可变的 slice，包含 32 或者 4 个 `Entry`（视 `level` 而定，详见后面逻辑）。
- `nodeMap` 与 `dataMap`：分别用位图 (bitmap) 标记“哪几个下标是子节点（sub-node）”和“哪几个下标是直接存储 Entry”。这样可以紧凑存储不同分支，而不必在 `entries` 中浪费太多空间。
  - 若 `dataMap.GetBit(i)` 为 `true`，表示 `entries[i]` 存的是一个单键的 `Entry`（或空槽等）。
  - 若 `nodeMap.GetBit(i)` 为 `true`，表示 `entries[i]` 存的是一个更深层的子树（`*node`）。
- `level`：表示在哈希前缀的哪一层（0~6）。

### 3. `collisionNode`

```go
type collisionNode struct {
    entries []Entry
}
```

- 用于处理“冲突”：在最深层 (level=6) 且无法进一步区分哈希的时候，如果 key 不同，就把它们都放在一个 `collisionNode` 里。

### 4. `Entry`

```go
type Entry interface {
    KeyHash() uint32
    Key() interface{}
    Value() interface{}
}
```

- 代表 Trie 存放的键值对或节点。
- 具体实现有：
  - `*entry{hash, key, value}`：单个键值对
  - `*node` / `*collisionNode` 也实现了该接口，部分方法返回空或不参与普通查找。

---

## 三、主要操作

### 1. 构造

```go
func New(hasher func(v interface{}) uint32) *Dtrie {
    if hasher == nil {
        hasher = defaultHasher
    }
    return &Dtrie{
        root:   emptyNode(0, 32),
        hasher: hasher,
    }
}
```

- 初始化一个空的根节点 `emptyNode(0, 32)` 表示 `level=0`、容量 32（因为初始层可容纳 32 路分支）。
- `hasher` 为 nil 则用 `defaultHasher`（基于 `fnv.New32a()` 或对整数类型做简单转换）。

### 2. 插入（Insert）

```go
func (d *Dtrie) Insert(key, value interface{}) *Dtrie {
    root := insert(d.root, &entry{d.hasher(key), key, value})
    return &Dtrie{root, d.hasher}
}
```

- 返回一个新的 `Dtrie` 实例（**持久化**），原 `Dtrie` 不变。
- `insert(n *node, entry Entry) *node` 具体逻辑：
  1. 计算下标 `index := mask(entry.KeyHash(), n.level)`，mask取 `(hash >> (5*level)) & 0x1f`。
  2. 根据 `dataMap`、`nodeMap` 判断该下标当前是空、存了单个条目、或存了子节点：
     - **若空**：把 `entry` 塞进去，`dataMap.SetBit(index)`。
     - **若是子节点** (`nodeMap.GetBit(index)`)：递归 `insert` 到子节点。
     - **若已有单条目**：
       - 如果 key 相同，则替换；
       - 否则要新建一个子节点，里面放原有条目和新条目（相当于把二者分到下一层），并更新 `nodeMap` / `dataMap`。
  3. 在 `level=6` (最深层)时，如果已有碰撞，则用/或创建 `collisionNode` 把多条目集合起来。

### 3. 查找（Get）

```go
func (d *Dtrie) Get(key interface{}) interface{} {
    node := get(d.root, d.hasher(key), key)
    if node != nil {
        return node.Value()
    }
    return nil
}
```

- 内部 `get(n *node, keyHash uint32, key interface{}) Entry`：
  1. 根据 level，计算 `index = mask(keyHash, n.level)`；
  2. 若 `dataMap.GetBit(index)` 为 true，就直接是单条目；
  3. 若 `nodeMap.GetBit(index)` 为 true，则说明是子树，递归下去；
  4. 若 `level=6` 且 `entries[index] != nil` 是 `collisionNode`，则线性扫描找匹配 key。
  5. 否则没找到，返回 nil。

### 4. 删除（Remove）

```go
func (d *Dtrie) Remove(key interface{}) *Dtrie {
    root := remove(d.root, d.hasher(key), key)
    return &Dtrie{root, d.hasher}
}
```

- 内部 `remove(n *node, keyHash uint32, key interface{}) *node`：
  1. 同样先取 index；
  2. 若 `dataMap.GetBit(index)` 为 true，说明存的是单条目，清空该槽即可；
  3. 若 `nodeMap.GetBit(index)` 为 true，则向下递归删除；
  - 如果子节点只剩 1 个 entry，就把它压缩（pull up）到当前节点，避免冗余子节点。
  4. 若是 `level=6` 中的 collisionNode，直接在 `entries` 里删除对应 key；若只剩 1 条，则拉升为单条。

> 这些操作也返回新的根节点。

### 5. 迭代（Iterator）

```go
func (d *Dtrie) Iterator(stop <-chan struct{}) <-chan Entry {
    return iterate(d.root, stop)
}
```

- 返回一个只读 channel，遍历 trie 中所有 `Entry`。
- `iterate` 函数开启一个 goroutine，调用 `pushEntries(n, stop, out)` 递归发送 `Entry` 到 channel：
  - 遍历 `entries`，若 `dataMap.GetBit(i)` 则是单条目；若 `nodeMap.GetBit(i)` 则递归子节点；若 `level=6` + collisionNode 则发送所有冲突条目。
  - 若 `stop` 通道被关闭，则提前返回，以避免 goroutine 泄漏（一直卡在发送端）。

### 6. 其它操作

- `Size()`：简单地遍历所有 entry 计数。注意大数据量时会 O(n)。
- `mask(hash, level)`: 提取 “hash 的第 [5*level .. 5*(level+1)-1] 位”。
- **默认哈希**：对常见数字类型做简单转换；否则使用 `fnv.New32a()` 对 `fmt.Sprintf("%#v", value)` 进行哈希。

---

## 四、位图（Bitmap）与紧凑存储

```go
type Bitmap32 uint32

func (b Bitmap32) SetBit(pos uint) Bitmap32    { ... }
func (b Bitmap32) ClearBit(pos uint) Bitmap32  { ... }
func (b Bitmap32) GetBit(pos uint) bool        { ... }
func (b Bitmap32) PopCount() int               { ... }
```

- 在哈希前缀树中，一层最多 32 个分支。若只有少量有效分支，直接用长度为 32 的切片会浪费空间。
- 这里用 `dataMap` 标记“存有单条目的分支索引”，`nodeMap` 标记“指向子节点的分支索引”。
- 根据 bit 的个数决定如何在 `entries` 里排布数据，从而**动态压缩**无效分支。

不过在此实现中，`entries` 的大小是固定 32 或 4（对最深层而言），而不是像正统的 HAMT 那样动态缩放 `entries` 只在 `PopCount()` 范围内使用。这里还是将 `entries` 分配为“可能最大”大小，但用 bitmap 来管理实际有效位置。

---

## 五、使用示例

```go
package main

import (
    "fmt"
    "dtrie"
)

func main() {
    // 1) 构造空 dtrie
    dt := dtrie.New(nil)

    // 2) 插入
    dt2 := dt.Insert("foo", 123)
    dt3 := dt2.Insert("bar", "hello")
    // dt 不变, dt2 不变, dt3 = [("foo", 123), ("bar", "hello")]

    // 3) 获取
    val := dt3.Get("foo")
    fmt.Println("foo ->", val) // foo -> 123

    // 4) 删除
    dt4 := dt3.Remove("foo")
    val2 := dt4.Get("foo")
    fmt.Println("foo ->", val2) // foo -> <nil>, because removed

    // 5) 遍历
    stop := make(chan struct{})
    defer close(stop)
    for e := range dt4.Iterator(stop) {
        fmt.Println(e.Key(), "=", e.Value())
    }
    // bar = hello

    // 6) Size
    fmt.Println("Size of dt4:", dt4.Size()) // 1
}
```

以上示例说明了 Dtrie 的**不可变**属性：每次操作都会返回一个新版本的 Dtrie，旧版本继续保持原数据。

---

## 六、特性与注意事项

1. **不可变**（持久化）：插入、删除均产生新 `Dtrie` 实例。旧版本依然可用。
2. **动态**：在哈希前缀的不同层（0 到 6）使用 `mask` 提取 5 位，用 `Bitmap32` 标记节点/数据存放位置，节省空间。
3. **第 6 层冲突**：当层数达到 6 (处理了 6\*5=30 bit，剩余 2 bit 或仍有 hash 碰撞) 就用 `collisionNode` 存放多个条目。
4. **时间复杂度**：
   - 查找/插入/删除：平均 O(log n) 或 O(1) w.r.t. 哈希分布（因为最大 7 层）。
   - 遍历：需要 DFS 形式，输出所有条目。
5. **内存占用**：此实现在每层都预分配 `entries` 数组大小为 32（或4），可能比完全“HAMT”更浪费些，但也简化了实现。
6. **线程安全**：`Dtrie` 的实例是不可变的，多个 goroutine 并发读同一个实例不会冲突。创建新版本不会影响旧版本，因此写操作也可并发执行（需要小心旧版本的合流等，但本质上不会修改旧版本内存）。

---

## 七、总结

1. **Dtrie** 是对 **Hash Array Mapped Trie** (HAMT) 的一个变形实现：借助 `dataMap` / `nodeMap` 位图来进行紧凑存储，同时保留“持久化/不可变”特性。
2. **操作**：
   - `Insert(key, value)`：返回新 dtrie；若 key 存在则更新；否则在合适节点插入。
   - `Get(key)`：按层级掩码解析 hash，找到对应条目或 nil。
   - `Remove(key)`：在树中移除后返回新 dtrie。
   - `Iterator(stop)`：生成一个 channel，遍历所有条目。
3. **冲突**：在最大层时构建 `collisionNode` 存放多个冲突 key。
4. **不可变**：任何修改均生成新 `Dtrie`，从而可以保留旧版本，实现函数式、并发安全地使用。

这种实现可在多键场景下提供相对高效的查找、插入、删除，且由于是不可变结构，非常利于并发或需要版本化数据的情况。
