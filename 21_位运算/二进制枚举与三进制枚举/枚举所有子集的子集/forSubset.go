// forSubset枚举某个状态的所有子集(枚举子集的子集)

package main

import "fmt"

func main() {
	state := 0b1101
	for g1 := state; g1 >= 0; {
		if g1 == state || g1 == 0 { // 排除空集和全集
			g1--
			continue
		}
		g2 := state ^ g1
		fmt.Println(g1, g2)
		if g1 == 0 {
			g1 = -1
		} else {
			g1 = (g1 - 1) & state
		}
	}
}
