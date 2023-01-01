// golang输入/golang输出
// 关键词:golangIO golang输入输出

package foo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

// https://github.dev/EndlessCheng/codeforces-go
// 带有 IO 缓冲区的输入输出
func bufferIO(reader io.Reader, writer io.Writer) {
	in := bufio.NewReader(reader)
	out := bufio.NewWriter(writer)
	defer out.Flush()

	// 读入一个整数
	var n int
	fmt.Fscan(in, &n)

	// 读入一个整数数组
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	// 读入一行字符串
	var s string
	fmt.Fscanln(in, &s)

	// 读入一个字符串数组
	strs := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscanln(in, &strs[i])
	}

	fmt.Fprintln(out, n)

	// 换行输出数组
	for _, v := range nums {
		fmt.Fprintln(out, v)
	}

	// 不换行输出数组
	for _, v := range nums {
		fmt.Fprint(out, v, " ")
	}
}

// 按行读入
func lineIO(reader io.Reader, writer io.Writer) {
	in := bufio.NewScanner(reader)
	in.Buffer(nil, 1e9) // 若单个 token 大小超过 65536 则加上这行
	out := bufio.NewWriter(writer)
	defer out.Flush()

	for in.Scan() {
		line := in.Bytes()
		sp := bytes.Split(line, []byte{' '})
		// do something
		for _, v := range sp {
			fmt.Println(string(v))
		}
	}
}

func main() {
	bufferIO(os.Stdin, os.Stdout)
}
