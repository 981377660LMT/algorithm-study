package main

func main() {

}

func WalshHadamardTransform(f []int, inverse bool) {
	n := len(f)
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				x, y := f[j], f[j|i]
				f[j], f[j|i] = x+y, x-y
			}
		}
	}

	if inverse {
		for i := range f {
			f[i] /= n
		}
	}
}
