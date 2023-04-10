// https://yukicoder.me/problems/116

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// 结束可以抢夺牌的nim游戏,拿到最后一张牌就胜利
	// 場に４つの山が用意される。
	// 4つの山それぞれの枚数が与えられます(1<=n<=13)。
	// プレイヤーは交互に山からカードを取り手札とする
	// !１回に取れるのは１つの山から１枚～３枚のみ（複数の山からまとめてとることはできない）
	// !パスはできず必ず１枚はカードを取らなければならない
	// ４つの山それぞれについて、最後のカードを取った場合、相手の手札の半分（奇数枚の場合は切り上げ）を奪うことができる
	// すべての山が無くなったとき、手札が多いほうが勝ち
	// 勝利するプレイヤー名（"Taro"か"Jiro"）、引き分けの場合は"Draw"を出力してください。

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	nums := [4]int{0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		fmt.Fscan(in, &nums[i])
	}

	dag := make([][]int, 14)
	for i := 1; i <= 13; i++ {
		for j := 1; j <= 3; j++ {
			if i-j >= 0 {
				dag[i] = append(dag[i], i-j)
			}
		}
	}
	grundy := GrundyNumber(dag) // !可以发现个数为i时,grundy数为 i%(k+1)

	xor := 0
	for i := 0; i < 4; i++ {
		xor ^= grundy[nums[i]]
	}

	if xor == 0 {
		fmt.Fprintln(out, "Jiro")
	} else {
		fmt.Fprintln(out, "Taro")
	}
}

// dag: 博弈的每个状态组成的有向无环图.
//  返回值: 每个状态的Grundy数.
//  grundy[i] = mex{grundy[j] | j in dag[i]}.
//  - 如果grundy为0,则先手必败,否则先手必胜.
//  - 若一个母状态可以拆分成多个相互独立的子状态，`则母状态的 SG 数等于各个子状态的 SG 数的异或。`
func GrundyNumber(dag [][]int) (grundy []int) {
	order, ok := topoSort(dag)
	if !ok {
		return
	}

	grundy = make([]int, len(dag))
	memo := make([]int, len(dag)+1)
	for j := len(order) - 1; j >= 0; j-- {
		i := order[j]
		if len(dag[i]) == 0 {
			continue
		}
		for _, v := range dag[i] {
			memo[grundy[v]]++
		}
		for memo[grundy[i]] > 0 {
			grundy[i]++
		}
		for _, v := range dag[i] {
			memo[grundy[v]]--
		}
	}

	return
}

func topoSort(dag [][]int) (order []int, ok bool) {
	n := len(dag)
	visited, temp := make([]bool, n), make([]bool, n)
	var dfs func(int) bool
	dfs = func(i int) bool {
		if temp[i] {
			return false
		}
		if !visited[i] {
			temp[i] = true
			for _, v := range dag[i] {
				if !dfs(v) {
					return false
				}
			}
			visited[i] = true
			order = append(order, i)
			temp[i] = false
		}
		return true
	}

	for i := 0; i < n; i++ {
		if !visited[i] {
			if !dfs(i) {
				return nil, false
			}
		}
	}

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order, true
}
