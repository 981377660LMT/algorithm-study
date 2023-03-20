// (https://github.com/EndlessCheng/codeforces-go/blob/029f576c04914ad4052dbe1073ff644dc219824a/copypasta/common.go#L390)

package main

import "fmt"

func main() {
	fmt.Println(mul64(2, 3, 1e10+7))
}

// 适用于 mod 超过 int32 范围的情况
func mul64(a, b, mod int64) (res int64) {
	for ; b > 0; b >>= 1 {
		if b&1 == 1 {
			res = (res + a) % mod
		}
		a = (a + a) % mod
	}
	return
}
