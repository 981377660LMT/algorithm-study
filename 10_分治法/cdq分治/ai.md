# CDQ分治算法框架（Golang实现）

CDQ分治是一种强大的算法框架，常用于解决多维偏序问题、带修改的查询问题等。
cdq 分治与通常的分治不同。常规的分治法中，对分出的两部分递归求解时，分出的两部分是没有关联的，可以独立求解；而 cdq 分治分出的两部分有关联，不能独立求解。

## 1. 基础框架实现

```go
package cdq

// Item 表示待处理的元素，可根据具体问题定制
type Item struct {
    X, Y, Z int // 三维坐标示例
    Val     int // 元素值
    Index   int // 原始索引
    Result  int // 存储结果
}

// CDQ算法主框架
func SolveCDQ(items []Item) []Item {
    if len(items) == 0 {
        return items
    }

    // 创建一个临时数组用于CDQ分治过程
    temp := make([]Item, len(items))

    // 启动CDQ分治
    cdqDivide(items, temp, 0, len(items)-1)
    return items
}

// cdqDivide 实现CDQ分治的递归过程
func cdqDivide(items, temp []Item, left, right int) {
    if left == right {
        // 处理单个元素的基本情况
        return
    }

    mid := left + (right-left)/2

    // 递归处理左右两部分
    cdqDivide(items, temp, left, mid)
    cdqDivide(items, temp, mid+1, right)

    // 合并阶段：处理跨越左右两部分的贡献
    merge(items, temp, left, mid, right)
}

// merge 处理跨越左右两部分的贡献
func merge(items, temp []Item, left, mid, right int) {
    // 初始化临时数组索引
    i, j, k := left, mid+1, left

    // 这里实现根据具体问题的合并逻辑
    // 通常包括：
    // 1. 按某种顺序归并左右两部分
    // 2. 计算跨越中点的贡献
    // 3. 更新结果

    // 下面是一个示例归并过程(按Y排序并处理贡献)
    bit := NewBIT(100001) // 假设坐标范围是[0, 100000]

    // 按Y排序归并
    for k <= right {
        if j > right || (i <= mid && items[i].Y <= items[j].Y) {
            // 来自左半部分的元素，将其Z加入树状数组
            temp[k] = items[i]
            bit.Update(items[i].Z, 1)
            i++
        } else {
            // 来自右半部分的元素，查询树状数组计算贡献
            temp[k] = items[j]
            // 计算有多少左半部分元素的Z小于当前元素的Z
            items[j].Result += bit.Query(items[j].Z)
            j++
        }
        k++
    }

    // 复位树状数组
    for i := left; i <= mid; i++ {
        bit.Update(items[i].Z, -1)
    }

    // 将临时数组的结果复制回原数组
    for i := left; i <= right; i++ {
        items[i] = temp[i]
    }
}

// BIT 树状数组实现
type BIT struct {
    tree []int
    size int
}

// NewBIT 创建新的树状数组
func NewBIT(size int) *BIT {
    return &BIT{
        tree: make([]int, size+1),
        size: size,
    }
}

// Update 更新树状数组中的值
func (b *BIT) Update(pos, val int) {
    pos++  // 树状数组索引从1开始
    for pos <= b.size {
        b.tree[pos] += val
        pos += pos & -pos // 加上lowbit
    }
}

// Query 查询前缀和
func (b *BIT) Query(pos int) int {
    pos++  // 树状数组索引从1开始
    sum := 0
    for pos > 0 {
        sum += b.tree[pos]
        pos -= pos & -pos // 减去lowbit
    }
    return sum
}
```

## 2. 应用示例：三维偏序问题

```go
package main

import (
    "fmt"
    "sort"
)

// 三维偏序问题：计算对于每个点(x,y,z)，有多少点(a,b,c)满足a<=x, b<=y, c<=z
func solveThreedimOrder(points []Item) []int {
    n := len(points)
    result := make([]int, n)

    // 按X坐标排序
    sort.Slice(points, func(i, j int) bool {
        if points[i].X != points[j].X {
            return points[i].X < points[j].X
        }
        if points[i].Y != points[j].Y {
            return points[i].Y < points[j].Y
        }
        return points[i].Z < points[j].Z
    })

    // 离散化坐标（需要时）
    discretize(points)

    // 应用CDQ分治
    SolveCDQ(points)

    // 整理结果
    for _, p := range points {
        result[p.Index] = p.Result
    }

    return result
}

// 离散化坐标
func discretize(points []Item) {
    // Y坐标离散化
    yValues := make([]int, len(points))
    for i, p := range points {
        yValues[i] = p.Y
    }
    sort.Ints(yValues)
    yMap := make(map[int]int)
    yIdx := 0
    for _, y := range yValues {
        if _, exists := yMap[y]; !exists {
            yMap[y] = yIdx
            yIdx++
        }
    }
    for i := range points {
        points[i].Y = yMap[points[i].Y]
    }

    // Z坐标离散化 (类似Y坐标的处理)
    // ...
}

func main() {
    // 示例用法
    points := []Item{
        {X: 1, Y: 2, Z: 3, Index: 0},
        {X: 2, Y: 1, Z: 3, Index: 1},
        {X: 3, Y: 3, Z: 1, Index: 2},
        // ...更多点
    }

    results := solveThreedimOrder(points)
    fmt.Println("Results:", results)
}
```

## 3. 优化CDQ框架

```go
package cdq

// CDQConfig 配置CDQ算法的参数和行为
type CDQConfig struct {
    // 自定义排序函数
    LessThan func(a, b Item) bool

    // 自定义合并函数
    CustomMerge func(items, temp []Item, left, mid, right int)

    // 自定义基本情况处理
    HandleBaseCase func(item *Item)

    // 其他配置选项
    MaxCoordinate int // 坐标最大值
    UseBIT        bool // 是否使用树状数组
    UseSegTree    bool // 是否使用线段树
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() CDQConfig {
    return CDQConfig{
        LessThan: func(a, b Item) bool {
            return a.X < b.X
        },
        CustomMerge: nil,
        HandleBaseCase: func(item *Item) {},
        MaxCoordinate: 100000,
        UseBIT: true,
        UseSegTree: false,
    }
}

// SolveCDQWithConfig 使用自定义配置的CDQ分治
func SolveCDQWithConfig(items []Item, config CDQConfig) []Item {
    if len(items) == 0 {
        return items
    }

    temp := make([]Item, len(items))

    // 如果提供了自定义排序，先排序
    if config.LessThan != nil {
        sort.Slice(items, func(i, j int) bool {
            return config.LessThan(items[i], items[j])
        })
    }

    // 调用CDQ分治
    cdqDivideWithConfig(items, temp, 0, len(items)-1, config)
    return items
}

func cdqDivideWithConfig(items, temp []Item, left, right int, config CDQConfig) {
    if left == right {
        // 处理基本情况
        if config.HandleBaseCase != nil {
            config.HandleBaseCase(&items[left])
        }
        return
    }

    mid := left + (right-left)/2

    // 递归处理左右两部分
    cdqDivideWithConfig(items, temp, left, mid, config)
    cdqDivideWithConfig(items, temp, mid+1, right, config)

    // 合并阶段
    if config.CustomMerge != nil {
        config.CustomMerge(items, temp, left, mid, right)
    } else {
        defaultMerge(items, temp, left, mid, right, config)
    }
}

func defaultMerge(items, temp []Item, left, mid, right int, config CDQConfig) {
    // 默认合并实现，可根据配置使用不同的数据结构
    // ...
}
```

这个CDQ分治框架实现了基本的分治逻辑和辅助数据结构，并提供了多种定制选项。可以根据具体问题调整配置，包括自定义排序、合并函数等，以实现高效的解决方案。
