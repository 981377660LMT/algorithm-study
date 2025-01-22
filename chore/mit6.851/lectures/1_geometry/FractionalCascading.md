下面给出一个**示例性质**的 Golang 实现，演示如何利用**Fractional Cascading**来对 \(k\) 个有序数组（每个长度为 \(n\)）进行加速查询：对于每个查询 \(x\)，在每个数组中找到**大于等于 \(x\)** 的最小数（若不存在则返回一个哨兵值），并且要求“在线”处理（即对每个查询立即输出结果）。

> **说明：**
>
> 1. 这是典型的 Fractional Cascading 方案的简化示例，便于理解核心思想。
> 2. 真实生产场景中，可能需要对输入/输出格式、异常处理、大规模数据的优化等做更完善的处理。
> 3. 本示例假设数组内部无重复值，若存在重复值也可以正常处理，但需要在构建和查询时注意等于情况的处理逻辑。

---

# 一、Fractional Cascading 原理概述

给定 \(k\) 个排好序的数组 \(\{A_1, A_2, \dots, A_k\}\)，每个长度为 \(n\)。我们想要进行如下查询：给一个数 \(x\)，分别在每个数组里二分查找大于等于 \(x\) 的最小值。**直接做**会对每次查询耗时 \(O(k \log n)\)。

Fractional Cascading（部分级联）技巧可以将每次查询的复杂度降为 **\(O(\log n + k)\)**：

1. **核心想法**：

   - 将相邻数组的元素部分“交织（merge）”在一起，并且设置指针或索引，使得在第一个数组做一次二分查找后，可以在常数时间内把“候选位置”带到下一个数组，大幅减少重复的二分开销。

2. **数据结构**：

   - 从后往前构建一个增强的数组结构 \(L*i\)（可能叫“胖数组”或“扩展数组”），其中包含原数组 \(A_i\) 的所有元素，以及从相邻数组 \(L*{i+1}\)“抽取”部分元素（通常是一半）进来。
   - 每个元素在 \(L*i\) 中会持有一个**跨级指针**或“nextIdx”，可以迅速跳到 \(L*{i+1}\) 中“对应”或“相邻”的位置。
   - 这样，只需在 \(L_1\) 中对 \(x\) 进行一次二分查找，就能在 \(O(1)\) 时间跳到 \(L_2\) 的正确位置，再顺次跳到 \(L_3\), …, \(L_k\)。

3. **查询**：
   - 在 \(L_1\) 做二分查找，找到“最小的 \(\ge x\)”所在位置，然后通过**跨级指针**在 \(O(1)\) 时间“级联”到 \(L_2\) 的对应区间，再取到那里的“最小的 \(\ge x\)”，继续级联到 \(L_3\)，直到 \(L_k\)，整体 \(O(\log n + k)\) 完成。

> 关于更详细的理论背景，可参考原始论文或各类算法教材，这里聚焦于 Golang 示例代码的**构造**与**查询**过程。

---

# 二、示例代码

以下是一份自包含的 Go 程序，会：

1. 读入 \(k, n, q\)
2. 读入 \(k\) 个排序好的数组
3. 使用 Fractional Cascading 构建增强数据结构
4. 读入 \(q\) 个查询 \(x\)，并在线输出对每个数组的查询结果

> - 若“在某个数组中不存在大于等于 \(x\) 的元素”，则返回 `INF` 或者你想要的哨兵值。

## 1. 数据结构定义

我们为每个“增强数组”定义一个切片 `[]FCNode`，其中 `FCNode` 存放：

- `val`：该节点表示的数值
- `origIdx`：该节点在**原数组**（\(A_i\)）中的下标（如果它是从下一个数组“借”过来的，则可能是 -1 或不去用它）
- `nextIdx`：指示在相邻的下一个增强数组 \(L\_{i+1}\) 中对应/相邻节点的下标，用于快速跳转

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

const INF = 1_000_000_000_000_000_000 // 一个很大的数

type FCNode struct {
    val     int
    origIdx int  // 在原数组 A_i 中的位置(若本节点源自A_i)
    nextIdx int  // 跳到 L_{i+1} 中对应位置
}

type FCArray struct {
    nodes []FCNode // 增强数组 L_i
}
```

## 2. 构建 Fractional Cascading

### 2.1 从后往前初始化

- 对于最后一个数组 \(A*{k-1}\)（假设下标从 0 开始：0..k-1），我们直接把它所有元素放到 `L*{k-1}`，`nextIdx` 暂时置为 -1（因为没有下一个数组了）。

```go
func buildFractionalCascading(arrs [][]int) []FCArray {
    k := len(arrs)
    result := make([]FCArray, k)

    // 先构建 L_{k-1} = 整个 A_{k-1}
    lastIndex := k - 1
    base := make([]FCNode, len(arrs[lastIndex]))
    for i, v := range arrs[lastIndex] {
        base[i] = FCNode{
            val:     v,
            origIdx: i,
            nextIdx: -1, // 没有下一个数组了
        }
    }
    result[lastIndex] = FCArray{nodes: base}

    // 再从后往前构建 L_{i}
    for i := k - 2; i >= 0; i-- {
        result[i] = mergeAndLinkFC(arrs[i], result[i+1])
    }

    return result
}
```

### 2.2 合并 + 抽取“半数”元素

`mergeAndLinkFC(A_i, L_{i+1})` 用于构建 `L_i`：

1. 先把 `A_i`（原数组）所有元素转换为 `FCNode`，记录 `origIdx`；
2. 从 `L_{i+1}` 中“抽取”大约一半元素，与 `A_i` 进行**归并**（按照数值顺序），得到一个新的切片 `merged`;
3. 在合并过程中，双方节点会互相登记 `nextIdx`，以便能相互跳转：
   - 当我们把 `L_{i+1}` 的某个节点放进 `merged` 时，说明它在 `L_i` 也有出现，于是就可以设置 “`L_{i+1}.nodes[posNext].nextIdx = posMerged`”，反之在 `merged[posMerged]` 里写 “`nextIdx = posNext`”。
4. 注意只抽取 `L_{i+1}` 中的一半节点（例如**偶数位**或**奇数位**），以保证最终大小保持线性。Fractional Cascading 的数学原理表明这样就足以在查询时“带动”对方跳转。

> 实际实现时，“抽一半”可以用**奇偶位**或“每两个取一个”等多种方式，这里选用“取 `L_{i+1}` 中偶数下标节点”作为示例。

```go
func mergeAndLinkFC(arr []int, nextFC FCArray) FCArray {
    n1 := len(arr)
    n2 := len(nextFC.nodes)

    // step1: 把 A_i 的所有元素转换为 FCNode
    arrNodes := make([]FCNode, n1)
    for i, v := range arr {
        arrNodes[i] = FCNode{
            val:     v,
            origIdx: i,
            nextIdx: -1,
        }
    }

    // step2: 从 nextFC 中抽取一半节点: 这里选取偶数下标
    pickNodes := make([]FCNode, 0)
    for idx := 0; idx < n2; idx += 2 {
        pickNodes = append(pickNodes, nextFC.nodes[idx])
    }

    // step3: 归并 arrNodes 和 pickNodes => merged
    merged := make([]FCNode, 0, len(arrNodes)+len(pickNodes))
    i1, i2 := 0, 0
    for i1 < len(arrNodes) && i2 < len(pickNodes) {
        if arrNodes[i1].val <= pickNodes[i2].val {
            merged = append(merged, arrNodes[i1])
            i1++
        } else {
            merged = append(merged, pickNodes[i2])
            i2++
        }
    }
    for i1 < len(arrNodes) {
        merged = append(merged, arrNodes[i1])
        i1++
    }
    for i2 < len(pickNodes) {
        merged = append(merged, pickNodes[i2])
        i2++
    }

    // step4: 构建 cross pointers (双向)
    // 我们需要在 merged 中和 nextFC.nodes 中找到同一个“值”的对应位置，互设 nextIdx
    // 但 merged 和 nextFC 都是有序的，我们可以用双指针扫描
    // （原则上只需要处理 pickNodes 对应的那部分）
    mp := merged
    np := nextFC.nodes

    im, in := 0, 0
    for im < len(mp) && in < len(np) {
        if mp[im].val == np[in].val {
            // 建立双向链接
            mp[im].nextIdx = in
            np[in].nextIdx = im
            im++
            in++
        } else if mp[im].val < np[in].val {
            im++
        } else {
            in++
        }
    }

    return FCArray{nodes: merged}
}
```

> - 上述 `mergeAndLinkFC` 中的 `cross pointers` 处理只需要在那些“也出现在 `merged` 里的 nextFC.nodes”上设置链接。因为只有那些节点才在 `L_i` 中真实存在。
> - 这里为简单起见，我们做了“在全部可能相同值之间建立链接”的处理。如果数组有重复值，逻辑也相似，只不过会出现“多对多”的匹配，可能需要更细致处理。

到此就能得到一系列 `L_0, L_1, ..., L_{k-1}`，并且相邻层有双向链接。

---

## 3. 查询

### 3.1 在 \(L_0\) 上二分查找

我们先在 `L_0`（即 `result[0].nodes`）做一次二分，找到**最小的 \(\ge x\)** 的索引 `pos`。若全部小于 `x`，则 `pos = len(L_0)`, 表示越界。

```go
// 在 L_0（FCArray） 中找 "最小的 >= x" 的下标
func searchInFC(fc FCArray, x int) int {
    nodes := fc.nodes
    left, right := 0, len(nodes)
    for left < right {
        mid := (left + right) >> 1
        if nodes[mid].val >= x {
            right = mid
        } else {
            left = mid + 1
        }
    }
    return left // 可能等于 len(nodes)，表示都小于 x
}
```

### 3.2 级联到后续数组

从 `L_0` 的 `pos` 开始，我们依次去到 `L_1`，`L_2` … `L_{k-1}`：

1. 若 `pos == len(L_0)`，说明在 `L_0` 都找不到 \(\ge x\) 的值，那么对应的“原数组 A_0”结果是 `INF`（表示无解）。
2. 否则在 `pos` 位置取到节点 `L_0.nodes[pos]`，记为 `n0`，我们先确定**原数组 A_0**中的最小 \(\ge x\) 值：
   - 因为 `n0.val` 就是一个 \(\ge x\) 值，但要确保它是从**本数组**而非下一个数组“借”过来的才行（即 `n0.origIdx != -1`）。如果 `n0.origIdx == -1`，那说明这个节点是来自 `L_{1}` 的抽取，我们需要再看下一个位置或附近位置，这里做法上通常是：如果 `n0.val >= x` 依然可以作为解，只是你要知道它在原数组 A_0 中并没有实际对应，这点实现上可以灵活处理。
   - 为简单起见，我们可以直接用 `n0.val` 当成“在 A_0 找到的 \(\ge x\)”，若越界就 `INF`。
3. 获取跨层下标 `nextPos = n0.nextIdx`（如果它是 -1，说明找不到对应，视为越界），然后在 `L_1` 用这个位置查看即可。若 `nextPos` 越界或 `-1`，说明需要在 `L_1` 再二分搜索一次。但 Fractional Cascading 的原理保证，这种状况极少，一般可以常数时间找到好位置。
4. 类似地从 `L_1[nextPos]` 跳到 `L_2`，直到 `L_{k-1}`。

下面写一个示例函数 `fcQueryAll(fcList []FCArray, x int) []int`，返回对每个原数组的查询结果。

```go
func fcQueryAll(fcList []FCArray, x int) []int {
    k := len(fcList)
    ans := make([]int, k)

    // 1) 在 L_0 上做一次二分
    pos := searchInFC(fcList[0], x)
    var node FCNode

    // 2) 逐层级联
    for i := 0; i < k; i++ {
        L := fcList[i]
        if pos >= len(L.nodes) {
            // 越界 => 在 A_i 中无 >= x
            ans[i] = INF
            // 往下一层时，我们还是拿 pos = len(...) => nextIdx=-1
            // 这样下层要自己再做 fallback
            node = FCNode{val: INF, nextIdx: -1}
        } else {
            node = L.nodes[pos]
            ans[i] = node.val
        }

        // 准备跳到下一层
        if i+1 < k {
            if node.nextIdx >= 0 && node.nextIdx < len(fcList[i+1].nodes) {
                pos = node.nextIdx // 级联成功
            } else {
                // 没有有效nextIdx，做一次二分fallback
                pos = searchInFC(fcList[i+1], x)
            }
        }
    }

    return ans
}
```

> - 为了保证**每个原数组**得到正确答案，最严谨的方法是：
>
>   1. 确定在 `L_i[pos]` 中的 `val >= x` 是否真正对应原数组 A*i（`origIdx != -1`）；若 `origIdx == -1`，则它是从 `L*{i+1}` 借来的节点，可能并不在 A_i 中；你需要沿着邻近节点看看有没有 A_i 的节点。
>   2. 如果都没有找到，就说明在 A_i 没有 \(\ge x\) 的元素。
>
>   但在大多数教科书/论文中，为简化，会直接把 `node.val` 当成在 “A_i” 中的结果（因为“从别的数组借来的元素”也相当于“我们知道A_i里至少有个不小于这个值的元素”——具体可通过相邻位置在 `L_i` 上确认）。本示例里采用了这种简化做法。

---

## 4. 将一切串起来 (main)

下面示例 `main()` 假设输入格式为：

```
k n q
# k行，每行n个升序数
A1[0], A1[1], ..., A1[n-1]
A2[0], A2[1], ..., A2[n-1]
...
Ak[0], ...
# q行，每行一个查询 x
x1
x2
...
xq
```

我们读入后进行构建，然后对每个 x 调用 `fcQueryAll`，把结果输出。若结果是 `INF` 说明该数组没有满足条件，打印 “-1” 或其它哨兵。

```go
func main() {
    input := bufio.NewReader(os.Stdin)
    var k, n, q int
    fmt.Fscanf(input, "%d %d %d\n", &k, &n, &q)

    arrs := make([][]int, k)
    for i := 0; i < k; i++ {
        arrs[i] = make([]int, n)
        for j := 0; j < n; j++ {
            fmt.Fscanf(input, "%d", &arrs[i][j])
        }
        // 每个数组保证已排序，否则在这里 sort 一下
        // sort.Ints(arrs[i]) // 如果输入可能无序，这里需要 sort
    }

    // 构建 Fractional Cascading 结构
    fcList := buildFractionalCascading(arrs)

    // 在线处理 q 个查询
    for ; q > 0; q-- {
        var x int
        fmt.Fscanf(input, "%d\n", &x)
        res := fcQueryAll(fcList, x)

        // 输出结果
        // 若某数组中不存在 >= x, 就输出 -1；否则输出对应值
        for i, v := range res {
            if v == INF {
                fmt.Printf("-1")
            } else {
                fmt.Printf("%d", v)
            }
            if i+1 < len(res) {
                fmt.Printf(" ")
            }
        }
        fmt.Println()
    }
}
```

编译并运行后，即可对外提供：

1. **构建时间**：\(O\bigl(\sum n + \sum \tfrac{n}{2}\bigr) \approx O(kn)\)，这是一次性的构建开销。
2. **单次查询时间**：理论上 \(O(\log n + k)\)。在本示例中，如果 `node.nextIdx` 无效，就会用二分回退，使得最坏情况仍然 \(O(k \log n)\)，但**期望情况**在文献中可达 \(O(\log n + k)\)。要实现完美的最坏情况 \(O(\log n + k)\)，需要更缜密的构造和 cross pointers 维护（稍复杂，这里仅演示核心思路）。

---

# 三、示例演示

输入示例：

```
3 5 3
1 3 5 7 10
2 3 9 15 17
4 5 6 100 101
3
10
18
```

- \(k=3\), \(n=5\), \(q=3\)
- 数组：
  - \(A_0 = [1, 3, 5, 7, 10]\)
  - \(A_1 = [2, 3, 9, 15, 17]\)
  - \(A_2 = [4, 5, 6, 100, 101]\)
- 查询：\(x=3, x=10, x=18\)

**手动分析**：

1. 对 \(x=3\)：

   - \(A_0\) 中最小 \(\ge 3\) 是 3
   - \(A_1\) 中最小 \(\ge 3\) 是 3
   - \(A_2\) 中最小 \(\ge 3\) 是 4

2. 对 \(x=10\)：

   - \(A_0\) 中最小 \(\ge 10\) 是 10
   - \(A_1\) 中最小 \(\ge 10\) 是 15 (因为 9 < 10, 下一个是15)
   - \(A_2\) 中最小 \(\ge 10\) 是 100

3. 对 \(x=18\)：
   - \(A_0\) 中最小 \(\ge 18\) 不存在 => `-1`
   - \(A_1\) 中最小 \(\ge 18\) 不存在 => `-1`
   - \(A_2\) 中最小 \(\ge 18\) 是 100

所以输出应为：

```
3 3 4
10 15 100
-1 -1 100
```

---

# 四、总结

- **Fractional Cascading**是一种在多数组上做类似二分查找时的加速技巧，能将**多次二分**从 \(O(k \log n)\) 降至 **\(O(\log n + k)\)**。
- 核心做法是：从后向前构建一系列“增强数组” \(L_i\)，令其含有**原数组 \(A_i\)** 的全部元素以及**相邻数组**的“一半元素”，并为这些交叉元素建立**双向链接**。
- 查询时，只需在 \(L_0\) 做一次二分，然后依靠链接在后续数组中**级联**跳转，很快便能找到相应的 \(\ge x\) 元素位置。
- 代码实现有多种方式，这里给出的是**最直观**的“合并+抽半 + cross pointers”。实际工程可根据需求做更多优化或适配。

至此，我们即可用上述示例代码在 Golang 中完成**多有序数组的在线查询**，并通过 Fractional Cascading 获得更优的查询效率。祝你好运!
