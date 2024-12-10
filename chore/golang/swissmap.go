// 使用 Go 语言实现 Swiss Table 哈希表的一个简单示例。
// Swiss Table 的核心思想：

// - 控制字节（Control Byte）：使用哈希值的低 8 位作为控制字节，用于加速查找操作。
// - 桶组（Bucket Group）：将桶（bucket）划分为组，每组包含 16 个桶，以便利用 SIMD 指令一次处理多个桶。
// - 元数据和数据分离：将控制字节（元数据）和实际的键值对（数据）分离存储，以提高缓存命中率。

package main

import (
	"errors"
	"fmt"
	"sync"
)

func main() {
	table := NewTable(64)

	table.Insert("apple", 1)
	table.Insert("banana", 2)
	table.Insert("cherry", 3)

	if val, ok := table.Find("banana"); ok {
		fmt.Println("banana:", val)
	}

	table.Delete("apple")
	if _, ok := table.Find("apple"); !ok {
		fmt.Println("apple not found")
	}
}

const (
	GroupSize = 16 // 每个桶组包含的桶数量
)

type controlByte uint8

const (
	Empty    controlByte = 0b10000000
	Deleted  controlByte = 0b11111110
	Sentinel controlByte = 0b11111111
)

// Entry 存储键值对
type Entry struct {
	key   interface{}
	value interface{}
}

// BucketGroup 存储控制字节和对应的 Entry
type BucketGroup struct {
	controls [GroupSize]controlByte
	entries  [GroupSize]*Entry
}

// Table 哈希表结构
type Table struct {
	buckets []*BucketGroup
	size    int
	lock    sync.RWMutex
}

// NewTable 创建新的哈希表
func NewTable(capacity int) *Table {
	numGroups := (capacity + GroupSize - 1) / GroupSize
	buckets := make([]*BucketGroup, numGroups)
	for i := 0; i < numGroups; i++ {
		buckets[i] = &BucketGroup{}
		// 初始化控制字节为 Empty
		for j := 0; j < GroupSize; j++ {
			buckets[i].controls[j] = Empty
		}
	}
	return &Table{
		buckets: buckets,
	}
}

// 哈希函数（简单示例，实际应用中应使用更好的哈希函数）
func hash(key interface{}) uint64 {
	return uint64((uintptr)(key.(string)[0])) // 假设 key 是字符串
}

// Insert 插入键值对
func (t *Table) Insert(key, value interface{}) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	h := hash(key)
	ctrl := controlByte(h & 0x7F) // 取哈希值低 7 位作为控制字节
	index := int(h % uint64(len(t.buckets)))
	group := t.buckets[index]

	for i := 0; i < GroupSize; i++ {
		if group.controls[i] == Empty || group.controls[i] == Deleted {
			group.controls[i] = ctrl
			group.entries[i] = &Entry{key, value}
			t.size++
			return nil
		}
	}

	// 处理冲突（简化处理，实际应重新哈希或扩容）
	return errors.New("hash table is full")
}

// Find 查找键对应的值
func (t *Table) Find(key interface{}) (interface{}, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	h := hash(key)
	ctrl := controlByte(h & 0x7F)
	index := int(h % uint64(len(t.buckets)))
	group := t.buckets[index]

	for i := 0; i < GroupSize; i++ {
		if group.controls[i] == ctrl {
			if group.entries[i].key == key {
				return group.entries[i].value, true
			}
		} else if group.controls[i] == Empty {
			return nil, false
		}
	}
	return nil, false
}

// Delete 删除指定键
func (t *Table) Delete(key interface{}) bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	h := hash(key)
	ctrl := controlByte(h & 0x7F)
	index := int(h % uint64(len(t.buckets)))
	group := t.buckets[index]

	for i := 0; i < GroupSize; i++ {
		if group.controls[i] == ctrl && group.entries[i].key == key {
			group.controls[i] = Deleted
			group.entries[i] = nil
			t.size--
			return true
		} else if group.controls[i] == Empty {
			return false
		}
	}
	return false
}
