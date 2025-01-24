// https://mp.weixin.qq.com/s?__biz=MzkxMjQzMjA0OQ==&mid=2247484228&idx=1&sn=f30aa1ec0ae19d934beaa3ab13e2c3ee&chksm=c10c4d9af67bc48ce5f5c414b2a8731f5822d0e3b8a8c0eb60ef5a022795d6fdaded70ce0a5a&cur_album_id=2709593649634033668&scene=189#wechat_redirect
// 如何实现一个并发安全的跳表
// 由于 redis 采用单线程处理模型，因此不存在并发访问跳表的诉求；
// rocksdb 则采用多线程处理模型，并发读写内存数据结构时，需要兼顾数据一致以及操作性能，此时就需要使用到并发安全的跳表结构.

package main

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

func main() {

}

// 基于节点粒度锁实现的并发安全的跳表结构
type ConcurrentSkipList struct {
	// 当前跳表中存在的元素个数，通过 atomic.Int32 保证增减操作的原子性
	cap atomic.Int32

	// 只有删除操作取写锁，其他操作均取读锁
	// 通过该锁实现了删除操作的单独互斥处理
	DeleteMutex sync.RWMutex

	// 同一个 key 的 put 操作需要串行化，会基于 sync.Map 进行并发安全地存储管理
	keyToMutex sync.Map

	// 跳表的头节点
	head *node

	// 对象池，复用跳表中创建和删除的 node 结构，减轻 gc 压力
	nodesCache sync.Pool

	// 比较 node key 大小的规则，倘若 key1 < key2 返回 true，否则返回 false
	compareFunc func(key1, key2 any) bool
}

type node struct {
	// key，val 对
	key, val any
	// nexts[i] 为第 i 层的 next 节点
	nexts []*node
	// 每个节点持有一把节点粒度的读写锁，后续可以作为左边界锁
	sync.RWMutex
}

// 构造一个并发安全的跳表，需要注入比较 key 大小的规则函数
func NewConcurrentSkipList(compareFunc func(key1, key2 any) bool) *ConcurrentSkipList {
	return &ConcurrentSkipList{
		// 跳表的初始高度为 1
		head: &node{
			nexts: make([]*node, 1),
		},
		// 初始化 node 对象池
		nodesCache: sync.Pool{
			New: func() any {
				return &node{}
			},
		},
		// 注入 key 比较函数
		compareFunc: compareFunc,
	}
}

// 根据 key 删除跳表中对应的 key-value
func (c *ConcurrentSkipList) Del(key any) {
	// 针对于 deleteMutex 加写锁，保证删除操作的全局互斥性
	// 由于此处保证了删除操作的独享性，因此后续操作都无需加锁
	c.DeleteMutex.Lock()
	defer c.DeleteMutex.Unlock()

	// 用于接收 key 对应的节点，后续将其投放回到对象池中
	var deleteNode *node

	// 从头节点最高层开始遍历
	move := c.head
	for level := len(c.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && c.compareFunc(move.nexts[level].key, key) {
			move = move.nexts[level]
		}

		// 此时已经来到对应层的左边界
		// 倘若当前层不存在 key 对应的 node，则直接跳过进入下一层
		if move.nexts[level] == nil || (move.nexts[level].key != key && !c.compareFunc(move.nexts[level].key, key)) {
			continue
		}

		// 在当前层找到存在 key 对应的 node，将其赋给 deleteNode
		if deleteNode == nil {
			deleteNode = move.nexts[level]
		}

		// 调整左边界指针引用实现删除 node 的效果
		move.nexts[level] = move.nexts[level].nexts[level]
	}

	// 倘若 key 不存在，直接返回即可
	if deleteNode == nil {
		return
	}

	// 走到此处意味着 key 存在
	// 方法返回前将跳表的节点个数器 - 1
	defer c.cap.Add(-1)
	// 提前将 deleteNode 的 nexts 置为空，节省空间
	deleteNode.nexts = nil
	// 将 deleteNode 放回对象池
	c.nodesCache.Put(deleteNode)

	// 尝试进行跳表高度缩容
	var dif int
	for level := len(c.head.nexts) - 1; level > 0; level-- {
		// 如果该层头节点next为空
		dif++
	}
	c.head.nexts = c.head.nexts[:len(c.head.nexts)-dif]
}

// 根据 key 获取对应的 value
func (c *ConcurrentSkipList) Get(key any) (any, bool) {
	// 取 deleteMutex 的读锁. 和 delete 操作实现互斥，但是 get 操作和 put 操作本身可以共享，基于更细粒度的节点锁实现互斥
	c.DeleteMutex.RLock()
	defer c.DeleteMutex.RUnlock()

	// 沿着头节点 head 从最高层出发进行检索
	move := c.head
	// 通过 last 记录上一层所在的节点位置，避免在下降的过程中对于同一个节点反复加多次左边界节点锁
	var last *node
	for level := len(c.head.nexts) - 1; level >= 0; level-- {
		// 在同一层中一路无锁穿越，直到来到左边界
		for move.nexts[level] != nil && c.compareFunc(move.nexts[level].key, key) {
			move = move.nexts[level]
		}

		// 走到左边界
		// 通过 last 指针保证对同一个节点只会加一次左边界节点锁
		if move != last {
			// get 操作针对左边界节点加读锁，保证多个 get 操作可以共享
			move.RLock()
			defer move.RUnlock()
			// 更新 last 指针引用
			last = move
		}

		// 倘若找到目标，则直接返回
		if move.nexts[level] != nil && move.nexts[level].key == key {
			return move.nexts[level].val, true
		}
	}

	// 遍历完也没找到目标，说明 key 不存在
	return 0, false
}

// 写入 key-value 对
func (c *ConcurrentSkipList) Put(key, val any) {
	// 取 deleteMutex 的读锁. 和 delete 操作实现互斥，但是 get 操作和 put 操作本身可以共享，基于更细粒度的节点锁实现互斥
	c.DeleteMutex.RLock()
	defer c.DeleteMutex.RUnlock()

	// 针对于同一个 key 的 put 操作需要互斥
	keyMutex := c.getKeyMutex(key)
	keyMutex.Lock()
	defer keyMutex.Unlock()

	// 通过 search 操作获取到 key 对应的 node，倘若 key 存在，则直接更新 node 的值然后返回即可
	if _node := c.search(key); _node != nil {
		_node.val = val
		return
	}

	// 走到此处，说明针对 key 的操作必然是一次插入操作，因此在方法返回前对节点个数计数器累加 1
	defer c.cap.Add(1)

	// 随机出新节点的高度
	rLevel := c.randomLevel()

	// 通过对象池创建出新的节点
	newNode, _ := c.nodesCache.Get().(*node)
	newNode.key, newNode.val = key, val
	newNode.nexts = make([]*node, rLevel+1)

	// 对创建出的新节点需要加写锁，避免在插入过程中成为 get 流程的左边界，产生脏数据
	newNode.Lock()
	defer newNode.Unlock()

	// 如果新节点导致跳表需要扩容，则需要对头节点加锁
	if rLevel > len(c.head.nexts)-1 {
		c.head.Lock()
		for rLevel > len(c.head.nexts)-1 {
			c.head.nexts = append(c.head.nexts, nil)
		}
		c.head.Unlock()
	}

	// 沿着头节点 head 从最高层出发进行检索
	move := c.head
	// 通过 last 记录上一层所在的节点位置，避免在下降的过程中对于同一个节点反复加多次左边界节点
	var last *node
	for level := rLevel; level >= 0; level-- {
		// 在同一层中一路无锁穿越，直到来到左边界
		for move.nexts[level] != nil && c.compareFunc(move.nexts[level].key, key) {
			move = move.nexts[level]
		}

		// 走到左边界
		// 通过 last 指针保证对同一个节点只会加一次左边界节点锁
		if move != last {
			// 插入操作需要对左边界加写锁，实现独享
			move.Lock()
			defer move.Unlock()
			last = move
		}

		// 插入新节点需要调整指针引用
		// 需要先设置 newNode 的指针，再设置 newNode 左边界的指针
		newNode.nexts[level] = move.nexts[level]
		move.nexts[level] = newNode
	}
}

// 根据 key 获取到对应的 key 锁
func (c *ConcurrentSkipList) getKeyMutex(key any) *sync.Mutex {
	// 基于 sync.Map 实现 key 锁的管理
	// 通过 LoadOrStore 方法获取 key 锁，倘若此前即存在则直接获取，否则实现存储操作，并获取到此时放入的锁
	rawMutex, _ := c.keyToMutex.LoadOrStore(key, &sync.Mutex{})
	mutex, _ := rawMutex.(*sync.Mutex)
	return mutex
}

// 随机出新节点的最大层数索引
func (c *ConcurrentSkipList) randomLevel() int {
	var level int
	for rand.Intn(2) > 0 {
		level++
	}
	return level
}

// 通过 key 检索到对应的 node
func (c *ConcurrentSkipList) search(key any) *node {

	// 沿着头节点 head 从最高层出发进行检索
	move := c.head
	// 通过 last 记录上一层所在的节点位置，避免在下降的过程中对于同一个节点反复加多次左边界节点
	var last *node
	for level := len(c.head.nexts) - 1; level >= 0; level-- {
		// 在同一层中一路无锁穿越，直到来到左边界
		for move.nexts[level] != nil && c.compareFunc(move.nexts[level].key, key) {
			move = move.nexts[level]
		}

		// 走到左边界
		// 通过 last 指针保证对同一个节点只会加一次左边界节
		if move != last {
			// 对左边界节点加读锁
			move.RLock()
			defer move.RUnlock()
			last = move
		}

		// 找到目标，直接返回即可
		if move.nexts[level] != nil && move.nexts[level].key == key {
			return move.nexts[level]
		}
	}
	return nil
}
