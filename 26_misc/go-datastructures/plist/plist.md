下面给出对这段 **不可变（持久化）链表**实现的详细分析，并附带一个简单的使用示例，帮助理解其内部逻辑与用法。

---

## 一、不可变（持久化）链表简介

- **不可变（Persistent）链表**：一旦创建以后就不会被修改，每次对链表的操作（插入、删除等）都会返回一个新的链表（共享原链表的不可变部分）。
- **应用场景**：在多线程环境下或需要保存历史版本（回溯）时，非常有用。可保证读线程看到的链表不会被其他写线程破坏。

在 Go 中，如果只用传统的切片或 `container/list` 进行插入、删除，往往是就地修改。但这个“PersistentList” 则始终保持过去版本的完整性。

---

## 二、主要类型与接口

### 1. `PersistentList` 接口

```go
type PersistentList interface {
    Head() (interface{}, bool)
    Tail() (PersistentList, bool)
    IsEmpty() bool
    Length() uint
    Add(head interface{}) PersistentList
    Insert(val interface{}, pos uint) (PersistentList, error)
    Get(pos uint) (interface{}, bool)
    Remove(pos uint) (PersistentList, error)
    Find(func(interface{}) bool) (interface{}, bool)
    FindIndex(func(interface{}) bool) int
    Map(func(interface{}) interface{}) []interface{}
}
```

- 该接口定义了不可变链表需要支持的基本操作，比如 `Head`, `Tail`, `Add`, `Insert`, `Remove` 等。

### 2. `emptyList` 与 `list`

源码里共有两个具体类型实现了 `PersistentList`：

1. **`emptyList`**：表示一个空链表。所有操作要么返回默认值（如 `nil`、`false` 等），要么构造一个新的非空链表实例。
2. **`list`**：表示一个非空的节点，包含：
   ```go
   type list struct {
       head interface{}
       tail PersistentList
   }
   ```
   - `head`：当前节点数据。
   - `tail`：指向后面剩余节点的 `PersistentList`（可共享结构）。

#### 为什么要有 `emptyList` 与 `list` 两种结构？

- **空链表**可以复用一个全局单例 `Empty`，让你在需要返回空时就直接使用这个实例。
- **非空链表**是通过“head + tail”来递归定义的。如果 `tail` 也是空，就到了列表的末尾；如果 `tail` 非空，则继续。

在代码里：

```go
var (
    Empty PersistentList = &emptyList{}
    ErrEmptyList         = errors.New("Empty list")
)
```

- `Empty` 是一个全局空链表。无需重复分配，可被多个地方共享。

---

## 三、核心方法解析

以下展示主要方法在 **空链表** (`emptyList`) 与 **非空链表** (`list`) 两种情况下的区别。

### 1. `IsEmpty()` / `Length()`

- 在 `emptyList`：
  ```go
  func (e *emptyList) IsEmpty() bool { return true }
  func (e *emptyList) Length() uint  { return 0 }
  ```
- 在 `list`：
  ```go
  func (l *list) IsEmpty() bool { return false }
  func (l *list) Length() uint {
      curr := l
      length := uint(0)
      for {
          length += 1
          tail, _ := curr.Tail()
          if tail.IsEmpty() {
              return length
          }
          curr = tail.(*list)
      }
  }
  ```
  - 对于非空链表，`IsEmpty()` = false。
  - `Length()` 遍历整个链表，累加节点数。

### 2. `Head()` / `Tail()`

- 在 `emptyList`：

  ```go
  func (e *emptyList) Head() (interface{}, bool) {
      return nil, false
  }

  func (e *emptyList) Tail() (PersistentList, bool) {
      return nil, false
  }
  ```

  - 空链表无法提供任何“头节点”或“尾链表”，所以返回 `(nil, false)` 表示无效。

- 在 `list`：

  ```go
  func (l *list) Head() (interface{}, bool) {
      return l.head, true
  }

  func (l *list) Tail() (PersistentList, bool) {
      return l.tail, true
  }
  ```

  - 非空链表直接返回 `head` 与 `tail`，同时给出 `true` 表示数据有效。

### 3. `Add()` —— 向前端插入

- 在 `emptyList`：

  ```go
  func (e *emptyList) Add(head interface{}) PersistentList {
      return &list{head, e}
  }
  ```

  - 给空链表增加一个元素，结果就是一个 **只包含这个元素，tail = emptyList** 的新链表。

- 在 `list`：
  ```go
  func (l *list) Add(head interface{}) PersistentList {
      return &list{head, l}
  }
  ```
  - 在现有链表前面插入一个新的节点：新的 `head` = 给定 `head`，而它的 `tail` = 旧的列表 `l`。
  - **关键：**旧链表 `l` 不变；我们创建一个新的 `list` 节点，共享旧链表作为其尾部。

### 4. `Insert(val, pos)` —— 在指定位置插入

- 空链表版本：

  ```go
  func (e *emptyList) Insert(val interface{}, pos uint) (PersistentList, error) {
      if pos == 0 {
          return e.Add(val), nil
      }
      return nil, ErrEmptyList
  }
  ```

  - 如果 `pos=0`，就相当于在空表前面插入一个节点；否则报错。

- 非空链表版本：
  ```go
  func (l *list) Insert(val interface{}, pos uint) (PersistentList, error) {
      if pos == 0 {
          return l.Add(val), nil
      }
      nl, err := l.tail.Insert(val, pos-1)
      if err != nil {
          return nil, err
      }
      return nl.Add(l.head), nil
  }
  ```
  - 若 `pos=0`，相当于“前面插入”，跟 `Add` 一样。
  - 否则递归调用 `l.tail.Insert(val, pos-1)`，成功后再把当前 `head` 加回去。
  - 因此，这种插入是**不可变**的：每一层都返回一个新的节点，并共享剩余结构不变。

### 5. `Get(pos)` —— 获取指定位置元素

- 空链表：直接返回 `(nil, false)` 表示无效。
- 非空链表：
  ```go
  func (l *list) Get(pos uint) (interface{}, bool) {
      if pos == 0 {
          return l.head, true
      }
      return l.tail.Get(pos - 1)
  }
  ```
  - 如果 `pos=0` 就是当前节点；否则递归到 `tail`。

### 6. `Remove(pos)` —— 移除指定位置节点

- 空链表：直接报错 `ErrEmptyList`。
- 非空链表：
  ```go
  func (l *list) Remove(pos uint) (PersistentList, error) {
      if pos == 0 {
          nl, _ := l.Tail()
          return nl, nil
      }
      nl, err := l.tail.Remove(pos - 1)
      if err != nil {
          return nil, err
      }
      return &list{l.head, nl}, nil
  }
  ```
  - 若 `pos=0` 就把当前节点去掉，返回 `tail`；否则递归到 tail，删除其 `pos-1` 位置，再构建一个新头节点。

### 7. `Find(pred)` / `FindIndex(pred)`

- 逐个匹配 `head`，若 `pred(head)` 返回 `true`，就返回当前节点或当前索引，否则继续往下查。
- **空链表**里则直接返回 `(nil, false)` 或 -1。

### 8. `Map(f)` —— 对每个元素应用函数，返回一个 `[]interface{}`

- 空链表返回 nil。
- 非空链表：
  ```go
  func (l *list) Map(f func(interface{}) interface{}) []interface{} {
      return append(l.tail.Map(f), f(l.head))
  }
  ```
  - 先对 tail 进行 `Map`（递归），得到一个切片，然后把当前 `head` 的映射值 `f(l.head)` 追加上去。

---

## 四、使用示例

以下展示如何使用这个不可变链表进行常见操作。假设代码在 `list` 包中，我们在一个示例 `main.go` 文件中导入并使用：

```go
package main

import (
    "fmt"
    "github.com/xxx/list" // 假设这个不可变链表放在此路径
)

func main() {
    // 1) 从空链表开始
    pl := list.Empty

    // 2) Add
    pl = pl.Add("A") // 头插入 "A"
    pl = pl.Add("B") // 头插入 "B"
    pl = pl.Add("C") // 头插入 "C"

    // 链表现在顺序: [C, B, A]
    headVal, _ := pl.Head()
    fmt.Println("Head:", headVal) // C

    // 3) Insert
    // 在位置1插入 "D" => 位置0是 "C", 位置1是 "B"
    pl, err := pl.Insert("D", 1)
    if err != nil {
        panic(err)
    }
    // 新链表: [C, D, B, A]

    // 4) Get
    val, ok := pl.Get(2)
    fmt.Println("Index2:", val, "OK:", ok) // Index2: B OK: true

    // 5) Remove
    pl2, err := pl.Remove(1) // 移除位置1 => "D"
    // pl2: [C, B, A]
    headVal2, _ := pl2.Head()
    fmt.Println("After remove, head:", headVal2) // C

    // pl 仍然是 [C, D, B, A], pl2 是 [C, B, A]
    // 不可变 => 原来的列表 pl 未改变

    // 6) Map
    mapped := pl.Map(func(x interface{}) interface{} {
        return fmt.Sprintf("mapped_%v", x)
    })
    fmt.Println("Mapped slice:", mapped)
    // => ["mapped_A", "mapped_B", "mapped_D", "mapped_C"] (顺序根据内部Append)

    // 7) Find / FindIndex
    foundVal, foundOk := pl.Find(func(x interface{}) bool {
        return x == "B"
    })
    idx := pl.FindIndex(func(x interface{}) bool {
        return x == "B"
    })
    fmt.Println("Found B:", foundVal, foundOk, "Index:", idx)

    // 8) IsEmpty, Length
    fmt.Println("IsEmpty:", pl.IsEmpty()) // false
    fmt.Println("Length:", pl.Length())   // 4
}
```

运行结果示例（注意 Map 的结果顺序是逆序追加）：

```
Head: C
Index2: B OK: true
After remove, head: C
Mapped slice: [mapped_A mapped_B mapped_D mapped_C]
Found B: B true Index: 2
IsEmpty: false
Length: 4
```

---

## 五、实现细节与特性

1. **不可变**：任何插入、删除操作都不会修改已有链表，而是返回一个指向新“头节点”的链表对象。旧链表的节点结构仍然可访问，不会被破坏。
2. **共享结构**：当我们 `Add` 或 `Insert` 时，新链表的 `tail` 通常直接指向旧链表，形成结构共享，节省内存。
3. **性能**：
   - 单次 `Add`（头部插入）为 \( O(1) \)。
   - `Insert(pos)`、`Remove(pos)`、`Get(pos)`、`Length()` 在最坏情况下为 \( O(n) \)（需要遍历），因为是单向链表。
   - 若需要快速随机访问，大量插入/删除中间元素，可能性能不佳。但若用于函数式风格的“只做头部操作”，性能可接受。
4. **空间使用**：不可变链表每次操作都分配一个新节点，但底层的 `tail` 可能与之前版本共享，不会完全拷贝。
5. **线程安全**：因其不可变性，任意协程都可以安全地并发读旧版本的列表，写操作产生新列表，不会破坏旧版本。

---

## 六、总结

- **PersistentList** 提供了一种函数式、不可变的数据结构，以“空节点 + (head, tail)”递归定义。
- 通过简单的递归实现各操作，保证旧版本链表不会被修改，形成多版本共享。
- 适用于不经常在中间位置修改，但需要多版本并发读取或保留历史记录的场合。
- 使用示例展示了常用操作：Add、Insert、Remove、Get、Map 等，理解了其不变性与共享结构的特性。

这段代码可当作一个基础“函数式链表”工具，在实际生产中，如果数据量大、需要大量中间修改或随机访问，也许需要更复杂的结构（如可持久化平衡树、指针+skip list 等），但在小规模或只做头部操作的场景中非常方便。
