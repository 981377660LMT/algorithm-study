package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Println()
	var n int
	fmt.Fscan(in, &n)

	// 警惕带缓冲的输出
	// 这里的输出会被缓冲，所以不会立即输出
	// 会先执行 os.Exit(0)，导致输出被丢弃
	// print 会自动换行,所以每次输出都可见
	// !实现特判+退出需要用return 而不是像py一样用exit
	// 在 Python 中，标准输出（sys.stdout）`默认是以行缓冲模式工作的`(与操作系统油管)。
	// 这意味着，在碰到换行符之前，输出不会立即写入标准输出，
	// 而是先储存在缓冲区中。`当遇到换行符时，缓冲区中的内容会自动刷新到标准输出中`。
	// 如果你想立即输出而不等待换行符，你可以使用 sys.stdout.flush() 来手动刷新缓冲区。
	// 使用flush就不会输入到缓冲区，而是直接输出到标准输出中

	fmt.Fprintln(os.Stdout, n) // 可以输出
	fmt.Fprintln(out, n)       // !不会输出 是没有flush 带缓冲的输出要flush 否则会停留在缓冲区
	out.Flush()                // !这里才会输出
	os.Exit(0)
}
