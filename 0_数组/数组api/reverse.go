package main

import "fmt"

func main() {
	s1 := "hello"
	fmt.Println(ReverseString(s1)) // olleh

	s2 := []int{1, 2, 3, 4, 5}
	Reverse(s2)
	fmt.Println(s2) // [5 4 3 2 1]

	s3 := []int{1, 2, 3, 4, 5}
	ReverseFunc(len(s3), func(i int) *int { return &s3[i] })
	fmt.Println(s3) // [5 4 3 2 1]
}

func ReverseString(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return string(runes)
}

func Reverse[S []E, E any](arr S) {
	n := len(arr)
	for i := 0; i < n/2; i++ {
		arr[i], arr[n-1-i] = arr[n-1-i], arr[i]
	}
}

func ReverseFunc[E any](n int, ptr func(i int) *E) {
	for i := 0; i < n/2; i++ {
		a, b := ptr(i), ptr(n-1-i)
		*a, *b = *b, *a
	}
}
