package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type Entry struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Tombstone bool   `json:"tombstone"` // 标记是否为删除操作
	Version   int64  `json:"version"`   // 可选：表示数据的版本或时间戳
}

// 内存中的有序数据结构，负责接收所有的写入操作。
// 考虑到并发度，一般使用跳表。
type MemTable interface {
	Insert(entry Entry) error     // 插入或更新键值对
	Get(key string) (Entry, bool) // 根据键获取对应的 Entry
	Delete(key string) error      // 标记键为删除（插入 Tombstone）
	Entries() []Entry             // 获取所有条目（用于刷写）
	Size() int                    // 当前 MemTable 的大小
	Clear()                       // 清空 MemTable
}

// Write-Ahead Log (WAL)
// 确保数据的持久性、故障恢复，所有写入操作在 MemTable 更新前首先记录到 WAL 中。
type WAL interface {
	Append(entries []Entry) error // 追加多个 Entry 到 WAL
	ReadAll() ([]Entry, error)    // 读取所有 Entry（用于恢复）
	Truncate() error              // 清空 WAL
	Close() error                 // 关闭 WAL 文件
}

// 磁盘上的不可变有序键值对集合，支持高效的查找和范围查询。
type SSTable interface {
	Load(path string) error                     // 从文件加载 SSTable
	Search(key string) (Entry, bool)            // 根据键搜索 Entry
	RangeQuery(startKey, endKey string) []Entry // 范围查询
	Write(path string, entries []Entry) error   // 将 Entry 写入 SSTable 文件
	MinKey() string                             // 获取 SSTable 的最小键
	MaxKey() string                             // 获取 SSTable 的最大键
}

// 用于快速判断某个键是否不存在于 SSTable 中，减少不必要的磁盘查找。
type BloomFilter interface {
	Add(key string)         // 添加键到过滤器
	Test(key string) bool   // 测试键是否可能存在
	Save(path string) error // 保存 Bloom Filter 到文件
	Load(path string) error // 从文件加载 Bloom Filter
}

// 负责执行合并和压缩操作，将多个 SSTable 合并为一个，去除重复和删除的数据。
type Compactor interface {
	Compact(sstables []SSTable) (SSTable, error) // 合并多个 SSTable，返回新的 SSTable
}

// LSM 树的核心接口，整合上述组件，并定义主要的操作方法。
type LSMTree interface {
	Insert(key, value string) error // 插入或更新键值对
	Get(key string) (string, error) // 根据键获取值
	Delete(key string) error        // 删除键（插入 Tombstone）
	Flush() error                   // 刷写 MemTable 到 SSTable
	Start() error                   // 启动 LSM 树（如启动压缩服务）
	Close() error                   // 关闭 LSM 树，释放资源
}

// - 插入操作（Insert）
//    1. 记录写入：将 Entry 记录到 Write-Ahead Log（WAL）中，确保数据持久性。
//    2. 更新 MemTable：将 Entry 插入到 MemTable 中，实现快速的内存写入。
//    3. 检查 MemTable 大小：如果 MemTable 达到预设的大小阈值，触发刷写操作（Flush）。
// - 读取操作（Get）
// 		1. 查询 MemTable：首先在 MemTable 中查找键值对，如果找到且非 Tombstone，返回对应值。
// 		2. 遍历 SSTable：
// 				按照 SSTable 的优先级（通常是从最新到最旧）进行遍历。
// 				利用 Bloom Filter 判断该键是否可能存在于当前 SSTable。
// 				如果 Bloom Filter 可能存在，则在 SSTable 中执行实际的键查找。
// 				一旦找到最新的非 Tombstone Entry，返回对应值并停止查找。
// - 删除操作（Delete）
// 		1. 插入 Tombstone：创建一个带有 Tombstone 标记的 Entry，表示该键已被删除。
// 		2. 记录删除：将 Tombstone 记录追加到 WAL 中。
// 		3. 更新 MemTable：将 Tombstone 插入到 MemTable 中。
// 		4. 检查 MemTable 大小：若 MemTable 满，触发刷写操作。
// - 压缩操作（Compaction）
// 		1. 选择待合并的 SSTables：根据合并策略（如 Leveled 或 Size-Tiered），选择需要合并的 SSTable 列表。
// 		2. 加载和排序数据：将待合并的 SSTable 中的所有 Entry 加载到内存，并按键排序。
// 		3. 去重与处理 Tombstones：保留最新版本的 Entry，移除已删除的键（根据 Tombstone）。
// 		4. 写入新 SSTable：将整理后的数据写入新的 SSTable 文件。
// 		5. 更新 SSTable 列表：将新 SSTable 添加到 SSTable 层级中，并移除已合并的旧 SSTable。
// - 关闭与资源清理
// 		1. 停止后台服务：如压缩服务、缓存刷新等。
//		2. 关闭 WAL 和 SSTable 文件：确保所有文件正确关闭，数据完整性得以保证。

// ====== 实现 LSMTree ======

var NewSimpleMemTable func() MemTable
var NewSimpleWAL func(path string) WAL
var NewSimpleCompactor func() Compactor
var NewSimpleSSTable func(entries []Entry, filePath string) SSTable

type SimpleLSMTree struct {
	mu         sync.RWMutex
	memTable   MemTable
	wal        WAL
	sstables   []SSTable
	compactor  Compactor
	dir        string
	maxMemSize int
}

func NewSimpleLSMTree(dir string, maxMemSize int) *SimpleLSMTree {
	memTable := NewSimpleMemTable()
	wal := NewSimpleWAL(dir + "/wal.log")
	compactor := NewSimpleCompactor()
	return &SimpleLSMTree{
		memTable:   memTable,
		wal:        wal,
		sstables:   []SSTable{},
		compactor:  compactor,
		dir:        dir,
		maxMemSize: maxMemSize,
	}
}

func (lsm *SimpleLSMTree) Insert(key, value string) error {
	entry := Entry{
		Key:       key,
		Value:     value,
		Tombstone: false,
	}

	// Append 到 WAL
	err := lsm.wal.Append([]Entry{entry})
	if err != nil {
		return err
	}

	// Insert 到 MemTable
	err = lsm.memTable.Insert(entry)
	if err != nil {
		return err
	}

	// 检查 MemTable 大小
	if lsm.memTable.Size() >= lsm.maxMemSize {
		return lsm.Flush()
	}

	return nil
}

func (lsm *SimpleLSMTree) Get(key string) (string, error) {
	// 查询 MemTable
	if entry, exists := lsm.memTable.Get(key); exists {
		if entry.Tombstone {
			return "", errors.New("key not found (deleted)")
		}
		return entry.Value, nil
	}

	// 遍历 SSTables，从最新到最旧
	for i := len(lsm.sstables) - 1; i >= 0; i-- {
		sst := lsm.sstables[i]
		entry, exists := sst.Search(key)
		if exists {
			if entry.Tombstone {
				return "", errors.New("key not found (deleted)")
			}
			return entry.Value, nil
		}
	}

	return "", errors.New("key not found")
}

func (lsm *SimpleLSMTree) Delete(key string) error {
	// 创建 Tombstone Entry
	entry := Entry{
		Key:       key,
		Value:     "",
		Tombstone: true,
	}

	// Append 到 WAL
	err := lsm.wal.Append([]Entry{entry})
	if err != nil {
		return err
	}

	// Insert Tombstone 到 MemTable
	err = lsm.memTable.Delete(key)
	if err != nil {
		return err
	}

	// 检查 MemTable 大小
	if lsm.memTable.Size() >= lsm.maxMemSize {
		return lsm.Flush()
	}

	return nil
}

func (lsm *SimpleLSMTree) Flush() error {
	// 获取 MemTable 的所有条目
	entries := lsm.memTable.Entries()

	// 清空 MemTable
	lsm.memTable.Clear()

	// 写入新的 SSTable
	newSST := NewSimpleSSTable(entries, fmt.Sprintf("%s/sstable-%d.sst", lsm.dir, len(lsm.sstables)+1))
	err := newSST.Write(lsm.dir, entries)
	if err != nil {
		return err
	}

	// 添加到 SSTable 列表
	lsm.sstables = append(lsm.sstables, newSST)

	return nil
}

func (lsm *SimpleLSMTree) Compact() error {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	// 选择需要合并的 SSTables，示例中选择前两个
	if len(lsm.sstables) < 2 {
		return nil // 不足两个进行合并
	}

	toCompact := lsm.sstables[:2]
	newSST, err := lsm.compactor.Compact(toCompact)
	if err != nil {
		return err
	}

	// 添加新 SSTable
	lsm.sstables = append(lsm.sstables, newSST)

	// 移除被合并的 SSTable
	lsm.sstables = lsm.sstables[2:]

	return nil
}

func (lsm *SimpleLSMTree) Start() error {
	// 启动后台压缩服务
	go func() {
		for {
			// 定期触发压缩操作
			err := lsm.Compact()
			if err != nil {
				log.Printf("Compaction error: %v", err)
			}
			// 休眠一段时间后再次触发
			// 示例中使用固定休眠时间，实际中可根据负载动态调整
			fmt.Println("Sleeping before next compaction...")
			// time.Sleep(time.Hour)
			break // 仅示例一次
		}
	}()
	return nil
}

func (lsm *SimpleLSMTree) Close() error {
	// 刷新 MemTable
	err := lsm.Flush()
	if err != nil {
		return err
	}

	// 关闭 WAL
	err = lsm.wal.Close()
	if err != nil {
		return err
	}

	return nil
}

// ====== 使用示例 ======

func main() {
	// 创建 LSMTree 实例
	lsm := NewSimpleLSMTree("data", 3)
	defer lsm.Close()

	// 启动 LSM 树（启动压缩服务）
	err := lsm.Start()
	if err != nil {
		log.Fatalf("Failed to start LSM Tree: %v", err)
	}

	// 插入数据
	err = lsm.Insert("key1", "value1")
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}
	err = lsm.Insert("key2", "value2")
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}
	err = lsm.Insert("key3", "value3")
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}

	// 触发一次 Flush（自动触发）
	err = lsm.Insert("key4", "value4")
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}

	// 读取数据
	val, err := lsm.Get("key2")
	if err != nil {
		fmt.Printf("Get key2 error: %v\n", err)
	} else {
		fmt.Printf("Get key2: %s\n", val)
	}

	// 删除数据
	err = lsm.Delete("key2")
	if err != nil {
		log.Fatalf("Delete error: %v", err)
	}

	// 读取已删除的数据
	val, err = lsm.Get("key2")
	if err != nil {
		fmt.Printf("Get key2 after deletion error: %v\n", err)
	} else {
		fmt.Printf("Get key2 after deletion: %s\n", val)
	}

	// 模拟压缩
	err = lsm.Compact()
	if err != nil {
		log.Fatalf("Compaction error: %v", err)
	}

	// 再次读取数据
	val, err = lsm.Get("key1")
	if err != nil {
		fmt.Printf("Get key1 error: %v\n", err)
	} else {
		fmt.Printf("Get key1: %s\n", val)
	}
}
