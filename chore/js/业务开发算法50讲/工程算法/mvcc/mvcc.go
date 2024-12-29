// 展示了一个简化的MVCC实现，包括事务管理、数据版本控制、可见性规则和垃圾回收。
// 实际的数据库系统中，MVCC的实现会更加复杂，涉及更多的优化和细节处理，例如并发控制、索引支持、锁机制等。

package main

import (
	"fmt"
	"sync"
	"time"
)

// Version 表示数据的一个版本
type Version struct {
	Data string // 数据内容

	// 元数据(隐藏列)
	Xmin      int64     // 创建该版本的事务ID
	Xmax      int64     // 软删除：删除该版本的事务ID，0表示未删除
	Timestamp time.Time // 版本创建时间
}

// Record 表示数据库中的一行记录，包含多个版本
type Record struct {
	ID       int64
	Versions []Version // !每条数据记录维护一个版本链
	mu       sync.RWMutex
}

// Transaction 表示一个事务
type Transaction struct {
	ID        int64
	StartTime time.Time
	// 可以扩展更多字段，如事务状态等
}

// Database 表示一个简单的内存数据库
type Database struct {
	records      map[int64]*Record
	transactions map[int64]*Transaction
	mu           sync.RWMutex
	nextTxID     int64
	nextRecID    int64
}

// NewDatabase 创建一个新的数据库实例
func NewDatabase() *Database {
	return &Database{
		records:      make(map[int64]*Record),
		transactions: make(map[int64]*Transaction),
		nextTxID:     1,
		nextRecID:    1,
	}
}

// BeginTransaction 开始一个新的事务
// TODO：每个事务在开始时会获取一个数据的快照，这个快照反映了该事务开始时的数据库状态
func (db *Database) BeginTransaction() *Transaction {
	db.mu.Lock()
	defer db.mu.Unlock()
	tx := &Transaction{
		ID:        db.nextTxID,
		StartTime: time.Now(),
	}
	db.transactions[db.nextTxID] = tx
	db.nextTxID++
	fmt.Printf("Transaction %d started.\n", tx.ID)
	return tx
}

// CommitTransaction 提交一个事务
func (db *Database) CommitTransaction(tx *Transaction) {
	db.mu.Lock()
	defer db.mu.Unlock()
	// 在这里可以添加冲突检测等逻辑
	delete(db.transactions, tx.ID)
	fmt.Printf("Transaction %d committed.\n", tx.ID)
}

// Insert 插入一条新记录
func (db *Database) Insert(tx *Transaction, data string) int64 {
	db.mu.Lock()
	defer db.mu.Unlock()
	recID := db.nextRecID
	db.nextRecID++

	version := Version{
		Data:      data,
		Xmin:      tx.ID,
		Xmax:      0,
		Timestamp: time.Now(),
	}

	record := &Record{
		ID:       recID,
		Versions: []Version{version},
	}
	db.records[recID] = record
	fmt.Printf("Transaction %d inserted record %d with data '%s'.\n", tx.ID, recID, data)
	return recID
}

// Read 读取一条记录的数据，根据事务的可见性规则，未被删除且创建时间早于事务开始时间
func (db *Database) Read(tx *Transaction, recID int64) (string, bool) {
	db.mu.RLock()
	record, exists := db.records[recID]
	db.mu.RUnlock()
	if !exists {
		return "", false
	}

	record.mu.RLock()
	defer record.mu.RUnlock()

	var visibleVersion *Version
	for i := len(record.Versions) - 1; i >= 0; i-- {
		v := &record.Versions[i]
		// !未被删除且创建时间早于事务开始时间
		if v.Xmin <= tx.ID && (v.Xmax == 0 || v.Xmax > tx.ID) {
			visibleVersion = v
			break
		}
	}

	if visibleVersion != nil {
		fmt.Printf("Transaction %d read record %d with data '%s'.\n", tx.ID, recID, visibleVersion.Data)
		return visibleVersion.Data, true
	}

	return "", false
}

// Update 更新一条记录的数据，创建新版本
func (db *Database) Update(tx *Transaction, recID int64, newData string) bool {
	db.mu.RLock()
	record, exists := db.records[recID]
	db.mu.RUnlock()
	if !exists {
		return false
	}

	record.mu.Lock()
	defer record.mu.Unlock()

	// 找到当前可见的最新版本
	var currentVersion *Version
	for i := len(record.Versions) - 1; i >= 0; i-- {
		v := &record.Versions[i]
		if v.Xmin <= tx.ID && (v.Xmax == 0 || v.Xmax > tx.ID) {
			currentVersion = v
			break
		}
	}

	if currentVersion == nil {
		return false
	}

	// 设置当前版本的Xmax为当前事务ID，表示版本失效
	currentVersion.Xmax = tx.ID

	// 创建新版本
	newVersion := Version{
		Data:      newData,
		Xmin:      tx.ID,
		Xmax:      0,
		Timestamp: time.Now(),
	}
	record.Versions = append(record.Versions, newVersion)
	fmt.Printf("Transaction %d updated record %d with new data '%s'.\n", tx.ID, recID, newData)
	return true
}

// GarbageCollect 清理不再被任何事务引用的旧版本数据
func (db *Database) GarbageCollect() {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, record := range db.records {
		record.mu.Lock()
		newVersions := []Version{}
		for _, v := range record.Versions {
			// !如果没有事务的Xmin小于该版本的Xmin，则可以回收
			canDelete := true
			for _, tx := range db.transactions {
				if tx.ID < v.Xmin {
					canDelete = false
					break
				}
			}
			if !canDelete {
				newVersions = append(newVersions, v)
			}
		}
		record.Versions = newVersions
		record.mu.Unlock()
	}
	fmt.Println("Garbage collection completed.")
}

func main() {
	db := NewDatabase()

	// 事务1开始
	tx1 := db.BeginTransaction()
	// 事务1插入一条记录
	recID := db.Insert(tx1, "Initial Data")
	// 事务1提交
	db.CommitTransaction(tx1)

	// 事务2开始
	tx2 := db.BeginTransaction()
	// 事务2读取记录
	data, exists := db.Read(tx2, recID)
	if exists {
		fmt.Printf("Transaction %d read data: %s\n", tx2.ID, data)
	}

	// 事务3开始
	tx3 := db.BeginTransaction()
	// 事务3更新记录
	db.Update(tx3, recID, "Updated Data")
	// 事务3提交
	db.CommitTransaction(tx3)

	// 事务2再次读取记录，应该读取旧版本
	data, exists = db.Read(tx2, recID)
	if exists {
		fmt.Printf("Transaction %d read data after update: %s\n", tx2.ID, data)
	}

	// 事务4开始
	tx4 := db.BeginTransaction()
	// 事务4读取记录，应该读取新版本
	data, exists = db.Read(tx4, recID)
	if exists {
		fmt.Printf("Transaction %d read data: %s\n", tx4.ID, data)
	}

	// 执行垃圾回收
	db.GarbageCollect()
}
