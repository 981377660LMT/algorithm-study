以下是对提供的 Go 语言跳表（Skiplist）实现代码的翻译和逐行详细分析。

## 概述

该代码实现了一个跳表（Skiplist），用于以递增顺序存储元素。跳表允许在大约 \( O(\log n) \) 的时间复杂度内进行查找、插入和删除操作。此外，还提供了获取下一个和上一个元素以及在不需要重新插入的情况下更改实际值的方法（前提是键保持不变）。跳表是平衡树的快速替代方案。

该实现基于 [huandu/skiplist](https://github.com/huandu/skiplist) 仓库，并遵循 MIT 许可证。

---

## 许可证声明

```go
// MIT License
//
// Copyright (c) 2018 Maurice Tollmien [...]

// [许可证全文省略]
```

**翻译：**

该部分是 MIT 许可证的声明，允许任何人免费使用、复制、修改、合并、发布、分发、再授权和/或销售软件副本，同时须保留版权声明和许可证通知。软件按“原样”提供，不附带任何明示或暗示的担保。

---

## 包和导入

```go
// Package skiplist is an implementation of a skiplist to store elements in increasing order.
// It allows finding, insertion and deletion operations in approximately O(n log(n)).
// Additionally, there are methods for retrieving the next and previous element as well as changing the actual value
// without the need for re-insertion (as long as the key stays the same!)
// Skiplist is a fast alternative to a balanced tree.

package skiplist

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"time"
)
```

**翻译及分析：**

- **包声明：** `package skiplist` 定义了一个名为 `skiplist` 的包，用于实现跳表数据结构。
- **导入的包：**
  - `fmt`：格式化输入输出。
  - `math`：数学函数。
  - `math/bits`：位操作函数。
  - `math/rand`：伪随机数生成。
  - `time`：时间相关函数。

这些导入的包为跳表的实现提供必要的功能支持。

---

## 常量定义

```go
const (
	// maxLevel denotes the maximum height of the skiplist. This height will keep the skiplist
	// efficient for up to 34m entries. If there is a need for much more, please adjust this constant accordingly.
	maxLevel = 25
	eps      = 0.00001
)
```

**翻译及分析：**

- `maxLevel = 25`：定义跳表的最大层数为 25。这一高度使得跳表对于最多 3400 万（34m）条目仍然高效。如果需要存储更多元素，可以适当调整此常量。
- `eps = 0.00001`：定义一个极小值（ε），用于在比较键值时判断两个键是否相等。由于浮点数计算可能存在精度问题，使用 `eps` 来确定键值的近似相等。

---

## 接口与结构体定义

### ListElement 接口

```go
// ListElement is the interface to implement for elements that are inserted into the skiplist.
type ListElement interface {
	// ExtractKey() returns a float64 representation of the key that is used for insertion/deletion/find. It needs to establish an order over all elements
	ExtractKey() float64
	// A string representation of the element. Can be used for pretty-printing the list. Otherwise just return an empty string.
	String() string
}
```

**翻译及分析：**

- **ListElement 接口：** 定义了插入到跳表中的元素需要实现的接口。包括两个方法：
  - `ExtractKey() float64`：返回用于插入、删除和查找的键的 `float64` 表示。这个键需要在所有元素之间建立一个有序关系。
  - `String() string`：返回元素的字符串表示，主要用于美观地打印跳表。如果不需要打印，可以返回空字符串。

### SkipListElement 结构体

```go
// SkipListElement represents one actual Node in the skiplist structure.
// It saves the actual element, pointers to the next nodes and a pointer to one previous node.
type SkipListElement struct {
	next  [maxLevel]*SkipListElement
	level int
	key   float64
	value ListElement
	prev  *SkipListElement
}
```

**翻译及分析：**

- **SkipListElement 结构体：** 表示跳表结构中的一个实际节点。包含以下字段：
  - `next [maxLevel]*SkipListElement`：一个长度为 `maxLevel` 的数组，存储指向各层下一个节点的指针。
  - `level int`：当前节点的层数。
  - `key float64`：节点的键值，用于排序。
  - `value ListElement`：节点存储的实际元素，实现了 `ListElement` 接口。
  - `prev *SkipListElement`：指向前一个节点的指针，仅在底层（第 0 层）中使用，便于双向遍历。

### SkipList 结构体

```go
// SkipList is the actual skiplist representation.
// It saves all nodes accessible from the start and end and keeps track of element count, eps and levels.
type SkipList struct {
	startLevels  [maxLevel]*SkipListElement
	endLevels    [maxLevel]*SkipListElement
	maxNewLevel  int
	maxLevel     int
	elementCount int
	eps          float64
}
```

**翻译及分析：**

- **SkipList 结构体：** 跳表的实际表示。包含以下字段：
  - `startLevels [maxLevel]*SkipListElement`：一个长度为 `maxLevel` 的数组，存储每一层的起始节点指针。
  - `endLevels [maxLevel]*SkipListElement`：一个长度为 `maxLevel` 的数组，存储每一层的末尾节点指针。
  - `maxNewLevel int`：用于生成新节点层数的当前最大层数。
  - `maxLevel int`：当前跳表的最大层数。
  - `elementCount int`：跳表中元素的计数。
  - `eps float64`：用于键值比较的极小值，确保浮点数比较的准确性。

---

## 构造函数

### NewSeedEps

```go
// NewSeedEps returns a new empty, initialized Skiplist.
// Given a seed, a deterministic height/list behaviour can be achieved.
// Eps is used to compare keys given by the ExtractKey() function on equality.
func NewSeedEps(seed int64, eps float64) SkipList {

	// Initialize random number generator.
	rand.Seed(seed)
	//fmt.Printf("SkipList seed: %v\n", seed)

	list := SkipList{
		startLevels:  [maxLevel]*SkipListElement{},
		endLevels:    [maxLevel]*SkipListElement{},
		maxNewLevel:  maxLevel,
		maxLevel:     0,
		elementCount: 0,
		eps:          eps,
	}

	return list
}
```

**翻译及分析：**

- **NewSeedEps 函数：** 返回一个新的、初始化的空跳表。
  - **参数：**
    - `seed int64`：用于随机数生成器的种子，从而实现确定性高度/列表行为。
    - `eps float64`：用于比较通过 `ExtractKey()` 函数提供的键值是否相等。
  - **步骤：**
    1. 使用 `rand.Seed(seed)` 初始化随机数生成器，以确保生成的层数是可预测的。
    2. 创建一个 `SkipList` 实例，初始化所有字段：
       - `startLevels` 和 `endLevels` 数组中的所有指针都设为 `nil`。
       - `maxNewLevel` 设为 `maxLevel`。
       - `maxLevel` 设为 `0`，表示当前跳表没有任何元素。
       - `elementCount` 设为 `0`，表示跳表中元素数量。
       - `eps` 设为传入的 `eps` 参数。
    3. 返回初始化后的跳表。

### NewEps

```go
// NewEps returns a new empty, initialized Skiplist.
// Eps is used to compare keys given by the ExtractKey() function on equality.
func NewEps(eps float64) SkipList {
	return NewSeedEps(time.Now().UTC().UnixNano(), eps)
}
```

**翻译及分析：**

- **NewEps 函数：** 返回一个新的、初始化的空跳表。
  - **参数：**
    - `eps float64`：用于键值比较。
  - **步骤：**
    1. 调用 `NewSeedEps`，使用当前 UTC 时间的纳秒数作为种子，确保跳表的层数生成具有随机性。
    2. 返回初始化后的跳表。

### NewSeed

```go
// NewSeed returns a new empty, initialized Skiplist.
// Given a seed, a deterministic height/list behaviour can be achieved.
func NewSeed(seed int64) SkipList {
	return NewSeedEps(seed, eps)
}
```

**翻译及分析：**

- **NewSeed 函数：** 返回一个新的、初始化的空跳表。
  - **参数：**
    - `seed int64`：用于随机数生成器的种子。
  - **步骤：**
    1. 调用 `NewSeedEps`，使用提供的种子和默认的 `eps` 常量。
    2. 返回初始化后的跳表。

### New

```go
// New returns a new empty, initialized Skiplist.
func New() SkipList {
	return NewSeedEps(time.Now().UTC().UnixNano(), eps)
}
```

**翻译及分析：**

- **New 函数：** 返回一个新的、初始化的空跳表。
  - **步骤：**
    1. 调用 `NewSeedEps`，使用当前 UTC 时间的纳秒数作为种子和默认的 `eps` 常量。
    2. 返回初始化后的跳表。

---

## 基本操作

### IsEmpty

```go
// IsEmpty checks, if the skiplist is empty.
func (t *SkipList) IsEmpty() bool {
	return t.startLevels[0] == nil
}
```

**翻译及分析：**

- **IsEmpty 方法：** 检查跳表是否为空。
  - **接收者：** 指向 `SkipList` 的指针 `t`。
  - **返回值：** 如果第 0 层的起始节点为 `nil`，则返回 `true`，表示跳表为空；否则返回 `false`。

### generateLevel

```go
func (t *SkipList) generateLevel(maxLevel int) int {
	level := maxLevel - 1
	// First we apply some mask which makes sure that we don't get a level
	// above our desired level. Then we find the first set bit.
	var x uint64 = rand.Uint64() & ((1 << uint(maxLevel-1)) - 1)
	zeroes := bits.TrailingZeros64(x)
	if zeroes <= maxLevel {
		level = zeroes
	}

	return level
}
```

**翻译及分析：**

- **generateLevel 方法：** 生成一个随机层数，用于新插入节点的层高。
  - **参数：** `maxLevel int` 指定生成层数的上限。
  - **步骤：**
    1. 初始层数 `level` 设为 `maxLevel - 1`。
    2. 生成一个 `uint64` 随机数 `x`，并通过掩码确保其不超过 `maxLevel - 1`。具体而言，掩码 `(1 << (maxLevel-1)) - 1` 生成一个具有 `maxLevel-1` 个二进制位的掩码。
    3. 使用 `bits.TrailingZeros64(x)` 计算 `x` 的最低有效位从右边数起的连续零的个数，即找到第一个设置为 1 的位的位置（从 0 开始）。
    4. 如果连续零的个数 `zeroes` 小于或等于 `maxLevel`，则将 `level` 设为 `zeroes`。
    5. 返回生成的 `level`。

**补充说明：**

这种层数生成方法基于概率，期望生成的层数遵循几何分布，从而确保跳表的层数与元素数量成对数关系，保持高效性。

### findEntryIndex

```go
func (t *SkipList) findEntryIndex(key float64, level int) int {
	// Find good entry point so we don't accidentally skip half the list.
	for i := t.maxLevel; i >= 0; i-- {
		if t.startLevels[i] != nil && t.startLevels[i].key <= key || i <= level {
			return i
		}
	}
	return 0
}
```

**翻译及分析：**

- **findEntryIndex 方法：** 查找适合开始搜索的层数索引，以避免跳过列表的一半。
  - **参数：**
    - `key float64`：要查找的键。
    - `level int`：最小层数。
  - **步骤：**
    1. 从当前跳表的最大层 `t.maxLevel` 开始，向下遍历层数。
    2. 对于每一层 `i`，检查：
       - 如果该层的起始节点不为 `nil` 且其键值小于或等于 `key`，或者当前层数 `i` 小于等于提供的 `level`，则返回该层的索引 `i`。
    3. 如果没有满足条件的层，则返回 0。

**补充说明：**

此方法用于优化查找路径，确保从适当的层数开始搜索，避免不必要的遍历。

### findExtended

```go
func (t *SkipList) findExtended(key float64, findGreaterOrEqual bool) (foundElem *SkipListElement, ok bool) {

	foundElem = nil
	ok = false

	if t.IsEmpty() {
		return
	}

	index := t.findEntryIndex(key, 0)
	var currentNode *SkipListElement

	currentNode = t.startLevels[index]
	nextNode := currentNode

	// In case, that our first element is already greater-or-equal!
	if findGreaterOrEqual && currentNode.key > key {
		foundElem = currentNode
		ok = true
		return
	}

	for {
		if math.Abs(currentNode.key-key) <= t.eps {
			foundElem = currentNode
			ok = true
			return
		}

		nextNode = currentNode.next[index]

		// Which direction are we continuing next time?
		if nextNode != nil && nextNode.key <= key {
			// Go right
			currentNode = nextNode
		} else {
			if index > 0 {

				// Early exit
				if currentNode.next[0] != nil && math.Abs(currentNode.next[0].key-key) <= t.eps {
					foundElem = currentNode.next[0]
					ok = true
					return
				}
				// Go down
				index--
			} else {
				// Element is not found and we reached the bottom.
				if findGreaterOrEqual {
					foundElem = nextNode
					ok = nextNode != nil
				}

				return
			}
		}
	}
}
```

**翻译及分析：**

- **findExtended 方法：** 在跳表中查找特定键的元素，或找到第一个大于或等于该键的元素。
  - **参数：**
    - `key float64`：要查找的键。
    - `findGreaterOrEqual bool`：是否寻找第一个大于或等于该键的元素。
  - **返回值：**
    - `foundElem *SkipListElement`：找到的元素节点指针。
    - `ok bool`：是否成功找到元素。
  - **步骤：**
    1. 初始化 `foundElem` 为 `nil`，`ok` 为 `false`。
    2. 如果跳表为空，直接返回。
    3. 使用 `findEntryIndex` 方法找到开始搜索的层数索引 `index`。
    4. 将 `currentNode` 初始化为该层的起始节点。
    5. 如果 `findGreaterOrEqual` 为 `true` 且 `currentNode.key > key`，则当前节点即为所需节点，返回。
    6. 进入循环，执行以下操作：
       - 如果 `currentNode.key` 与 `key` 的绝对差值小于等于 `eps`，则找到目标节点，返回。
       - 获取当前节点在当前层的下一个节点 `nextNode`。
       - 如果 `nextNode` 不为 `nil` 且其键值小于或等于 `key`，则向右移动到 `nextNode`。
       - 否则：
         - 如果当前层数 `index` 大于 0，尝试向下移动一层：
           - 进行提前退出检查：如果 `currentNode.next[0]` 不为 `nil` 且其键值与 `key` 近似相等，则找到目标节点，返回。
           - 向下移动一层 (`index--`)。
         - 如果已经在底层 (`index == 0`)，则：
           - 如果 `findGreaterOrEqual` 为 `true`，将 `foundElem` 设为 `nextNode`，并根据 `nextNode` 是否存在设置 `ok`。
           - 返回。

**补充说明：**

此方法是跳表查找的核心，通过层级跳跃快速定位目标元素或其相邻元素，提高查找效率。

---

## 查找操作

### Find

```go
// Find tries to find an element in the skiplist based on the key from the given ListElement.
// elem can be used, if ok is true.
// Find runs in approx. O(log(n))
func (t *SkipList) Find(e ListElement) (elem *SkipListElement, ok bool) {

	if t == nil || e == nil {
		return
	}

	elem, ok = t.findExtended(e.ExtractKey(), false)
	return
}
```

**翻译及分析：**

- **Find 方法：** 根据给定 `ListElement` 的键在跳表中查找元素。
  - **参数：** `e ListElement`：要查找的元素，需实现 `ListElement` 接口。
  - **返回值：**
    - `elem *SkipListElement`：找到的节点指针。
    - `ok bool`：是否成功找到元素。
  - **步骤：**
    1. 如果跳表或者元素 `e` 为 `nil`，则直接返回。
    2. 调用 `findExtended` 方法，传入 `e.ExtractKey()` 作为键，`false` 表示不需要查找大于或等于的元素，只查找等于的元素。
    3. 返回 `findExtended` 的结果。

**时间复杂度：** 约为 \( O(\log n) \)。

### FindGreaterOrEqual

```go
// FindGreaterOrEqual finds the first element, that is greater or equal to the given ListElement e.
// The comparison is done on the keys (So on ExtractKey()).
// FindGreaterOrEqual runs in approx. O(log(n))
func (t *SkipList) FindGreaterOrEqual(e ListElement) (elem *SkipListElement, ok bool) {

	if t == nil || e == nil {
		return
	}

	elem, ok = t.findExtended(e.ExtractKey(), true)
	return
}
```

**翻译及分析：**

- **FindGreaterOrEqual 方法：** 查找第一个大于或等于给定 `ListElement` 的元素。
  - **参数：** `e ListElement`：参考元素，需实现 `ListElement` 接口。
  - **返回值：**
    - `elem *SkipListElement`：找到的节点指针。
    - `ok bool`：是否成功找到元素。
  - **步骤：**
    1. 如果跳表或者元素 `e` 为 `nil`，则直接返回。
    2. 调用 `findExtended` 方法，传入 `e.ExtractKey()` 作为键，`true` 表示需要查找大于或等于的元素。
    3. 返回 `findExtended` 的结果。

**时间复杂度：** 约为 \( O(\log n) \)。

---

## 删除操作

### Delete

```go
// Delete removes an element equal to e from the skiplist, if there is one.
// If there are multiple entries with the same value, Delete will remove one of them
// (Which one will change based on the actual skiplist layout)
// Delete runs in approx. O(log(n))
func (t *SkipList) Delete(e ListElement) {

	if t == nil || t.IsEmpty() || e == nil {
		return
	}

	key := e.ExtractKey()

	index := t.findEntryIndex(key, 0)

	var currentNode *SkipListElement
	nextNode := currentNode

	for {

		if currentNode == nil {
			nextNode = t.startLevels[index]
		} else {
			nextNode = currentNode.next[index]
		}

		// Found and remove!
		if nextNode != nil && math.Abs(nextNode.key-key) <= t.eps {

			if currentNode != nil {
				currentNode.next[index] = nextNode.next[index]
			}

			if index == 0 {
				if nextNode.next[index] != nil {
					nextNode.next[index].prev = currentNode
				}
				t.elementCount--
			}

			// Link from start needs readjustments.
			if t.startLevels[index] == nextNode {
				t.startLevels[index] = nextNode.next[index]
				// This was our currently highest node!
				if t.startLevels[index] == nil {
					t.maxLevel = index - 1
				}
			}

			// Link from end needs readjustments.
			if nextNode.next[index] == nil {
				t.endLevels[index] = currentNode
			}
			nextNode.next[index] = nil
		}

		if nextNode != nil && nextNode.key < key {
			// Go right
			currentNode = nextNode
		} else {
			// Go down
			index--
			if index < 0 {
				break
			}
		}
	}

}
```

**翻译及分析：**

- **Delete 方法：** 从跳表中删除与给定 `ListElement` 相等的元素（基于键值）。如果有多个相同键值的节点，仅删除其中一个。
  - **参数：** `e ListElement`：要删除的元素，需实现 `ListElement` 接口。
  - **步骤：**
    1. 如果跳表为 `nil`、为空或元素 `e` 为 `nil`，则直接返回。
    2. 提取元素 `e` 的键值 `key`。
    3. 使用 `findEntryIndex` 方法找到开始搜索的层数索引 `index`。
    4. 初始化 `currentNode` 为 `nil`，`nextNode` 也为 `nil`。
    5. 进入循环，执行以下操作：
       - 如果 `currentNode` 为 `nil`，则将 `nextNode` 设为当前层的起始节点。
       - 否则，将 `nextNode` 设为 `currentNode` 在当前层的下一个节点。
       - 如果 `nextNode` 不为 `nil` 且其键值与 `key` 的绝对差值小于等于 `eps`，则找到目标节点进行删除：
         - 如果 `currentNode` 不为 `nil`，则将其在当前层的下一个指针指向 `nextNode.next[index]`，跳过 `nextNode`。
         - 如果当前层数为 0（底层），则：
           - 如果 `nextNode.next[index]` 不为 `nil`，将其 `prev` 指针指向 `currentNode`（维护双向链表）。
           - 元素计数 `elementCount` 减 1。
         - 如果当前层的起始节点为 `nextNode`，则将该层的起始节点更新为 `nextNode.next[index]`。
           - 如果更新后该层的起始节点为 `nil`，则减少跳表的最大层数 `maxLevel`。
         - 如果 `nextNode.next[index]` 为 `nil`，则将该层的末尾节点更新为 `currentNode`。
         - 将 `nextNode.next[index]` 设为 `nil`，断开与后续节点的连接。
       - 如果 `nextNode` 不为 `nil` 且其键值小于 `key`，则向右移动到 `nextNode`。
       - 否则，向下移动一层 (`index--`)。
       - 如果层数索引 `index` 小于 0，则退出循环。

**时间复杂度：** 约为 \( O(\log n) \)。

---

## 插入操作

### Insert

```go
// Insert inserts the given ListElement into the skiplist.
// Insert runs in approx. O(log(n))
func (t *SkipList) Insert(e ListElement) {

	if t == nil || e == nil {
		return
	}

	level := t.generateLevel(t.maxNewLevel)

	// Only grow the height of the skiplist by one at a time!
	if level > t.maxLevel {
		level = t.maxLevel + 1
		t.maxLevel = level
	}

	elem := &SkipListElement{
		next:  [maxLevel]*SkipListElement{},
		level: level,
		key:   e.ExtractKey(),
		value: e,
	}

	t.elementCount++

	newFirst := true
	newLast := true
	if !t.IsEmpty() {
		newFirst = elem.key < t.startLevels[0].key
		newLast = elem.key > t.endLevels[0].key
	}

	normallyInserted := false
	if !newFirst && !newLast {

		normallyInserted = true

		index := t.findEntryIndex(elem.key, level)

		var currentNode *SkipListElement
		nextNode := t.startLevels[index]

		for {

			if currentNode == nil {
				nextNode = t.startLevels[index]
			} else {
				nextNode = currentNode.next[index]
			}

			// Connect node to next
			if index <= level && (nextNode == nil || nextNode.key > elem.key) {
				elem.next[index] = nextNode
				if currentNode != nil {
					currentNode.next[index] = elem
				}
				if index == 0 {
					elem.prev = currentNode
					if nextNode != nil {
						nextNode.prev = elem
					}
				}
			}

			if nextNode != nil && nextNode.key <= elem.key {
				// Go right
				currentNode = nextNode
			} else {
				// Go down
				index--
				if index < 0 {
					break
				}
			}
		}
	}

	// Where we have a left-most position that needs to be referenced!
	for i := level; i >= 0; i-- {

		didSomething := false

		if newFirst || normallyInserted {

			if t.startLevels[i] == nil || t.startLevels[i].key > elem.key {
				if i == 0 && t.startLevels[i] != nil {
					t.startLevels[i].prev = elem
				}
				elem.next[i] = t.startLevels[i]
				t.startLevels[i] = elem
			}

			// link the endLevels to this element!
			if elem.next[i] == nil {
				t.endLevels[i] = elem
			}

			didSomething = true
		}

		if newLast {
			// Places the element after the very last element on this level!
			// This is very important, so we are not linking the very first element (newFirst AND newLast) to itself!
			if !newFirst {
				if t.endLevels[i] != nil {
					t.endLevels[i].next[i] = elem
				}
				if i == 0 {
					elem.prev = t.endLevels[i]
				}
				t.endLevels[i] = elem
			}

			// Link the startLevels to this element!
			if t.startLevels[i] == nil || t.startLevels[i].key > elem.key {
				t.startLevels[i] = elem
			}

			didSomething = true
		}

		if !didSomething {
			break
		}
	}
}
```

**翻译及分析：**

- **Insert 方法：** 将给定的 `ListElement` 插入到跳表中。
  - **参数：** `e ListElement`：要插入的元素，需实现 `ListElement` 接口。
  - **步骤：**
    1. 如果跳表为 `nil` 或元素 `e` 为 `nil`，则直接返回。
    2. 调用 `generateLevel` 方法生成新节点的层数 `level`。
    3. 检查生成的 `level` 是否超过当前跳表的最大层数 `maxLevel`：
       - 如果 `level > maxLevel`，则将 `level` 设为 `maxLevel + 1`，并更新 `maxLevel`。
    4. 创建一个新的 `SkipListElement` 实例 `elem`，初始化其字段：
       - `next` 数组全部设为 `nil`。
       - `level` 设为生成的层数。
       - `key` 从 `e.ExtractKey()` 获取。
       - `value` 设为 `e`。
    5. 增加元素计数 `elementCount`。
    6. 初始化 `newFirst` 和 `newLast` 为 `true`，表示可能是新的最小或最大元素。
       - 如果跳表不为空：
         - `newFirst` 设为 `elem.key < t.startLevels[0].key`，即新元素是否小于当前最小元素。
         - `newLast` 设为 `elem.key > t.endLevels[0].key`，即新元素是否大于当前最大元素。
    7. 初始化 `normallyInserted` 为 `false`。
       - 如果新元素既不是新的最小元素，也不是新的最大元素，则进行正常插入：
         - 设置 `normallyInserted = true`。
         - 使用 `findEntryIndex` 找到插入位置开始的层数 `index`。
         - 初始化 `currentNode` 为 `nil`，`nextNode` 为当前层的起始节点。
         - 进入循环，执行以下操作：
           - 如果 `currentNode` 为 `nil`，将 `nextNode` 设为当前层的起始节点；否则，将其设为 `currentNode` 在当前层的下一个节点。
           - 如果当前层数 `index` 小于等于新节点的层数 `level`，且 `nextNode` 为 `nil` 或其键值大于新节点的键值：
             - 将新节点的 `next[index]` 设为 `nextNode`。
             - 如果 `currentNode` 不为 `nil`，将其 `next[index]` 设为新节点 `elem`。
             - 如果当前层 `index == 0`，则更新新节点的 `prev` 指针为 `currentNode`，并将 `nextNode` 的 `prev` 指针设为新节点（如果 `nextNode` 不为 `nil`）。
           - 如果 `nextNode` 不为 `nil` 且其键值小于或等于新节点的键值，则向右移动到 `nextNode`。
           - 否则，向下移动一层 (`index--`)。
           - 如果层数索引 `index < 0`，则退出循环。
    8. 处理新节点作为最小或最大元素的情况：
       - 从新节点的层数 `level` 开始，向下遍历到第 0 层，执行以下操作：
         - 初始化 `didSomething` 为 `false`。
         - 如果 `newFirst` 或 `normallyInserted` 为 `true`：
           - 如果当前层的起始节点为 `nil` 或其键值大于新节点的键值：
             - 如果 `i == 0` 且当前层的起始节点不为 `nil`，则将当前层起始节点的 `prev` 指针设为新节点。
             - 将新节点的 `next[i]` 设为当前层的起始节点。
             - 将当前层的起始节点更新为新节点 `elem`。
           - 如果新节点的 `next[i]` 为 `nil`，则将当前层的末尾节点设为新节点。
           - 设置 `didSomething = true`。
         - 如果 `newLast` 为 `true`：
           - 如果 `!newFirst`（即新节点不是同时是新的最小和最大的节点），则：
             - 如果当前层的末尾节点不为 `nil`，将其 `next[i]` 设为新节点。
             - 如果 `i == 0`，则将新节点的 `prev` 指针设为当前层的末尾节点。
             - 更新当前层的末尾节点为新节点。
           - 如果当前层的起始节点为 `nil` 或其键值大于新节点的键值，则将当前层的起始节点设为新节点。
           - 设置 `didSomething = true`。
         - 如果 `didSomething` 为 `false`，则退出循环。

**时间复杂度：** 约为 \( O(\log n) \)。

---

## 其他操作

### GetValue

```go
// GetValue extracts the ListElement value from a skiplist node.
func (e *SkipListElement) GetValue() ListElement {
	return e.value
}
```

**翻译及分析：**

- **GetValue 方法：** 从跳表节点中提取 `ListElement` 的值。
  - **接收者：** 指向 `SkipListElement` 的指针 `e`。
  - **返回值：** 节点存储的 `ListElement` 实例。
  - **步骤：**
    1. 返回节点的 `value` 字段。

### GetSmallestNode

```go
// GetSmallestNode returns the very first/smallest node in the skiplist.
// GetSmallestNode runs in O(1)
func (t *SkipList) GetSmallestNode() *SkipListElement {
	return t.startLevels[0]
}
```

**翻译及分析：**

- **GetSmallestNode 方法：** 返回跳表中最小的节点。
  - **接收者：** 指向 `SkipList` 的指针 `t`。
  - **返回值：** 指向第 0 层起始节点的指针。
  - **步骤：**
    1. 返回跳表第 0 层的起始节点 `t.startLevels[0]`。

**时间复杂度：** \( O(1) \)。

### GetLargestNode

```go
// GetLargestNode returns the very last/largest node in the skiplist.
// GetLargestNode runs in O(1)
func (t *SkipList) GetLargestNode() *SkipListElement {
	return t.endLevels[0]
}
```

**翻译及分析：**

- **GetLargestNode 方法：** 返回跳表中最大的节点。
  - **接收者：** 指向 `SkipList` 的指针 `t`。
  - **返回值：** 指向第 0 层末尾节点的指针。
  - **步骤：**
    1. 返回跳表第 0 层的末尾节点 `t.endLevels[0]`。

**时间复杂度：** \( O(1) \)。

### Next

```go
// Next returns the next element based on the given node.
// Next will loop around to the first node, if you call it on the last!
func (t *SkipList) Next(e *SkipListElement) *SkipListElement {
	if e.next[0] == nil {
		return t.startLevels[0]
	}
	return e.next[0]
}
```

**翻译及分析：**

- **Next 方法：** 返回给定节点的下一个节点。如果调用在最后一个节点上，则循环回第一个节点。
  - **参数：** `e *SkipListElement`：当前节点。
  - **返回值：** 下一个节点的指针。
  - **步骤：**
    1. 如果当前节点在第 0 层的下一个节点为 `nil`，则返回第 0 层的起始节点（循环回第一个节点）。
    2. 否则，返回当前节点在第 0 层的下一个节点。

### Prev

```go
// Prev returns the previous element based on the given node.
// Prev will loop around to the last node, if you call it on the first!
func (t *SkipList) Prev(e *SkipListElement) *SkipListElement {
	if e.prev == nil {
		return t.endLevels[0]
	}
	return e.prev
}
```

**翻译及分析：**

- **Prev 方法：** 返回给定节点的上一个节点。如果调用在第一个节点上，则循环回最后一个节点。
  - **参数：** `e *SkipListElement`：当前节点。
  - **返回值：** 上一个节点的指针。
  - **步骤：**
    1. 如果当前节点的 `prev` 指针为 `nil`，则返回第 0 层的末尾节点（循环回最后一个节点）。
    2. 否则，返回当前节点的 `prev` 指针指向的节点。

### GetNodeCount

```go
// GetNodeCount returns the number of nodes currently in the skiplist.
func (t *SkipList) GetNodeCount() int {
	return t.elementCount
}
```

**翻译及分析：**

- **GetNodeCount 方法：** 返回当前跳表中的节点数。
  - **接收者：** 指向 `SkipList` 的指针 `t`。
  - **返回值：** 跳表中元素的计数 `t.elementCount`。

### ChangeValue

```go
// ChangeValue can be used to change the actual value of a node in the skiplist
// without the need of Deleting and reinserting the node again.
// Be advised, that ChangeValue only works, if the actual key from ExtractKey() will stay the same!
// ok is an indicator, wether the value is actually changed.
func (t *SkipList) ChangeValue(e *SkipListElement, newValue ListElement) (ok bool) {
	// The key needs to stay correct, so this is very important!
	if math.Abs(newValue.ExtractKey()-e.key) <= t.eps {
		e.value = newValue
		ok = true
	} else {
		ok = false
	}
	return
}
```

**翻译及分析：**

- **ChangeValue 方法：** 用于在不需要删除和重新插入节点的情况下更改跳表中节点的实际值。
  - **参数：**
    - `e *SkipListElement`：要更改的节点。
    - `newValue ListElement`：新的值，实现了 `ListElement` 接口。
  - **返回值：**
    - `ok bool`：指示值是否成功更改。
  - **步骤：**
    1. 检查新值的键值 `newValue.ExtractKey()` 与当前节点的键值 `e.key` 之差的绝对值是否小于等于 `eps`。
       - 如果是，则更改节点的 `value` 为 `newValue`，并返回 `ok = true`。
       - 否则，不更改，返回 `ok = false`。

**注意：** 这个方法仅在新值的键值与当前节点的键值保持一致时有效，以确保跳表的有序性不被破坏。

### String

```go
// String returns a string format of the skiplist. Useful to get a graphical overview and/or debugging.
func (t *SkipList) String() string {
	s := ""

	s += " --> "
	for i, l := range t.startLevels {
		if l == nil {
			break
		}
		if i > 0 {
			s += " -> "
		}
		next := "---"
		if l != nil {
			next = l.value.String()
		}
		s += fmt.Sprintf("[%v]", next)

		if i == 0 {
			s += "    "
		}
	}
	s += "\n"

	node := t.startLevels[0]
	for node != nil {
		s += fmt.Sprintf("%v: ", node.value)
		for i := 0; i <= node.level; i++ {

			l := node.next[i]

			next := "---"
			if l != nil {
				next = l.value.String()
			}

			if i == 0 {
				prev := "---"
				if node.prev != nil {
					prev = node.prev.value.String()
				}
				s += fmt.Sprintf("[%v|%v]", prev, next)
			} else {
				s += fmt.Sprintf("[%v]", next)
			}
			if i < node.level {
				s += " -> "
			}

		}
		s += "\n"
		node = node.next[0]
	}

	s += " --> "
	for i, l := range t.endLevels {
		if l == nil {
			break
		}
		if i > 0 {
			s += " -> "
		}
		next := "---"
		if l != nil {
			next = l.value.String()
		}
		s += fmt.Sprintf("[%v]", next)
		if i == 0 {
			s += "    "
		}
	}
	s += "\n"
	return s
}
```

**翻译及分析：**

- **String 方法：** 返回跳表的字符串表示，便于图形化概览和调试。
  - **接收者：** 指向 `SkipList` 的指针 `t`。
  - **返回值：** 跳表的字符串表示。
  - **步骤：**
    1. 初始化字符串 `s`。
    2. 构建顶部指针行：
       - 遍历所有层数的起始节点 `t.startLevels`。
       - 对于每一层，添加其起始节点的字符串表示，格式为 `[节点值]`；如果没有节点，则显示 `---`。
       - 层数之间使用 `" -> "` 分隔。
    3. 换行。
    4. 遍历第 0 层的所有节点，逐行构建每个节点的详细连接信息：
       - 对于每个节点，添加其值 `node.value` 的字符串表示。
       - 遍历该节点的所有层数，添加其在每层的前后连接情况：
         - 对于第 0 层，显示 `[prev|next]`，表示前一个和下一个节点的值；如果没有前一个或下一个节点，则显示 `---`。
         - 对于其他层，显示 `[next]`，表示该层的下一个节点的值；如果没有下一个节点，则显示 `---`。
         - 层数之间使用 `" -> "` 分隔。
       - 换行。
    5. 构建底部指针行，与顶部指针行类似，遍历所有层数的末尾节点 `t.endLevels`。
    6. 换行并返回最终的字符串 `s`。

**示例输出格式：**

```
 --> [Value1] -> [Value2] -> [Value3]
Value1: [---|Value2] -> [---]
Value2: [Value1|Value3] -> [Value1|---]
Value3: [Value2|---] -> [Value2|---]
 --> [Value1] -> [Value3]
```

**补充说明：**

此方法主要用于调试和可视化跳表的内部结构，包括各层的连接情况和节点的前后关系。

---

## 总结

该跳表实现提供了基本的插入、删除、查找和遍历操作，保证了在 \( O(\log n) \) 的时间复杂度内完成这些操作。通过多层级的索引机制，跳表能够高效地管理大量有序元素，并支持快速的遍历和更新。

**关键点总结：**

- **多层索引：** 跳表通过多层索引提高查找效率，每一层的节点数量大约为下一层的一半，从而使整个结构类似于多层平衡树。
- **随机层数生成：** 新节点的层数通过随机数生成，使得跳表的层数和节点分布具有随机性，避免最坏情况的出现。
- **双向遍历：** 节点在第 0 层维护了双向指针（`prev` 和 `next`），支持双向遍历操作。
- **空间效率：** 尽管跳表需要额外的空间来存储多层指针，但其空间复杂度仍然是线性的，相较于平衡树具有更好的空间效率。

通过以上分析，可以看出该跳表实现结构清晰，功能完备，适用于需要高效有序数据管理的场景。
