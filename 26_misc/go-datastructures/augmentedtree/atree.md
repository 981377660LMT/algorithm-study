# 一、代码整体概述

这个代码主要实现了一个**多维度的区间树**（Interval Tree），支持在**1 个以上的维度**存储区间。其特点包括：

1. **节点结构 `node`**：

   - 每个节点存储一个实现 `Interval` 接口的区间对象，以及节点自身的 `min`、`max`、子节点指针、颜色标记（`red`），以及一个 `id` 用于区分重复区间。
   - 基于某种红黑树变体（可以看到 `red` 节点、旋转操作）来维持平衡。
   - 还维持子树的 `min` / `max` 值，方便快速判断是否可能和查询区间重叠。

2. **插入、删除**：

   - 在插入时，需要做红黑树的插入修正（变色、旋转）；
   - 在删除时，同样执行相应修正并保证平衡。
   - 插入/删除时，会更新节点的 `max` / `min` 以表示子树中的最小和最大边界。

3. **查询**：

   - `Query(interval Interval)`：找出所有与给定区间重叠的节点；
   - 通过节点的 `query` 方法，结合 `overlaps(...)` 进行剪枝。
   - 在多维场景下，还会检查所有维度（这里以 `maxDimension` 为参考）来判定是否重叠。

4. **多维度**：

   - 其实重点在 `intervalOverlaps(...)` 中，会根据 `maxDimension` 调用 `OverlapsAtDimension` 来判断其他维度是否也重叠；第一维的重叠先简单判定，然后若需要再检查其余维度。

5. **接口与使用**：
   - 外部通过 `New(dimensions uint64) Tree` 构造一个 `Tree`;
   - 然后对该 `Tree` 调用 `Add(intervals)`, `Delete(intervals)`, `Query(interval)` 等；
   - 其中 `Interval` 是一个需要你自行实现的接口，定义坐标、ID、OverlapsAtDimension 等。

> **总结**：这是一个支持多维度区间查询的红黑树实现，会在节点上维护区间 `[LowAtDimension(d), HighAtDimension(d)]` 的最小/最大值，以及颜色属性，用来在插入/删除时保持平衡。在查询时利用 `min/max` 做剪枝，并在必要时检查多维度的 Overlaps。

---

# 二、核心数据结构

## 1. `Interval` 接口

```go
type Interval interface {
    LowAtDimension(uint64) int64
    HighAtDimension(uint64) int64
    OverlapsAtDimension(Interval, uint64) bool
    ID() uint64
}
```

- **LowAtDimension(d)/HighAtDimension(d)**：返回该 `Interval` 在第 `d` 维度的下界/上界（整数）。
- **OverlapsAtDimension(other, d)**：判断在第 `d` 维度上是否与另一个区间 `other` 有重叠。
- **ID**：返回唯一 ID，用于区分相同区间对象的不同实例。

> 例如在 2 维场景下，你的区间对象可以包含 `(x1, x2, y1, y2)`，在 `LowAtDimension(1)` 返回 `x1`, `HighAtDimension(1)` 返回 `x2`, `LowAtDimension(2)` 返回 `y1`, `HighAtDimension(2)` 返回 `y2`。  
> 同时 `OverlapsAtDimension(iv, 1)` 会检查 `[x1, x2]` 是否和 `[iv.x1, iv.x2]` 重叠，以此类推。

## 2. `node` 结构

```go
type node struct {
    interval Interval
    max, min int64
    children [2]*node
    red      bool
    id       uint64
}
```

- `interval`：存储当前节点代表的区间。
- `max, min`：该节点子树（包括自身）中在**第一维**的最大值、最小值，用于快速剪枝。
- `children[0]`/`children[1]`：左右子树指针。
- `red`: 是否是红节点（红黑树属性）；
- `id`: 缓存 `interval.ID()`，以减少方法调用。

## 3. `tree` 结构

```go
type tree struct {
    root                 *node
    maxDimension, number uint64
    dummy                node
}
```

- `root`: 指向根节点；
- `maxDimension`: 表示区间的最大维度（若 = 2，表示 2 维；=1 表示 1 维，等等）；
- `number`: 当前树的节点数 (Len)；
- `dummy`: 一个空节点，用于在插入/删除旋转时的辅助。

> 这些字段在操作（Add/Delete/Query 等）中被反复使用，维护区间树的完整性和统计信息。

---

# 三、重要方法的讲解

## 1. `func (tree *tree) Add(intervals ...Interval)`

- 逐个调用 `tree.add(iv)`。
- `add(iv)`：核心是按照**第一维**的 Low 值来定位插入位置，并进行红黑树旋转、变色等操作；同时更新节点的 `min` / `max` 字段以包含新区间。
- 如果要插入的 ID 已存在，则**不会**重复插入（可以看判断 `node.id == id` 处) 或者覆盖区间。
- 过程中若产生 2-3 级红节点冲突（连续两个红节点），会做`rotate` 或 `doubleRotate` 修正。

## 2. `func (tree *tree) Delete(intervals ...Interval)`

- 逐个调用 `tree.delete(iv)`。
- `delete`：先找到要删的节点，在红黑树中执行“标记 + 替换 + 调整”流程；
- “found != nil” 表示找到匹配节点后，在 RB-tree 中交换节点数据并物理删除；
- 最后将 `tree.root` 设为黑节点以保证根为黑。
- 之后 `adjustRanges()` 用来重新计算 `min` / `max` 等信息。

## 3. `func (tree *tree) Query(interval Interval) Intervals`

- 若 `tree.root == nil`，直接返回空；否则：
- 先取出 interval 的 `[ivLow, ivHigh]` 在**第一维**，再调用 `root.query(ivLow, ivHigh, interval, tree.maxDimension, fn)`.
- `query` 方法会：
  1. 检查左子树是否和 `[low, high]` 有重叠（通过 `overlaps(child.max, high, child.min, low)`），若重叠则递归查询左子树。
  2. 再 `intervalOverlaps(n, low, high, interval, maxDimension)` 判断当前节点是否与查询区间在**所有维度**都重叠，若是就执行回调 `fn(n)`.
  3. 同理对右子树做处理。
- 这样就能得到所有与之重叠的节点的 `interval`，时间复杂度在平均情况下约 \(O(\log n + k)\) （k 为结果大小）。

## 4. `intervalOverlaps(n *node, low, high int64, interval Interval, maxDimension uint64) bool`

- 先用简单 `overlaps(...)` 方法看是否在第1维度重叠：`n.interval.HighAtDimension(1) >= low && n.interval.LowAtDimension(1) <= high`。
- 若 `interval == nil`，直接返回 true；表示只要第一维重叠就行。
- 否则再对 `[2..maxDimension]` 调用 `n.interval.OverlapsAtDimension(interval, i)`，全部通过才认为多维重叠。

---

# 四、如何使用

## 1. 实现你的 `Interval` 类型

你需要提供一个结构体，实现以下方法：

```go
func (my *MyInterval) LowAtDimension(dim uint64) int64 { ... }
func (my *MyInterval) HighAtDimension(dim uint64) int64 { ... }
func (my *MyInterval) OverlapsAtDimension(iv Interval, dim uint64) bool { ... }
func (my *MyInterval) ID() uint64 { ... }
```

假设是 2 维，可能 `MyInterval` 存 `(x1, x2, y1, y2, idVal)`, 其中

- `LowAtDimension(1)` -> `x1`,
- `HighAtDimension(1)` -> `x2`,
- `LowAtDimension(2)` -> `y1`,
- `HighAtDimension(2)` -> `y2`,
- `OverlapsAtDimension(other, 1)` -> `[x1, x2]` overlap `[other.x1, other.x2]`,
- `ID()` 返回一个全局唯一 ID（比如 `uint64` counter）。

## 2. 创建树

```go
tree := New(2) // 表示2维
```

- 传入 `dimensions = 2` 表示有2个维度。

## 3. 向树中添加区间

```go
iv1 := &MyInterval{x1: 0, x2: 10, y1: 5, y2: 15, idVal: 1}
iv2 := &MyInterval{x1: 2, x2: 6,  y1: 7, y2: 12, idVal: 2}
tree.Add(iv1, iv2)
```

## 4. 查询

```go
queryInterval := &MyInterval{x1: 3, x2: 8, y1: 0, y2: 10}
results := tree.Query(queryInterval)
// results: Intervals that overlap with [3..8] x [0..10]
for _, r := range results {
    fmt.Printf("Overlapped with ID=%d\n", r.ID())
}
results.Dispose() // 记得释放 intervalsPool 池
```

> Query 会返回一个 `Intervals` 切片，需要调用 `Dispose()` 来归还到 `sync.Pool`。

## 5. 删除

```go
tree.Delete(iv1) // 删除id=1的interval
```

## 6. 遍历

```go
tree.Traverse(func(iv Interval) {
    fmt.Println("Interval ID:", iv.ID())
})
```

---

# 五、核心要点与注意事项

1. **红黑树 + 维度区间**

   - 从代码可见，很多旋转/变色操作 (rotate/doubleRotate, isRed) 都与红黑树类似，只是这里还维护了区间 `min`/`max`。
   - `min`/`max` 仅对应第 1 维度；若要在更多维度也加快查询，需要类似 `intervalOverlaps(...)` 对其他维度做进一步检查。

2. **多维查询**

   - 在查找是否与目标区间 `interval` 重叠时，先看第 1 维的 overlap 来快速剪枝；
   - 若还需更精准，可以检查各维度的 `OverlapsAtDimension`。
   - 这是典型**多维区间树**的做法，但性能不如“KD-tree”或 R-tree 适合高维；不过在 2~3 维时还算可以。

3. **Pool 回收**

   - `intervalsPool` 是一个 `sync.Pool`，存储 `Intervals` 切片以减少 GC 压力；
   - 每次 Query 都会返回一个新的 `Intervals`，用完后需要 `Dispose()` 来归还到池中。

4. **ID 重复处理**

   - 若相同 ID 的区间再次插入，会**不**重复添加，因为插入过程若发现 `node.id == id` 就会 `break`。
   - 你可以根据业务需要对 `ID()` 进行“时戳+自增”之类策略，或只做对象指针的 hash 也行，但必须保证全局唯一。

5. **自定义 OverlapsAtDimension**

   - 在多维下，需要你自己决定**闭区间**/**开区间**/是否存在特别边界判断；
   - 代码默认**区间是包含端点**(`[low, high]`)重叠。

6. **自定义最大维度**
   - 一般在构造时 `New(2)` 即可用于二维情形；如果只是一维用 `New(1)`，其实就相当于一个带 `min/max` 的区间树 / 红黑树了。

---

# 六、小示例

下面用 1 维场景示例（`dim=1`）展示最小可行用法：

```go
package main

import (
    "fmt"
    // 引入你的interval tree所在包:
    // "github.com/your/repo/intervaltree"
)

// 实现 Interval:
type MyInterval struct {
    start, end int64
    idVal      uint64
}

func (m *MyInterval) LowAtDimension(d uint64) int64 {
    return m.start
}
func (m *MyInterval) HighAtDimension(d uint64) int64 {
    return m.end
}
func (m *MyInterval) OverlapsAtDimension(iv Interval, d uint64) bool {
    // 1维: [start..end] vs [iv.start..iv.end]
    otherLow := iv.LowAtDimension(d)
    otherHigh:= iv.HighAtDimension(d)
    return !(m.end < otherLow || m.start > otherHigh)
}
func (m *MyInterval) ID() uint64 {
    return m.idVal
}

func main() {
    // 1. Create tree with 1 dimension
    t := New(1)

    // 2. Add intervals
    iv1 := &MyInterval{start: 1, end:5, idVal:101}
    iv2 := &MyInterval{start: 3, end:7, idVal:102}
    t.Add(iv1, iv2)

    // 3. Query
    q := &MyInterval{start:4, end:4, idVal:999} // [4..4]
    results := t.Query(q)
    for _, r := range results {
        fmt.Println("Overlapped ID:", r.ID())
    }
    // Overlapped ID: 101
    // Overlapped ID: 102
    results.Dispose()

    // 4. Delete
    t.Delete(iv1)
    fmt.Println("Tree size after delete:", t.Len()) // should be 1

    // 5. Traverse
    t.Traverse(func(iv Interval) {
        fmt.Println("Traverse ID:", iv.ID())
    })
    // Traverse ID: 102
}
```

运行后可见 `[1..5]` 和 `[3..7]` 都与 `[4..4]` 重叠，因此同时被查询出来。再删除 `[1..5]` 后，只剩 `[3..7]`。

---

# 七、总结

- 这是一套多维**区间红黑树**的实现，核心在**节点维护** `(interval, min, max, color, children)`。
- **增/删**：典型红黑树流程 + 更新 `min/max`。
- **查询**：对 `[low, high]` 与节点做**剪枝判断**（first dimension） + **多维 overlap**。
- 提供一个 `Interval` 接口，你必须实现**各维度上下界**、**OverlapsAtDimension**、**唯一 ID**。
- 使用时：
  1. `tree := New(dimensions)`
  2. `tree.Add(myIntervals...)`
  3. `res := tree.Query(someInterval)`
  4. `tree.Delete(...)`
  5. `tree.Traverse(...)` 等
- 适合**在 2~3 维**的区间查询场景下使用，或在 1 维直接当作区间树。较高维度时性能可能衰减，需要更专业的结构（如 R-tree, KD-tree 等）。
- 代码中使用 `sync.Pool` 来缓存 `Intervals` 切片，减少 GC 压力，记得在 Query 后手动 `Dispose()` 释放切片。

这就是该多维 Interval Tree 的实现与用法要点，希望对你有所帮助。
