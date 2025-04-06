// 近似的LRU缓存

package main

import "fmt"

type Page struct {
	ID           int
	ReferenceBit bool
}

type ClockBuffer struct {
	frames   []Page      // 缓冲区帧数组
	capacity int         // 缓冲区容量
	hand     int         // 时钟指针位置
	pageMap  map[int]int // 页面ID到帧索引的映射
}

func NewClockBuffer(capacity int) *ClockBuffer {
	frames := make([]Page, capacity)
	for i := range frames {
		frames[i] = Page{ID: -1}
	}
	return &ClockBuffer{
		frames:   frames,
		capacity: capacity,
		hand:     0,
		pageMap:  make(map[int]int),
	}
}

// 访问页面，设置引用位为1
func (c *ClockBuffer) Access(pageID int) bool {
	if idx, exists := c.pageMap[pageID]; exists {
		c.frames[idx].ReferenceBit = true
		return true
	}
	return false
}

// 添加页面到缓冲区
func (c *ClockBuffer) Add(pageID int) (evictedID int, evicted bool) {
	// 如果页面已在缓冲区中，只需设置引用位
	if c.Access(pageID) {
		return -1, false
	}

	// 查找空闲帧
	for i, frame := range c.frames {
		if frame.ID == -1 {
			c.frames[i] = Page{ID: pageID, ReferenceBit: true}
			c.pageMap[pageID] = i
			return -1, false
		}
	}

	// 缓冲区已满，需要替换页面
	// 使用CLOCK算法寻找要替换的页面
	for {
		if !c.frames[c.hand].ReferenceBit {
			// 找到引用位为0的页面，将其替换
			evictedID = c.frames[c.hand].ID
			delete(c.pageMap, evictedID)

			// 添加新页面
			c.frames[c.hand] = Page{ID: pageID, ReferenceBit: true}
			c.pageMap[pageID] = c.hand

			// 移动时钟指针
			c.hand = (c.hand + 1) % c.capacity

			return evictedID, true
		}

		// 引用位为1，设置为0并移动指针
		c.frames[c.hand].ReferenceBit = false
		c.hand = (c.hand + 1) % c.capacity
	}
}

// 演示CLOCK算法的运行
func main() {
	// 创建容量为4的缓冲池
	buffer := NewClockBuffer(4)

	fmt.Println("添加页面序列: 1, 2, 3, 4")
	buffer.Add(1)
	buffer.Add(2)
	buffer.Add(3)
	buffer.Add(4)
	printBufferState(buffer)

	fmt.Println("\n访问页面2")
	buffer.Access(2)
	printBufferState(buffer)

	fmt.Println("\n添加新页面5（触发替换）")
	evictedID, _ := buffer.Add(5)
	fmt.Printf("被替换的页面: %d\n", evictedID)
	printBufferState(buffer)

	fmt.Println("\n添加新页面6（触发替换）")
	evictedID, _ = buffer.Add(6)
	fmt.Printf("被替换的页面: %d\n", evictedID)
	printBufferState(buffer)
}

// 打印缓冲池状态
func printBufferState(c *ClockBuffer) {
	fmt.Printf("时钟指针位置: %d\n", c.hand)
	fmt.Println("缓冲池状态:")
	for i, page := range c.frames {
		fmt.Printf("[%d] ID=%d, 引用位=%v\n", i, page.ID, page.ReferenceBit)
	}
}
