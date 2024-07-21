package main

// !WARNING:Not Verified.
// 遍历row行col列的矩阵的对角线.
//
//	direction:
//	 - `0: ↘`, 从左上到右下. 同一条对角线上 `r-c` 为定值.
//	 - `1: ↖`, 从右下到左上. 同一条对角线上 `r-c` 为定值.
//	 - `2: ↙`, 从右上到左下. 同一条对角线上 `r+c` 为定值.
//	 - `3: ↗`, 从左下到右上. 同一条对角线上 `r+c` 为定值.
//
//	upToDown: 是否从上到下遍历.
func EnumerateDiagnal(row, col int, direction int8, upToDown bool, f func(group [][2]int)) {
	switch direction {
	case 0:
		if upToDown {
			for key := -col + 1; key < row; key++ {
				r := key
				if r < 0 {
					r = 0
				}
				c := r - key
				group := [][2]int{}
				for r < row && c < col {
					group = append(group, [2]int{r, c})
					r++
					c++
				}
				if len(group) > 0 {
					f(group)
				}
			}
		} else {
			for key := row - 1; key > -col; key-- {
				r := key
				if r < 0 {
					r = 0
				}
				c := r - key
				group := [][2]int{}
				for r < row && c < col {
					group = append(group, [2]int{r, c})
					r++
					c++
				}
				if len(group) > 0 {
					f(group)
				}
			}
		}
	case 1:
		if upToDown {
			for key := -col + 1; key < row; key++ {
				var r int
				if key > row-col {
					r = row - 1
				} else {
					r = key + col - 1
				}
				c := r - key
				group := [][2]int{}
				for r >= 0 && c >= 0 {
					group = append(group, [2]int{r, c})
					r--
					c--
				}
				if len(group) > 0 {
					f(group)
				}
			}
		} else {
			for key := row - 1; key > -col; key-- {
				var r int
				if key > row-col {
					r = row - 1
				} else {
					r = key + col - 1
				}
				c := r - key
				group := [][2]int{}
				for r >= 0 && c >= 0 {
					group = append(group, [2]int{r, c})
					r--
					c--
				}
				if len(group) > 0 {
					f(group)
				}
			}
		}
	case 2:
		if upToDown {
			for key := 0; key < row+col-1; key++ {
				var r int
				if key < col {
					r = 0
				} else {
					r = key - col + 1
				}
				c := key - r
				group := [][2]int{}
				for r < row && c >= 0 {
					group = append(group, [2]int{r, c})
					r++
					c--
				}
				if len(group) > 0 {
					f(group)
				}
			}
		} else {
			for key := row + col - 2; key >= 0; key-- {
				var r int
				if key < col {
					r = 0
				} else {
					r = key - col + 1
				}
				c := key - r
				group := [][2]int{}
				for r < row && c >= 0 {
					group = append(group, [2]int{r, c})
					r++
					c--
				}
				if len(group) > 0 {
					f(group)
				}
			}
		}
	case 3:
		if upToDown {
			for key := 0; key < row+col-1; key++ {
				var r int
				if key < row {
					r = key
				} else {
					r = row - 1
				}
				c := key - r
				group := [][2]int{}
				for r >= 0 && c < col {
					group = append(group, [2]int{r, c})
					r--
					c++
				}
				if len(group) > 0 {
					f(group)
				}
			}
		} else {
			for key := row + col - 2; key >= 0; key-- {
				var r int
				if key < row {
					r = key
				} else {
					r = row - 1
				}
				c := key - r
				group := [][2]int{}
				for r >= 0 && c < col {
					group = append(group, [2]int{r, c})
					r--
					c++
				}
				if len(group) > 0 {
					f(group)
				}
			}
		}
	default:
		panic("direction must be in (0, 1, 2, 3)")
	}
}
