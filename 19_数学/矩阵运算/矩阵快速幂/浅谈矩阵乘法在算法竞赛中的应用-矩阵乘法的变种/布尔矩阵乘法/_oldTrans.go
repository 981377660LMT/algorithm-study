// package main

// import "math"

// func main() {

// }

// // 原地修改a求传递闭包.
// func TransitiveClosure(a []int, l, r int, n int) {
// 	if l+1 == r {
// 		return
// 	}
// 	mid := (l + r + 1) >> 1
// 	n1, n2 := mid-l, r-mid
// 	TransitiveClosure(a, l, mid, n)
// 	TransitiveClosure(a, mid, r, n)
// 	size := n * n
// 	b, c, d := make([]int, size), make([]int, size), make([]int, size)
// 	for i := 0; i < n1; i++ {
// 		for j := 0; j < n1; j++ {
// 			b[i*n1+j] = a[(l+i)*n+(l+j)]
// 		}
// 	}
// 	for i := 0; i < n1; i++ {
// 		for j := 0; j < n2; j++ {
// 			c[i*n1+j] = a[(l+i)*n+(mid+j)]
// 		}
// 		if n1 > n2 {
// 			c[i*n1+n2] = 0
// 		}
// 	}
// 	MatMul(b, c, d, n1)
// 	for i := 0; i < n2; i++ {
// 		for j := 0; j < n2; j++ {
// 			c[i*n1+j] = a[(mid+i)*n+(mid+j)]
// 		}
// 		if n1 > n2 {
// 			c[i*n1+n2] = 0
// 			c[n2*n1+i] = 0
// 		}
// 	}
// 	MatMul(d, c, b, n1)
// 	for i := 0; i < n1; i++ {
// 		for j := 0; j < n2; j++ {
// 			a[(l+i)*n+(mid+j)] = b[i*n1+j]
// 		}
// 	}
// }

// func MatMul(a, b, c []int, n int) {
// 	if n <= 4 {
// 		MatMulBf(a, b, c, n)
// 		return
// 	}
// 	if n <= 64 {
// 		MatMulWord(a, b, c, n)
// 		return
// 	}
// 	f := make([][][]uint64, n/64+1)
// 	B := make([][]uint64, n/3+1)
// 	res := make([][]uint64, n)
// 	// int L=max((int)floor(log2(max(n/W,1)))-1,1),n1=(n-1)/L+1,n2=(n-1)/W+1;
// 	L := max(int(math.Floor(math.Log2(float64(n/64)))), 1)
// }

// func MatMulWord(a, b, c []int, n int) {
// 	size := 8 * n
// 	A, B := make([]int, size), make([]int, size)
// 	for i := 0; i < n; i++ {
// 		for j := 0; j < n; j++ {
// 			A[i] |= a[i*n+j] << j
// 			B[j] |= b[i*n+j] << i
// 		}
// 	}
// 	for i := 0; i < n; i++ {
// 		ai := A[i]
// 		for j := 0; j < n; j++ {
// 			tmp := ai & B[j]
// 			if tmp > 0 {
// 				c[i*n+j] = 1
// 			}
// 		}
// 	}
// }

// func MatMulBf(a, b, c []int, n int) {
// 	for i := 0; i < n; i++ {
// 		for k := 0; k < n; k++ {
// 			if a[i*n+k] > 0 {
// 				offsetB := k * n
// 				offsetC := i * n
// 				for j := 0; j < n; j++ {
// 					c[offsetC+j] |= b[offsetB+j]
// 				}
// 			}
// 		}
// 	}
// }

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }
// https://leetcode.cn/problems/course-schedule-iv/solution/chuan-di-bi-bao-gui-yue-dao-ju-zhen-chen-47yu/

package main

func main() {

}
