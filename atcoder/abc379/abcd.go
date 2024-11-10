package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var N int
	fmt.Fscan(in, &N)
	S, _ := in.ReadString('\n')
	S = S[:len(S)-1]

	if len(S) < N {
		temp, _ := in.ReadString('\n')
		S += temp[:len(temp)-1]
	}

	totalSum := big.NewInt(0)
	multiplier := big.NewInt(0)
	powerOfTen := big.NewInt(1)
	ten := big.NewInt(10)

	for i := N - 1; i >= 0; i-- {
		digit := int64(S[i] - '0')
		num := big.NewInt(digit)

		multiplier.Add(multiplier, powerOfTen)

		index := big.NewInt(int64(i + 1))

		temp := new(big.Int).Mul(num, index)
		temp.Mul(temp, multiplier)

		totalSum.Add(totalSum, temp)

		powerOfTen.Mul(powerOfTen, ten)
	}

	fmt.Println(totalSum.String())
}
