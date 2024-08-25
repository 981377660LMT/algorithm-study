package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	abc368_c_1()
}

// C - Triple Attack
// https://atcoder.jp/contests/abc368/tasks/abc368_c_1
// !攻击怪兽。普通攻击-1，三连击-3，每三个回合可以使用一次三连击。怪兽血量为h[i]，问最少需要多少回合才能击败所有怪兽。
// 先处理循环节，然后处理剩余部分。
func abc368_c_1() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	turn := 0
	for _, v := range nums {
		cycle := v / 5
		turn += cycle * 3
		v -= cycle * 5
		for v > 0 {
			turn++
			if turn%3 == 0 {
				v -= 3
			} else {
				v--
			}
		}
	}

	fmt.Fprintln(out, turn)
}
