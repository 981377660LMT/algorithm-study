/* 构造题
n 个小朋友围成一圈，按顺时针顺序依次编号为 1∼n。
有 7 种颜色的帽子，每种颜色的帽子的数量都足够多。
7 种颜色不妨表示为 R、O、Y、G、B、I、V。
现在，要给每个小朋友都发一个帽子，要求：

!每种颜色的帽子都至少有一个小朋友戴。
!任意四个相邻小朋友的帽子颜色都各不相同。
请你提供一种分发帽子的方案。

!把7种颜色分成两组 0123 456
!初始为456 接下来0123循环即可
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const G1 = "ROY"
const G2 = "GBIV"

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	res := []string{G1}
	div, mod := (n-3)/4, (n-3)%4
	res = append(res, strings.Repeat(G2, div), G2[:mod])

	fmt.Fprintln(out, strings.Join(res, ""))
}
