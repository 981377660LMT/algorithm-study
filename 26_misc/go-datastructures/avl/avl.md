下面给出对这段“可持久化（不可变）AVL树”示例代码的详细分析，并提供相应的使用方法示例，帮助理解其实现原理与实际用法。

---

## 一、核心特性

1. **不可变（Immutable）**  
   该 AVL 树的插入、删除操作并不会就地修改原有数据结构，而是通过“分支复制”（branch copy）生成一棵新树。对外暴露的方法 `Insert`、`Delete`、`Get` 都返回新的版本或结果，因此原先版本始终保持不变。  
   这是一种 **“持久化（Persistent）”** 或 **“不可变（Immutable）”** 数据结构：在并发场景中，可以提供更高的读取一致性与安全性。

2. **AVL 树**

   - AVL 是平衡二叉搜索树的一种，能保证插入、删除、查找均为 \(O(\log n)\)。
   - 每个节点维护一个 `balance` 值（高为 +1 或 -1），插入和删除时根据平衡因子做相应的旋转或双旋转。

3. **可扩展性**  
   结合其他结构（如 B+ 树）可用于构建不可变索引、不可变映射等。

---

## 二、代码结构与主要类型

### 1. \`Immutable\` 结构

```go
// Immutable represents an immutable AVL tree.
type Immutable struct {
    root   *node   // 根节点
    number uint64  // 节点总数
    dummy  node    // 用于插入操作的辅助节点
}
```

- `root`：AVL 树的根节点指针，类型为 `*node`。
- `number`：记录 AVL 树中元素（节点）数量。
- `dummy`：在插入操作时会被用作临时的辅助节点，避免在原树上就地修改。

### 2. \`node\` 结构

```go
type node struct {
    balance  int8      // 平衡因子，|balance| <= 1
    children [2]*node  // 左右子节点，children[0]是左，children[1]是右
    entry    Entry     // 节点存储的数据
}
```

- `balance`：AVL 树维护的平衡因子，取值 -1、0、1 等；当插入或删除导致失衡时，需要旋转或双旋转恢复平衡。
- `children`：左右子节点。
- `entry`：实现了 `Entry` 接口的数据对象（可比较），在树中保存。

### 3. \`Entry\` 接口

```go
// Entry represents all items that can be placed into the AVL tree.
type Entry interface {
    // Compare should return -1, 0, or 1
    //  -1 表示 this < other
    //   0 表示 this = other
    //   1 表示 this > other
    Compare(Entry) int
}
```

- 使用者需要实现此接口，才能将自定义类型存入 AVL 树。例如整型、字符串或自定义结构都可以，只要实现 `Compare`。

### 4. 其他类型

- `Entries []Entry`：批量操作时，用切片装载多个 `Entry`。
- `nodes []*node`：用于在删除操作中作为栈来暂存节点。
- 若干帮助函数（`rotate`, `doubleRotate`, `adjustBalance`, `normalizeComparison` 等）实现旋转、平衡调整等细节。

---

## 三、主要方法解析

### 1. `NewImmutableAVL`

```go
func NewImmutableAVL() *Immutable {
    immutable := &Immutable{}
    immutable.init() // 初始化 dummy 节点等
    return immutable
}
```

- 分配并返回一个新的不可变 AVL 树对象。
- 内部调用 `immutable.init()` 来初始化 `dummy` 节点（辅助插入时使用）。

### 2. `Get` / `get`

```go
func (immutable *Immutable) Get(entries ...Entry) Entries {
    result := make(Entries, 0, len(entries))
    for _, e := range entries {
        result = append(result, immutable.get(e))
    }
    return result
}

func (immutable *Immutable) get(entry Entry) Entry {
    n := immutable.root
    for n != nil {
        cmpResult := n.entry.Compare(entry)
        switch {
        case cmpResult == 0:
            return n.entry
        case cmpResult > 0:
            n = n.children[0] // go left
        case cmpResult < 0:
            n = n.children[1] // go right
        }
    }
    return nil
}
```

- 接受一个或多个要查找的 `Entry`，在树中依次寻找，返回对应的元素（或 `nil`）。
- `get` 按照二叉搜索树的规则遍历：如果当前节点大于查找目标，就向左子树走；小于查找目标，就向右子树走；相等则返回。

### 3. `Insert` / `insert`

```go
func (immutable *Immutable) Insert(entries ...Entry) (*Immutable, Entries) {
    if len(entries) == 0 {
        return immutable, Entries{}
    }

    overwritten := make(Entries, 0, len(entries))
    cp := immutable.copy()   // 分支复制
    for _, e := range entries {
        overwritten = append(overwritten, cp.insert(e))
    }

    return cp, overwritten
}
```

- 对外的批量插入接口：会先对当前的 `immutable` 做一次 `copy()` 生成新树 `cp`，然后把 `entries` 逐个插入到 `cp` 中。
- 最后返回新版本的树 `cp` 以及被覆盖的旧值（如果某个插入条目已存在，会把旧条目以 `Entry` 形式返回，否则返回 `nil`）。

#### `immutable.copy()`

```go
func (immutable *Immutable) copy() *Immutable {
    var root *node
    if immutable.root != nil {
        root = immutable.root.copy()
    }
    cp := &Immutable{
        root:   root,
        number: immutable.number,
        dummy:  *newNode(nil),
    }
    return cp
}
```

- 这里的“拷贝”只是复制根节点与 `number` 计数，还有 `dummy` 重置。
- 在后续真正往下继续插入时，还会在需要时对沿途节点做“分支复制”，从而保证原树不变。

#### `insert` 方法核心

```go
func (immutable *Immutable) insert(entry Entry) Entry {
    if immutable.root == nil {
        immutable.root = newNode(entry)
        immutable.number++
        return nil
    }

    // 重置 dummy，作为插入过程中的辅助节点
    immutable.resetDummy()

    // 1. 找插入位置：自顶向下，比较 entry 大小，选择走左/右子树
    // 2. 在走的过程中做 copy()，避免修改原节点
    // 3. 找到空位就插入一个新节点 newNode(entry)
    // 4. 在回溯时更新平衡因子 balance，如果出现 |balance| > 1，做对应旋转

    // ...
    // 省略若干行，见源码
    // ...

    return nil
}
```

- 整体遵循 AVL 树插入步骤，只是每次需要修改节点时会先 `copy()` 一份，避免影响原节点。
- 插入完成后，如出现失衡（|balance| > 1）就通过 `insertBalance` 做旋转或双旋转。
- `balance` 的维护方式：
  - 如果向左子树插入，当前节点 `balance--`；
  - 如果向右子树插入，当前节点 `balance++`。

### 4. `Delete` / `delete`

```go
func (immutable *Immutable) Delete(entries ...Entry) (*Immutable, Entries) {
    if len(entries) == 0 {
        return immutable, Entries{}
    }

    deleted := make(Entries, 0, len(entries))
    cp := immutable.copy()
    for _, e := range entries {
        deleted = append(deleted, cp.delete(e))
    }

    return cp, deleted
}
```

- 与 `Insert` 类似，会先复制当前树，再依次删除指定的条目，返回新树和被删除的旧值。

#### `delete` 方法

```go
func (immutable *Immutable) delete(entry Entry) Entry {
    if immutable.root == nil {
        return nil
    }

    // 1. 找到要删除的节点，用一个栈 (cache + dirs) 记录沿途经过的节点与方向
    // 2. 确认找到后，再做 branch copy（复制沿途节点）
    // 3. 根据删除情况 (叶子节点、有一个子节点、或有两个子节点) 做不同处理
    // 4. 重新计算 balance 并执行 removeBalance
    // ...
    // 返回被删除节点的 oldEntry
}
```

- 这是 AVL 树标准删除流程：
  1. 如果要删的是叶节点，直接将父节点指向空。
  2. 如果只有一个子节点，就让父节点指向这个子节点。
  3. 如果有左右子节点，则先找到后继或前驱（源码里是往右子树中寻找最左孩子），替换当前节点内容，然后再把后继节点删除。
- 最终通过 `removeBalance` 调整平衡。

### 5. AVL 平衡调整函数

- `insertBalance`：插入后的平衡恢复。
- `removeBalance`：删除后的平衡恢复。
- `rotate`、`doubleRotate`：单旋转、双旋转。
- `adjustBalance`：旋转前后微调子节点的平衡因子。

它们的实现都是标准的 AVL 树操作，只是要注意 **这里会调用节点的 `copy`**，以保证原节点不被就地修改。

---

## 四、代码使用方法

### 1. 实现 `Entry` 接口

假设我们要在这棵 AVL 树中存放整型 `intEntry`，可以这样实现：

```go
type intEntry int

func (i intEntry) Compare(other avl.Entry) int {
    o := other.(intEntry)
    switch {
    case i < o:
        return -1
    case i > o:
        return 1
    default:
        return 0
    }
}
```

### 2. 创建并操作树

下面是一个简单示例，演示如何创建树、插入、查找、删除等操作：

```go
package main

import (
    "fmt"
    "github.com/xxx/avl"  // 假设该包路径为 github.com/xxx/avl
)

func main() {
    // 1. 创建不可变AVL树
    tree := avl.NewImmutableAVL()

    // 2. 插入
    // Insert返回 (新树, 被覆盖的旧值列表)
    nums := []avl.Entry{intEntry(10), intEntry(5), intEntry(15), intEntry(10)}
    newTree, overwritten := tree.Insert(nums...)
    // 因为 10 重复插入一次，所以可能返回被覆盖的旧值
    fmt.Printf("overwritten: %v\n", overwritten)
    // oldTree 保持不变，newTree 是新的版本

    // 3. 查找
    // Get返回[]Entry，如果找不到则返回nil作为对应下标
    found := newTree.Get(intEntry(5), intEntry(10), intEntry(999))
    fmt.Printf("found: %v\n", found)
    // found[2] 可能就是 nil，因为 999 不存在

    // 4. 删除
    // Delete返回 (新树, 被删除的旧值列表)，无法删除就返回nil
    newTree2, deleted := newTree.Delete(intEntry(10))
    fmt.Printf("deleted: %v\n", deleted)

    // 5. 再查一次看看
    found2 := newTree2.Get(intEntry(10))
    fmt.Printf("found2 after delete: %v\n", found2) // 应该是 [nil]

    // 注意：tree/newTree/newTree2 是三个版本的树，并且可以并行安全地使用它们
}
```

### 3. 关于不可变（持久化）的好处

- **多版本并行**：插入、删除返回新版本树，而老版本仍然可以使用。
- **读操作安全**：在并发场景下，只要不进行破坏性写操作，就无需额外加锁；老版本可以被任意线程安全地读取。
- **注意内存消耗**：为了保持历史版本，在频繁插入或删除时，会有一定的分支复制开销。但由于只在“沿途路径”上进行复制，实际复制量一般是 \(O(\log n)\)，并不需要整树复制。

---

## 五、总结

1. **不可变 AVL 树**：在每次插入、删除时会创建并返回一棵新树，保持原树不变。
2. **关键实现**：通过“分支复制”技术，仅复制插入/删除路径上的节点来保持持久化。
3. **标准 AVL 平衡**：平衡调整（插入、删除后的旋转）完全符合 AVL 逻辑。
4. **使用接口**：提供 `Insert`、`Delete`、`Get`、`Len` 等常用操作，需自定义 `Entry.Compare` 来确定排序顺序。
5. **使用场景**：适合需要多版本并行、提高读并发安全性、或实现历史回溯的场合。

示例代码中展示的用法，涵盖了大部分常用操作和特性。在实际应用里，根据需求自己实现 `Entry` 类型的 `Compare` 方法即可投入使用。
