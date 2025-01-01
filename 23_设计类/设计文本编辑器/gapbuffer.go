// Gap Buffer 适用于需要在局部频繁编辑文本的应用场景，如文本编辑器和集成开发环境。通过合理管理间隙位置和缓冲区容量，可以实现高效的插入和删除操作。
// 然而，在处理大规模文本或频繁的非局部编辑时，Gap Buffer 的性能可能不如其他数据结构（如 Rope 或 Piece Table），因此在实际应用中需要根据具体需求选择合适的数据结构或结合多种优化策略。

package main

import (
	"fmt"
	"strings"
)

// GapBuffer 结构体表示间隙缓冲区
type GapBuffer struct {
	buffer   []rune // 存储文本的缓冲区
	gapStart int    // 间隙起始位置
	gapEnd   int    // 间隙结束位置
	capacity int    // 缓冲区总容量
}

// NewGapBuffer 初始化一个新的 GapBuffer
// initialContent 为初始文本内容
func NewGapBuffer(initialContent string) *GapBuffer {
	// 将初始内容转换为 rune 切片以支持 Unicode
	runes := []rune(initialContent)
	initialCapacity := len(runes) + 10 // 初始缓冲区容量，比内容稍大以留出间隙
	buffer := make([]rune, initialCapacity)
	copy(buffer, runes)
	gapStart := len(runes)
	gapEnd := initialCapacity
	return &GapBuffer{
		buffer:   buffer,
		gapStart: gapStart,
		gapEnd:   gapEnd,
		capacity: initialCapacity,
	}
}

// Length 返回文本的当前长度
func (gb *GapBuffer) Length() int {
	return gb.capacity - (gb.gapEnd - gb.gapStart)
}

// Insert 在当前位置插入一个字符
func (gb *GapBuffer) Insert(char rune) {
	if gb.gapStart == gb.gapEnd {
		gb.expandBuffer()
	}
	gb.buffer[gb.gapStart] = char
	gb.gapStart++
}

// Delete 删除当前位置的字符（向后删除）
func (gb *GapBuffer) Delete() {
	if gb.gapEnd == gb.gapStart {
		// Gap 已经覆盖末尾，无字符可删
		return
	}
	gb.gapEnd++
}

// MoveGap 移动间隙到指定位置
func (gb *GapBuffer) MoveGap(position int) {
	if position < 0 || position > gb.Length() {
		panic("MoveGap: position out of bounds")
	}

	if position < gb.gapStart {
		// 需要将间隙向左移动
		moveSize := gb.gapStart - position
		copy(gb.buffer[gb.gapEnd-moveSize:gb.gapEnd], gb.buffer[position:gb.gapStart])
		gb.gapStart = position
		gb.gapEnd -= moveSize
	} else if position > gb.gapStart {
		// 需要将间隙向右移动
		moveSize := position - gb.gapStart
		copy(gb.buffer[gb.gapStart:gb.gapStart+moveSize], gb.buffer[gb.gapEnd:gb.gapEnd+moveSize])
		gb.gapStart += moveSize
		gb.gapEnd += moveSize
	}
	// 如果 position == gapStart，间隙已在目标位置，无需移动
}

// GetContent 获取当前缓冲区的全部内容
func (gb *GapBuffer) GetContent() string {
	var sb strings.Builder
	sb.Grow(gb.Length())
	sb.WriteString(string(gb.buffer[:gb.gapStart]))
	sb.WriteString(string(gb.buffer[gb.gapEnd:gb.capacity]))
	return sb.String()
}

// expandBuffer 扩展缓冲区容量，并调整间隙大小
func (gb *GapBuffer) expandBuffer() {
	// 扩展策略：将缓冲区容量翻倍，并在间隙中添加新的空闲空间
	newCapacity := gb.capacity * 2
	newBuffer := make([]rune, newCapacity)

	// 复制前部分
	copy(newBuffer, gb.buffer[:gb.gapStart])

	// 设置新的 gapStart 和 gapEnd
	gb.gapEnd = newCapacity - (gb.capacity - gb.gapEnd)
	gb.buffer = newBuffer
	gb.capacity = newCapacity
}

// Example usage
func main() {
	// 初始化 GapBuffer
	gb := NewGapBuffer("Hello World")
	fmt.Println("Initial content:", gb.GetContent()) // 输出: Hello World

	// 插入字符
	gb.Insert('!')
	fmt.Println("After insertion:", gb.GetContent()) // 输出: Hello World!

	// 移动间隙到位置 5（光标移动到 'Hello| World!' 中的 | 位置）
	gb.MoveGap(5)

	// 插入另一个字符
	gb.Insert(',')
	fmt.Println("After moving gap and inserting:", gb.GetContent()) // 输出: Hello, World!

	// 删除字符（删除空格）
	gb.Delete()
	fmt.Println("After deletion:", gb.GetContent()) // 输出: Hello,World!

	// 插入多个字符
	message := " GoLang"
	for _, char := range message {
		gb.Insert(char)
	}
	fmt.Println("After bulk insertion:", gb.GetContent()) // 输出: Hello,World! GoLang
}
