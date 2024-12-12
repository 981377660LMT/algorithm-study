// 使用跳转标签让break和continue语句更简洁
// !Loop labels for cleaner breaks and continues

package main

import "fmt"

func main() {
	// gotoDemo1()
	gotoDemo2()
}

func gotoDemo1() {
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	target := 5

OuterLoop:
	for i, row := range matrix {
		for j, val := range row {
			if val == target {
				fmt.Printf("found %d (row: %d, col: %d)\n", target, i, j)
				break OuterLoop
			}

		}
	}
}

func gotoDemo2() {
	userChoice := 2
SwitchChoice:
	switch userChoice {
	case 1:
		fmt.Println("You chose 1")
	case 2:
		fmt.Println("You chose 2")
		for i := 0; i < 10; i++ {
			fmt.Println(i)
			if i == 5 {
				break SwitchChoice
			}
		}
	default:
		fmt.Println("You chose something else")
	}
}
