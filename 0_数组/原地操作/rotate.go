package main

func SwapRange(arr []int, i, j int) {
	for i < j {
		arr[i], arr[j] = arr[j], arr[i]
		i++
		j--
	}
}

func RotateLeft(arr []int, start, end, step int) {
	n := len(arr)
	if step >= n {
		step %= n
	}
	SwapRange(arr, start, start+step)
	SwapRange(arr, start+step, end)
	SwapRange(arr, start, end)
}

func RotateRight(arr []int, start, end, step int) {
	n := len(arr)
	if step >= n {
		step %= n
	}
	SwapRange(arr, start, end-step)
	SwapRange(arr, end-step, end)
	SwapRange(arr, start, end)
}
