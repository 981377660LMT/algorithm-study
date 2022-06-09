package foo

import (
	"bufio"
	"fmt"
	"io"
)

// https://github.dev/EndlessCheng/codeforces-go
// 带有 IO 缓冲区的输入输出
func bufferIO(reader io.Reader, writer io.Writer) {
	in := bufio.NewReader(reader)
	out := bufio.NewWriter(writer)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	fmt.Fprintln(out, n)
}
