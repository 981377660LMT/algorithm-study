package main

import (
	"container/heap"
	"fmt"
)

// 访问历史记录
type HistoryItem struct {
	timestamp int64 // 访问时间戳
}

// 页面结构
type Page struct {
	id          int           // 页面ID
	history     []HistoryItem // 保存K次访问历史
	accessCount int           // 访问计数
	lastKAccess int64         // 第K次最近访问的时间戳
	correlated  bool          // 是否有足够访问历史做相关性预测
}

// LRUK缓冲池
type LRUKBuffer struct {
	capacity    int           // 缓冲池容量
	k           int           // K值
	pages       map[int]*Page // 页面映射
	currentTime int64         // 当前逻辑时间
	pageQueue   PriorityQueue // 优先队列用于有效替换
}

// 创建新的LRUK缓冲池
func NewLRUKBuffer(capacity, k int) *LRUKBuffer {
	buffer := &LRUKBuffer{
		capacity:    capacity,
		k:           k,
		pages:       make(map[int]*Page),
		currentTime: 0,
		pageQueue:   make(PriorityQueue, 0),
	}
	heap.Init(&buffer.pageQueue)
	return buffer
}

// 访问页面（记录历史）
func (b *LRUKBuffer) Access(pageID int) {
	b.currentTime++

	// 如果页面已存在，更新其访问历史
	if page, exists := b.pages[pageID]; exists {
		// 添加新的访问记录
		b.updateHistory(page)

		// 更新优先队列
		heap.Fix(&b.pageQueue, page.accessCount-1)
		return
	}

	// 页面不存在，需要添加
	if len(b.pages) >= b.capacity {
		// 缓冲池已满，需要替换页面
		b.evictPage()
	}

	// 创建新页面并添加到缓冲池
	page := &Page{
		id:          pageID,
		history:     make([]HistoryItem, 0, b.k),
		accessCount: 0,
		lastKAccess: -1,
		correlated:  false,
	}
	b.pages[pageID] = page
	b.updateHistory(page)

	// 将页面添加到优先队列
	item := &Item{
		page:     page,
		priority: b.calculatePriority(page),
		index:    len(b.pageQueue),
	}
	heap.Push(&b.pageQueue, item)
}

// 更新页面访问历史
func (b *LRUKBuffer) updateHistory(page *Page) {
	// 记录当前访问
	page.accessCount++

	// 添加访问历史
	historyItem := HistoryItem{timestamp: b.currentTime}

	if len(page.history) < b.k {
		page.history = append(page.history, historyItem)
	} else {
		// 移除最旧的历史记录，添加新的
		page.history = append(page.history[1:], historyItem)
	}

	// 如果有足够的历史，更新lastKAccess
	if len(page.history) == b.k {
		page.lastKAccess = page.history[0].timestamp
		page.correlated = true
	}
}

// 驱逐页面
func (b *LRUKBuffer) evictPage() {
	if b.pageQueue.Len() == 0 {
		return
	}

	// 从优先队列中移除优先级最低的页面
	item := heap.Pop(&b.pageQueue).(*Item)
	evictedPage := item.page

	// 从缓冲池中移除页面
	delete(b.pages, evictedPage.id)

	fmt.Printf("驱逐页面: %d (访问次数: %d, 优先级: %f)\n",
		evictedPage.id, evictedPage.accessCount, item.priority)
}

// 计算页面优先级（用于替换决策）
func (b *LRUKBuffer) calculatePriority(page *Page) float64 {
	// 如果页面访问次数少于K次
	if !page.correlated {
		// 使用FIFO策略（基于首次访问时间）
		if len(page.history) > 0 {
			return float64(page.history[0].timestamp)
		}
		return float64(b.currentTime)
	}

	// 使用第K次最近访问的时间戳（较旧的优先级更高）
	return float64(page.lastKAccess)
}

//---- 优先队列实现 ----

type Item struct {
	page     *Page
	priority float64 // 优先级
	index    int     // 在堆中的索引
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// !较小的优先级值意味着较高的驱逐优先级
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// 演示LRU-K算法
func main() {
	// 创建容量为3，K=2的LRU-K缓冲池
	buffer := NewLRUKBuffer(3, 2)

	// 访问序列：1, 2, 3, 4, 1, 2, 5
	fmt.Println("访问页面 1")
	buffer.Access(1)

	fmt.Println("访问页面 2")
	buffer.Access(2)

	fmt.Println("访问页面 3")
	buffer.Access(3)

	fmt.Println("访问页面 4 (触发替换)")
	buffer.Access(4)

	fmt.Println("访问页面 1 (再次)")
	buffer.Access(1)

	fmt.Println("访问页面 2 (再次)")
	buffer.Access(2)

	fmt.Println("访问页面 5 (触发替换)")
	buffer.Access(5)

	fmt.Println("\n最终缓冲池状态:")
	for id, page := range buffer.pages {
		fmt.Printf("页面 %d: 访问次数=%d, 第K次访问时间=%d\n",
			id, page.accessCount, page.lastKAccess)
	}
}
