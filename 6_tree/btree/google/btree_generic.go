// 版权所有（C）2014-2022 Google Inc.
//
// 依据 Apache License 2.0（以下简称 “许可证”）授权；
// 使用本文件需遵守许可证条款。除非适用法律要求或书面同意，
// 您不可使用本文件。可在以下地址查看许可证副本：
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// 除非适用法律要求或已书面同意，按照许可证发布的软件
// “按原样”提供，不附带任何明示或暗示的担保或条件。
// 有关许可证下的权限与限制的详细信息，请参阅上述许可证。

//go:build go1.18
// +build go1.18

// 在 Go 1.18 及以上版本中，会创建一个名为 BTreeG 的泛型版本，
// 而 BTree 则是该泛型针对 Item 接口的特定实例化版本，同时保留向后兼容的 API。
// 在 go1.18 之前，不支持泛型，BTree 只能基于 Item 接口来实现。

// 包 btree 实现了可在内存中使用的可调阶 B-Tree。
//
// btree 用于在内存中维护一个有序数据结构，并不适用于持久化存储方案。
//
// 与等价的红黑树或其他二叉树相比，它的结构更“扁平”，
// 在某些情况下可以带来更好的内存利用率和/或性能。
// 可参见以下链接中对此的部分讨论：
//
//	http://google-opensource.blogspot.com/2013/01/c-containers-that-save-memory-and-time.html
//
// 需要注意的是，本项目与上述链接中的 C++ B-Tree 实现并无任何关联。
//
// 在本实现中，每个节点包含一个 items 切片，以及一个（可能为 nil）
// 的 children 切片。对于基本数值类型或原始结构体，
// 与等价的 C++ 模板代码（在节点内部使用数组存储）相比，
// 可能会出现一些效率差异：
//   - 由于使用 interface 存储时的额外开销（每个值本身、再加上两字长的 interface 指针，
//     分别指向值及其类型），导致内存使用量更高。
//   - 接口指针可能指向内存中的任意位置，因而值很可能并不是连续存放，
//     这会导致更多的缓存未命中。
//
// 不过，当处理字符串或其他堆分配的结构时，这些问题通常并不显著，
// 因为等价的 C++ 结构中同样要存储指针，并在堆上分配它们的值。
//
// 本实现旨在可无缝替代 gollrb.LLRB（http://github.com/petar/gollrb），
// 这是目前 Go 生态中最常用的有序树实现之一，且功能十分优秀。
// 因此本实现的函数接口在可能的情况下与 llrb.LLRB 完全对应。
// 不同的是，gollrb 支持在树中存储多个相同值，本实现目前并不支持这一点。
//
// 该包有两个实现版本：
// 带 “G” 后缀的版本使用泛型，可处理任意类型，并要求调用方提供一个 “less” 函数来定义排序方式；
// 不带 “G” 后缀的版本则专门针对 “Item” 接口，使用其中的 “Less” 方法进行排序。
//
//
//
// 该版本针对 Go 1.18+ 增加了泛型支持，内部使用了“副本写时复制”（copy-on-write）的技巧，
// 以便在 Clone 之后，共享原有节点并在修改时才进行真正的复制。
//
// Api:

// - **数据结构**：
//   1. `BTreeG[T]`：泛型 B-Tree 的结构，包含树的阶、root、元素总数，以及 copy-on-write 的上下文。
//   2. `node[T]`：具体的 B-Tree 节点，内部有 `items` (关键字) 与 `children` (子节点)。
//   3. `copyOnWriteContext[T]`：管理节点的“所有权”，实现克隆后节点共享，写操作时复制。

// - **核心操作**：
//   1. **插入**：`ReplaceOrInsert`，若节点满则拆分，递归向下直至叶子；支持替换已存在的等价元素。
//   2. **删除**：`Delete` / `DeleteMin` / `DeleteMax` 等，通过借用/合并子节点来保持 B-Tree 性质。
//   3. **查找**：`Get`/`Has`，节点中二分查找后递归下去，时间复杂度约 O(log n)。
//   4. **遍历**：`Ascend`/`Descend`，支持区间查询、从小到大或从大到小依次处理元素。
//   5. **克隆**：`Clone()`，使用写时复制手段，能在读多写少的场景中高效共享节点。
//   6. **清空**：`Clear()` 快速释放树内节点，可选地把节点放回 `freelist` 以重用。

// - **应用场景**：
//   - 纯内存的有序数据存储，对插入、删除、范围查询需求较多。
//   - 通过“扁平”结构（节点内存放多个 key）减小树的高度，与红黑树或 AVL 树相比，实际性能可能更好，尤其在读多写少时。
//   - 不适用于持久化到磁盘的场景（没有专门的页级结构或 I/O 优化），而是更适合在内存中维护大量有序数据，比如缓存、搜索数据结构等。

package main

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
)

func main() {
	flag.Int("n", 0, "number of lines to read")
	tree := NewOrderedG[int](32)

	tree.ReplaceOrInsert(1)

	fmt.Println(tree.Get(1))
}

// Item 表示树中的一个元素。
type Item interface {
	// Less 用于判断当前元素是否小于给定的元素。
	//
	// 该函数需提供一个严格弱序关系。
	// 如果 !a.Less(b) && !b.Less(a)，则视为 a == b（即树中仅会保留 a 或 b 其中之一）。
	Less(than Item) bool
}

const (
	DefaultFreeListSize = 32
)

// FreeListG：节点重用机制
// FreeListG 表示一个存放 btree 节点的空闲列表。默认情况下，每个
// BTree 都有自己的 FreeList，但多个 BTree 可以共享同一个 FreeList，
// 尤其是在它们通过 Clone 创建时。
// 如果两个 BTree 使用同一个 freelist，则可以安全地在并发写场景下使用。
type FreeListG[T any] struct {
	freelist []*node[T] // 用于缓存已被释放的节点，以免频繁地向 Go runtime 申请/GC
}

// NewFreeListG 创建一个新的空闲列表。
// size 是该列表所能容纳的最大节点数量。
func NewFreeListG[T any](size int) *FreeListG[T] {
	return &FreeListG[T]{freelist: make([]*node[T], 0, size)}
}

func (f *FreeListG[T]) newNode() (n *node[T]) {
	index := len(f.freelist) - 1
	if index < 0 {
		return new(node[T])
	}
	n = f.freelist[index]
	f.freelist[index] = nil
	f.freelist = f.freelist[:index]
	return
}

func (f *FreeListG[T]) freeNode(n *node[T]) (out bool) {
	if len(f.freelist) < cap(f.freelist) {
		f.freelist = append(f.freelist, n)
		out = true
	}
	return
}

// ItemIteratorG 用于在 {A/De}scend* 中按照顺序遍历树的某些部分。
// 当此函数返回 false 时，遍历会停止，对应的 Ascend* 函数会立即返回。
type ItemIteratorG[T any] func(item T) bool

// Ordered 表示可以直接使用 '<' 运算符比较的类型集合。
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

// Less[T] 返回一个对支持 '<' 运算符的类型所使用的默认比较函数。
func Less[T Ordered]() LessFunc[T] {
	return func(a, b T) bool { return a < b }
}

// NewOrderedG 创建一棵针对有序类型 T 的 B-Tree。
func NewOrderedG[T Ordered](degree int) *BTreeG[T] {
	return NewG[T](degree, Less[T]())
}

// NewG 根据给定的阶（degree）和比较函数（less）创建一棵泛型 B-Tree。
// 比如 NewG(2) 会创建一个 2-3-4 树（每个节点包含 1-3 个元素以及 2-4 个孩子）。
func NewG[T any](degree int, less LessFunc[T]) *BTreeG[T] {
	return NewWithFreeListG(degree, less, NewFreeListG[T](DefaultFreeListSize))
}

// NewWithFreeListG 创建一棵使用指定空闲列表的 B-Tree。
func NewWithFreeListG[T any](degree int, less LessFunc[T], f *FreeListG[T]) *BTreeG[T] {
	if degree <= 1 {
		panic("bad degree")
	}
	return &BTreeG[T]{
		degree: degree,
		cow:    &copyOnWriteContext[T]{freelist: f, less: less},
	}
}

// items 是节点中存放元素的切片。
type items[T any] []T

// insertAt 在给定索引处插入一个值，后续所有值向后移动。
func (s *items[T]) insertAt(index int, item T) {
	var zero T
	*s = append(*s, zero)
	if index < len(*s) {
		copy((*s)[index+1:], (*s)[index:])
	}
	(*s)[index] = item
}

// removeAt 移除给定索引处的值，后续所有值向前移动。
func (s *items[T]) removeAt(index int) T {
	item := (*s)[index]
	copy((*s)[index:], (*s)[index+1:])
	var zero T
	(*s)[len(*s)-1] = zero
	*s = (*s)[:len(*s)-1]
	return item
}

// pop 移除并返回切片的最后一个元素。
func (s *items[T]) pop() (out T) {
	index := len(*s) - 1
	out = (*s)[index]
	var zero T
	(*s)[index] = zero
	*s = (*s)[:index]
	return
}

// truncate 将切片截断到 index，使其只包含前 index 个元素。index 必须 <= 当前长度。
func (s *items[T]) truncate(index int) {
	var toClear items[T]
	*s, toClear = (*s)[:index], (*s)[index:]
	var zero T
	for i := 0; i < len(toClear); i++ {
		toClear[i] = zero
	}
}

// find 返回给定元素应插入到此切片的索引位置。如果在该位置已经有相等元素，则返回 found = true。
func (s items[T]) find(item T, less func(T, T) bool) (index int, found bool) {
	i := sort.Search(len(s), func(i int) bool {
		return less(item, s[i])
	})
	if i > 0 && !less(s[i-1], item) {
		return i - 1, true
	}
	return i, false
}

// node 表示 B-Tree 中的一个节点。维护 B-Tree 插入/删除时的局部操作。
// 它必须保持以下不变式：
//   - 若 len(children) == 0，则 len(items) 可任意
//   - 否则必须保证 len(children) == len(items) + 1
type node[T any] struct {
	items    items[T] // 存放节点的 key（有序）
	children items[*node[T]]
	cow      *copyOnWriteContext[T] // copy-on-write（写时复制）的上下文指针，说明该节点属于哪个“上下文”
}

// !“写操作”前的可写校验，可理解为cloneNode()
// mutableFor 检查当前节点是否与给定的 COW 上下文匹配，若不匹配则复制。
func (n *node[T]) mutableFor(cow *copyOnWriteContext[T]) *node[T] {
	if n.cow == cow {
		// 如果节点“归属”就是当前树的上下文，就可以直接改
		return n
	}

	// 否则（它属于另一棵树或另一 COW 上下文），必须复制结点
	out := cow.newNode()
	// 复制 items
	if cap(out.items) >= len(n.items) {
		out.items = out.items[:len(n.items)]
	} else {
		out.items = make(items[T], len(n.items), cap(n.items))
	}
	copy(out.items, n.items)
	// 复制 children
	if cap(out.children) >= len(n.children) {
		out.children = out.children[:len(n.children)]
	} else {
		out.children = make(items[*node[T]], len(n.children), cap(n.children))
	}
	copy(out.children, n.children)
	return out
}

// mutableChild 获取并返回第 i 个子节点的可写副本。
func (n *node[T]) mutableChild(i int) *node[T] {
	c := n.children[i].mutableFor(n.cow)
	n.children[i] = c
	return c
}

// split 在给定索引处对节点进行拆分。当前节点缩小，返回拆分出的 item 和新节点。
func (n *node[T]) split(i int) (T, *node[T]) {
	item := n.items[i]
	next := n.cow.newNode()
	next.items = append(next.items, n.items[i+1:]...)
	n.items.truncate(i)
	if len(n.children) > 0 {
		next.children = append(next.children, n.children[i+1:]...)
		n.children.truncate(i + 1)
	}
	return item, next
}

// maybeSplitChild 检查并在需要时拆分子节点。若拆分发生则返回 true。
func (n *node[T]) maybeSplitChild(i, maxItems int) bool {
	if len(n.children[i].items) < maxItems {
		return false
	}
	first := n.mutableChild(i)
	item, second := first.split(maxItems / 2)
	n.items.insertAt(i, item)
	n.children.insertAt(i+1, second)
	return true
}

// 递归插入。
// insert 向以当前节点为根的子树中插入一个元素，确保没有节点会超过 maxItems 个元素。
// 如果插入发现已有相等元素，则替换并返回旧元素与是否找到。
func (n *node[T]) insert(item T, maxItems int) (_ T, _ bool) {
	i, found := n.items.find(item, n.cow.less)
	if found {
		out := n.items[i]
		n.items[i] = item
		return out, true
	}
	if len(n.children) == 0 {
		n.items.insertAt(i, item)
		return
	}
	if n.maybeSplitChild(i, maxItems) {
		inTree := n.items[i]
		switch {
		case n.cow.less(item, inTree):
			// 不变，去第一个拆分出的节点
		case n.cow.less(inTree, item):
			i++ // 去第二个拆分出的节点
		default:
			out := n.items[i]
			n.items[i] = item
			return out, true
		}
	}
	return n.mutableChild(i).insert(item, maxItems)
}

// get 在子树中查找给定 key 并返回。
func (n *node[T]) get(key T) (_ T, _ bool) {
	i, found := n.items.find(key, n.cow.less)
	if found {
		return n.items[i], true
	} else if len(n.children) > 0 {
		return n.children[i].get(key)
	}
	return
}

// min 返回子树中的最小元素。
func min[T any](n *node[T]) (_ T, found bool) {
	if n == nil {
		return
	}
	for len(n.children) > 0 {
		n = n.children[0]
	}
	if len(n.items) == 0 {
		return
	}
	return n.items[0], true
}

// max 返回子树中的最大元素。
func max[T any](n *node[T]) (_ T, found bool) {
	if n == nil {
		return
	}
	for len(n.children) > 0 {
		n = n.children[len(n.children)-1]
	}
	if len(n.items) == 0 {
		return
	}
	return n.items[len(n.items)-1], true
}

// toRemove 表示在 node.remove 调用时要删除何种元素。
type toRemove int

const (
	removeItem toRemove = iota // 删除指定元素
	removeMin                  // 删除子树中的最小元素
	removeMax                  // 删除子树中的最大元素
)

// remove 从以当前节点为根的子树中删除一个元素。
func (n *node[T]) remove(item T, minItems int, typ toRemove) (_ T, _ bool) {
	var i int
	var found bool
	switch typ {
	case removeMax:
		if len(n.children) == 0 {
			return n.items.pop(), true
		}
		i = len(n.items)
	case removeMin:
		if len(n.children) == 0 {
			return n.items.removeAt(0), true
		}
		i = 0
	case removeItem:
		i, found = n.items.find(item, n.cow.less)
		if len(n.children) == 0 {
			if found {
				return n.items.removeAt(i), true
			}
			return
		}
	default:
		panic("invalid type")
	}
	// 如果到这里说明还有子节点
	if len(n.children[i].items) <= minItems {
		return n.growChildAndRemove(i, item, minItems, typ)
	}
	child := n.mutableChild(i)
	// 此时要么子节点足够大，要么做过借用/合并
	if found {
		// 说明在当前节点找到该元素，需用左子树中的前驱替换它
		out := n.items[i]
		var zero T
		n.items[i], _ = child.remove(zero, minItems, removeMax)
		return out, true
	}
	// 否则目标在子节点且子节点大小足够
	return child.remove(item, minItems, typ)
}

// growChildAndRemove 用于在移除元素前，对子节点进行扩充确保其大小 >= minItems，
// 然后再执行 remove。
func (n *node[T]) growChildAndRemove(i int, item T, minItems int, typ toRemove) (T, bool) {
	if i > 0 && len(n.children[i-1].items) > minItems {
		// 向左兄弟借一个元素
		child := n.mutableChild(i)
		stealFrom := n.mutableChild(i - 1)
		stolenItem := stealFrom.items.pop()
		child.items.insertAt(0, n.items[i-1])
		n.items[i-1] = stolenItem
		if len(stealFrom.children) > 0 {
			child.children.insertAt(0, stealFrom.children.pop())
		}
	} else if i < len(n.items) && len(n.children[i+1].items) > minItems {
		// 向右兄弟借一个元素
		child := n.mutableChild(i)
		stealFrom := n.mutableChild(i + 1)
		stolenItem := stealFrom.items.removeAt(0)
		child.items = append(child.items, n.items[i])
		n.items[i] = stolenItem
		if len(stealFrom.children) > 0 {
			child.children = append(child.children, stealFrom.children.removeAt(0))
		}
	} else {
		if i >= len(n.items) {
			i--
		}
		child := n.mutableChild(i)
		// 与右侧子节点合并
		mergeItem := n.items.removeAt(i)
		mergeChild := n.children.removeAt(i + 1)
		child.items = append(child.items, mergeItem)
		child.items = append(child.items, mergeChild.items...)
		child.children = append(child.children, mergeChild.children...)
		n.cow.freeNode(mergeChild)
	}
	return n.remove(item, minItems, typ)
}

type direction int

const (
	descend = direction(-1)
	ascend  = direction(+1)
)

type optionalItem[T any] struct {
	item  T
	valid bool
}

func optional[T any](item T) optionalItem[T] {
	return optionalItem[T]{item: item, valid: true}
}
func empty[T any]() optionalItem[T] {
	return optionalItem[T]{}
}

// iterate 提供了对树中元素的遍历方法。
//
// 当 ascending 时，start < stop；descending 时，start > stop。
// includeStart = true 表示当等于 start 时也包括在内，即实现一种 “>=” 或 “<=” 的逻辑。
// hit 表示是否已开始输出符合条件的元素。
func (n *node[T]) iterate(dir direction, start, stop optionalItem[T], includeStart bool, hit bool, iter ItemIteratorG[T]) (bool, bool) {
	var ok, found bool
	var index int
	switch dir {
	case ascend:
		if start.valid {
			index, _ = n.items.find(start.item, n.cow.less)
		}
		for i := index; i < len(n.items); i++ {
			if len(n.children) > 0 {
				if hit, ok = n.children[i].iterate(dir, start, stop, includeStart, hit, iter); !ok {
					return hit, false
				}
			}
			if !includeStart && !hit && start.valid && !n.cow.less(start.item, n.items[i]) {
				hit = true
				continue
			}
			hit = true
			if stop.valid && !n.cow.less(n.items[i], stop.item) {
				return hit, false
			}
			if !iter(n.items[i]) {
				return hit, false
			}
		}
		if len(n.children) > 0 {
			if hit, ok = n.children[len(n.children)-1].iterate(dir, start, stop, includeStart, hit, iter); !ok {
				return hit, false
			}
		}
	case descend:
		if start.valid {
			index, found = n.items.find(start.item, n.cow.less)
			if !found {
				index = index - 1
			}
		} else {
			index = len(n.items) - 1
		}
		for i := index; i >= 0; i-- {
			if start.valid && !n.cow.less(n.items[i], start.item) {
				if !includeStart || hit || n.cow.less(start.item, n.items[i]) {
					continue
				}
			}
			if len(n.children) > 0 {
				if hit, ok = n.children[i+1].iterate(dir, start, stop, includeStart, hit, iter); !ok {
					return hit, false
				}
			}
			if stop.valid && !n.cow.less(stop.item, n.items[i]) {
				return hit, false
			}
			hit = true
			if !iter(n.items[i]) {
				return hit, false
			}
		}
		if len(n.children) > 0 {
			if hit, ok = n.children[0].iterate(dir, start, stop, includeStart, hit, iter); !ok {
				return hit, false
			}
		}
	}
	return hit, true
}

// print 用于测试/调试，打印节点内容。
func (n *node[T]) print(w io.Writer, level int) {
	fmt.Fprintf(w, "%sNODE:%v\n", strings.Repeat("  ", level), n.items)
	for _, c := range n.children {
		c.print(w, level+1)
	}
}

// BTreeG 是一个泛型 B-Tree 的实现。维护 B-Tree 的 根节点 和 全局信息。
//
// BTreeG 存储类型为 T 的元素，以有序结构的方式支持插入、删除和遍历操作。
//
// 多协程并发读是安全的，但写操作不支持多协程同时进行。
type BTreeG[T any] struct {
	degree int // 树的阶，决定每个节点最少/最多能含有多少元素.[degree-1, 2*degree-1]
	length int // 整棵树的元素数量
	root   *node[T]
	cow    *copyOnWriteContext[T] // 这棵树的 copy-on-write 上下文。任何写操作都需要检查节点是否归属于该上下文
}

// LessFunc[T] 用于决定类型 T 的排序方式。对于 a < b 的情况返回 true。
type LessFunc[T any] func(a, b T) bool

// copyOnWriteContext 用于确定节点的所有权……
//
// 当一个树和节点的 write context 相同，则可直接修改该节点；
// 如果不同，必须先复制节点并更新其 cow。
type copyOnWriteContext[T any] struct {
	freelist *FreeListG[T]
	less     LessFunc[T]
}

// B-Tree 在做克隆之后，原树和新树可以共享同一份节点。但一旦有写操作，就要复制节点并设置新的上下文，以免互相影响
// Clone 会对 btree 进行惰性克隆。Clone 本身不应在并发场景下调用，
// 但调用完成后，原树 t 与新树 t2 可以并发使用。
//
// 原树的内部结构会被标记为只读并与新树共享。
// 对 t 和 t2 的写操作会采用写时复制逻辑，只有在修改原树节点时才创建新副本。
// 读操作不会有性能下降。
// 两棵树的写操作在最初会因为额外的分配和复制而稍有减速，但之后会逐渐恢复到原先水平。
//
// !树的节点没有立即复制，只有当后续对某个节点做写操作时，发现 n.cow != 当前树.cow，才会真正去复制节点
func (t *BTreeG[T]) Clone() (t2 *BTreeG[T]) {
	cow1, cow2 := *t.cow, *t.cow
	out := *t

	// !给原树 t 和 新树 out 分配不同的 COW 上下文
	t.cow = &cow1
	out.cow = &cow2
	return &out
}

// maxItems 返回节点可容纳的最大元素数。
func (t *BTreeG[T]) maxItems() int {
	return t.degree*2 - 1
}

// minItems 返回节点可容纳的最少元素数（根节点可不遵守此限制）。
func (t *BTreeG[T]) minItems() int {
	return t.degree - 1
}

func (c *copyOnWriteContext[T]) newNode() (n *node[T]) {
	n = c.freelist.newNode()
	n.cow = c
	return
}

type freeType int

const (
	ftFreelistFull freeType = iota // freelist 满了，节点无法放进去，只能交给 Go 的 GC 回收
	ftStored                       // 节点成功放入 freelist
	ftNotOwned                     // 节点并不属于当前 COW，上下文无权释放。表示“这不是我管辖的节点，我不处理它”。
)

// freeNode 会在给定 COW 上下文中释放一个节点（若归其所有）。
// 返回值指示了节点最终状态（见 freeType）。
func (c *copyOnWriteContext[T]) freeNode(n *node[T]) freeType {
	if n.cow == c {
		// 1. 清空节点内容，以便 GC
		n.items.truncate(0)
		n.children.truncate(0)
		n.cow = nil

		// 2. 尝试放入 freelist
		if c.freelist.freeNode(n) {
			return ftStored
		} else {
			return ftFreelistFull
		}
	} else {
		return ftNotOwned
	}
}

// ReplaceOrInsert 向树中添加给定元素 item。
// 如果树中已存在等价元素，则将其移除并返回，同时返回 true 表示替换成功；
// 否则返回 (零值, false)。
//
// 若插入 nil 会 panic。
func (t *BTreeG[T]) ReplaceOrInsert(item T) (_ T, _ bool) {
	if t.root == nil {
		t.root = t.cow.newNode()
		t.root.items = append(t.root.items, item)
		t.length++
		return
	} else {
		t.root = t.root.mutableFor(t.cow)
		if len(t.root.items) >= t.maxItems() {
			item2, second := t.root.split(t.maxItems() / 2)
			oldroot := t.root
			t.root = t.cow.newNode()
			t.root.items = append(t.root.items, item2)
			t.root.children = append(t.root.children, oldroot, second)
		}
	}
	out, outb := t.root.insert(item, t.maxItems())
	if !outb {
		t.length++
	}
	return out, outb
}

// Delete 从树中删除一个与传入 item 相等的元素并返回。若未找到则返回 (零值, false)。
func (t *BTreeG[T]) Delete(item T) (T, bool) {
	return t.deleteItem(item, removeItem)
}

// DeleteMin 删除树中的最小元素并返回。若树为空则返回 (零值, false)。
func (t *BTreeG[T]) DeleteMin() (T, bool) {
	var zero T
	return t.deleteItem(zero, removeMin)
}

// DeleteMax 删除树中的最大元素并返回。若树为空则返回 (零值, false)。
func (t *BTreeG[T]) DeleteMax() (T, bool) {
	var zero T
	return t.deleteItem(zero, removeMax)
}

func (t *BTreeG[T]) deleteItem(item T, typ toRemove) (_ T, _ bool) {
	if t.root == nil || len(t.root.items) == 0 {
		return
	}
	t.root = t.root.mutableFor(t.cow)
	out, outb := t.root.remove(item, t.minItems(), typ)
	if len(t.root.items) == 0 && len(t.root.children) > 0 {
		oldroot := t.root
		t.root = t.root.children[0]
		t.cow.freeNode(oldroot)
	}
	if outb {
		t.length--
	}
	return out, outb
}

// AscendRange 对 [greaterOrEqual, lessThan) 范围内的每个元素调用 iterator，
// 当 iterator 返回 false 时停止。
func (t *BTreeG[T]) AscendRange(greaterOrEqual, lessThan T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, optional[T](greaterOrEqual), optional[T](lessThan), true, false, iterator)
}

// AscendLessThan 对 [first, pivot) 范围内的每个元素调用 iterator，
// 当 iterator 返回 false 时停止。
func (t *BTreeG[T]) AscendLessThan(pivot T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, empty[T](), optional(pivot), false, false, iterator)
}

// AscendGreaterOrEqual 对 [pivot, last] 范围内的每个元素调用 iterator，
// 当 iterator 返回 false 时停止。
func (t *BTreeG[T]) AscendGreaterOrEqual(pivot T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, optional[T](pivot), empty[T](), true, false, iterator)
}

// Ascend 调用 iterator 遍历 [first, last] 范围内的每个元素，
// 当 iterator 返回 false 时停止。
func (t *BTreeG[T]) Ascend(iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, empty[T](), empty[T](), false, false, iterator)
}

// DescendRange 对 [lessOrEqual, greaterThan) 范围内的每个元素调用 iterator，
// 当 iterator 返回 false 时停止。
func (t *BTreeG[T]) DescendRange(lessOrEqual, greaterThan T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, optional[T](lessOrEqual), optional[T](greaterThan), true, false, iterator)
}

// DescendLessOrEqual 对 [pivot, first] 范围内的每个元素调用 iterator，
// 当 iterator 返回 false 时停止。
func (t *BTreeG[T]) DescendLessOrEqual(pivot T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, optional[T](pivot), empty[T](), true, false, iterator)
}

// DescendGreaterThan 对 [last, pivot) 范围内的每个元素调用 iterator，
// 当 iterator 返回 false 时停止。
func (t *BTreeG[T]) DescendGreaterThan(pivot T, iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, empty[T](), optional[T](pivot), false, false, iterator)
}

// Descend 调用 iterator 遍历 [last, first] 范围内的每个元素，
// 当 iterator 返回 false 时停止。
func (t *BTreeG[T]) Descend(iterator ItemIteratorG[T]) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, empty[T](), empty[T](), false, false, iterator)
}

// Get 查找树中与 key 相等的元素并返回；未找到则返回 (零值, false)。
func (t *BTreeG[T]) Get(key T) (_ T, _ bool) {
	if t.root == nil {
		return
	}
	return t.root.get(key)
}

// Min 返回树中最小的元素；若树为空，返回 (零值, false)。
func (t *BTreeG[T]) Min() (_ T, _ bool) {
	return min(t.root)
}

// Max 返回树中最大的元素；若树为空，返回 (零值, false)。
func (t *BTreeG[T]) Max() (_ T, _ bool) {
	return max(t.root)
}

// Has 若树中包含 key 则返回 true，否则返回 false。
func (t *BTreeG[T]) Has(key T) bool {
	_, ok := t.Get(key)
	return ok
}

// Len 返回树中当前的元素数量。
func (t *BTreeG[T]) Len() int {
	return t.length
}

// 允许你快速清空整棵树，比逐条删除效率更高，还能把节点放回 freelist（若选 addNodesToFreelist=true）
//
// Clear 清空 btree 中的所有元素。如果 addNodesToFreelist 为 true，
// 会尝试将节点加入 freelist（直到满），否则直接将 root 置空，
// 由 Go 的 GC 处理。
//
// 此操作可能比对所有元素依次 Delete 效率更高，
// 因为后者需逐个查找并调整树结构。
// 也比新建一棵替换旧树略快一些，因为旧树的节点会被回收进 freelist，
// 而不是丢给 GC。
//
// 该操作耗时与以下因素相关：
//   - O(1)：当 addNodesToFreelist=false 时，只是一次性操作。
//   - O(1)：当 freelist 已满时，会立即停止。
//   - O(空闲列表容量)：当 freelist 未满且所有节点都属于本树时，会将节点加入 freelist 直到其满。
//   - O(树大小)：当节点属于其他树时，需要遍历所有节点，但由于所有权不匹配，无法加入 freelist。
func (t *BTreeG[T]) Clear(addNodesToFreelist bool) {
	if t.root != nil && addNodesToFreelist {
		t.root.reset(t.cow)
	}
	t.root, t.length = nil, 0
}

// reset 会递归地把子节点也释放到 freelist，如果没满的话
// reset 将子树的节点释放到 freelist。如果 freelist 已满会立刻结束。
// 若返回 false，表示可以直接结束上层的 reset 调用。
func (n *node[T]) reset(c *copyOnWriteContext[T]) bool {
	for _, child := range n.children {
		if !child.reset(c) {
			return false
		}
	}
	return c.freeNode(n) != ftFreelistFull
}
