以下是支持区间修改（带懒标记）的RadixTree修改代码：

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"time"
)

type RadixTree[E any] struct {
	e         func() E
	op        func(a, b E) E      // 聚合操作
	addOp     func(a E, b E) E    // 标记叠加操作
	mergeLazy func(a E, b E) E    // 标记合并操作
	log       int
	blockSize int

	n           int
	data        []E
	levels      [][]E
	levelShifts []int
	lazy        [][]E // 每层的懒标记
}

func NewRadixTree[E any](e func() E, op, addOp, mergeLazy func(a, b E) E, log int) *RadixTree[E] {
	if log < 1 {
		log = 6
	}
	return &RadixTree[E]{
		e:         e,
		op:        op,
		addOp:     addOp,
		mergeLazy: mergeLazy,
		log:       log,
		blockSize: 1 << log,
	}
}

func (m *RadixTree[E]) Build(n int, f func(i int) E) {
	m.n = n
	m.data = make([]E, n)
	for i := 0; i < n; i++ {
		m.data[i] = f(i)
	}
	m.levels = [][]E{}
	m.levelShifts = []int{}
	m.lazy = [][]E{}

	build := func(pre []E) []E {
		cur := make([]E, (len(pre)+m.blockSize-1)>>m.log)
		for i := range cur {
			start := i << m.log
			end := min(start+m.blockSize, len(pre))
			v := m.e()
			for j := start; j < end; j++ {
				v = m.op(v, pre[j])
			}
			cur[i] = v
		}
		return cur
	}

	preLevel := m.data
	preShift := 1
	for len(preLevel) > 1 {
		curLevel := build(preLevel)
		m.levels = append(m.levels, curLevel)
		m.levelShifts = append(m.levelShifts, m.log*preShift)
		m.lazy = append(m.lazy, make([]E, len(curLevel)))
		for i := range m.lazy[len(m.lazy)-1] {
			m.lazy[len(m.lazy)-1][i] = m.e()
		}
		preLevel = curLevel
		preShift++
	}
}

// 核心：下传懒标记
func (m *RadixTree[E]) pushDown(k, bid int) {
	if k == 0 {
		return // 最后一层不需要下传
	}

	if m.lazy[k][bid] != m.e() {
		// 计算子块范围
		shift := m.levelShifts[k-1]
		start := bid << (m.levelShifts[k] - m.levelShifts[k-1])
		end := min(start+(1<<(m.levelShifts[k]-m.levelShifts[k-1])), len(m.levels[k-1]))

		// 更新子层级的值和懒标记
		for i := start; i < end; i++ {
			m.levels[k-1][i] = m.addOp(m.levels[k-1][i], m.lazy[k][bid])
			if k-1 > 0 {
				m.lazy[k-1][i] = m.mergeLazy(m.lazy[k-1][i], m.lazy[k][bid])
			}
		}
		m.lazy[k][bid] = m.e()
	}
}

func (m *RadixTree[E]) UpdateRange(l, r int, value E) {
	if l < 0 {
		l = 0
	}
	if r > m.n {
		r = m.n
	}
	if l >= r {
		return
	}
	m.updateRangeRecursive(l, r, value, len(m.levels)-1)
}

func (m *RadixTree[E]) updateRangeRecursive(l, r int, value E, k int) {
	if k < 0 {
		for i := l; i < r; i++ {
			m.data[i] = m.addOp(m.data[i], value)
		}
		return
	}

	shift := m.levelShifts[k]
	startBlock := l >> shift
	endBlock := (r - 1) >> shift

	// 处理完全覆盖的块
	if startBlock == endBlock {
		m.pushDown(k+1, startBlock>>(m.levelShifts[k+1]-shift))
		m.updateRangeRecursive(l, r, value, k-1)
		m.refresh(k, startBlock)
		return
	}

	// 处理中间完整块
	for i := startBlock + 1; i < endBlock; i++ {
		m.levels[k][i] = m.addOp(m.levels[k][i], value)
		m.lazy[k][i] = m.mergeLazy(m.lazy[k][i], value)
	}

	// 处理左半部分
	leftEnd := (startBlock + 1) << shift
	if leftEnd > l {
		m.pushDown(k+1, startBlock>>(m.levelShifts[k+1]-shift))
		m.updateRangeRecursive(l, leftEnd, value, k-1)
		m.refresh(k, startBlock)
	}

	// 处理右半部分
	rightStart := endBlock << shift
	if rightStart < r {
		m.pushDown(k+1, endBlock>>(m.levelShifts[k+1]-shift))
		m.updateRangeRecursive(rightStart, r, value, k-1)
		m.refresh(k, endBlock)
	}
}

func (m *RadixTree[E]) refresh(k, bid int) {
	if k == 0 {
		return
	}
	shift := m.levelShifts[k] - m.levelShifts[k-1]
	start := bid << shift
	end := min(start+(1<<shift), len(m.levels[k-1]))

	v := m.e()
	for i := start; i < end; i++ {
		v = m.op(v, m.levels[k-1][i])
	}
	m.levels[k][bid] = v
}

// 其他原有方法需要增加懒标记处理...

// 示例修改后的QueryRange方法
func (m *RadixTree[E]) QueryRange(l, r int) E {
	if l < 0 {
		l = 0
	}
	if r > m.n {
		r = m.n
	}
	if l >= r {
		return m.e()
	}
	return m.queryRangeRecursiveWithLazy(l, r, len(m.levels)-1)
}

func (m *RadixTree[E]) queryRangeRecursiveWithLazy(l, r, k int) E {
	if k >= 0 {
		m.pushDown(k+1, l>>m.levelShifts[k+1])
	}

	// 原有查询逻辑...
}

// 其他需要修改的方法（Update/Set等）需要处理懒标记...
```

主要修改点说明：

1. **新增字段**：

   - `addOp`: 定义如何将标记应用到元素
   - `mergeLazy`: 定义如何合并标记
   - `lazy`: 多层懒标记存储

2. **Build初始化**：

   - 为每层初始化懒标记数组

3. **新增核心方法**：

   - `pushDown`: 将当前层的标记下传到下层
   - `refresh`: 更新父层聚合值
   - `UpdateRange`: 区间更新入口

4. **修改查询方法**：
   - 在查询前需要执行pushDown保证数据最新

使用时需要提供三个操作函数：

- `op`: 聚合操作（原功能）
- `addOp`: 元素与标记的叠加操作（如加法：func(a,b int)int {return a+b}）
- `mergeLazy`: 标记的合并操作（如标记叠加：func(a,b int)int {return a+b}）

示例初始化：

```go
// 加法示例
rt := NewRadixTree(
    func() int { return 0 },
    func(a, b int) int { return a + b },
    func(a, b int) int { return a + b },
    func(a, b int) int { return a + b },
    4
)
```

注意事项：

1. 需要保证操作满足结合律
2. 幺元需要与操作匹配
3. 标记合并操作需要满足幂等性
4. 更新时需要先pushDown再操作
5. 查询时需要保证路径上的标记都已下传

完整实现需要考虑所有涉及聚合操作的方法都需要处理懒标记，包括现有的Update和Set方法。具体实现需要根据实际需求调整标记处理逻辑。
