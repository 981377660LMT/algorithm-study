// 扩展gcd求逆元

package main

import "fmt"

func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

func modInv(a, mod int) (inv int, err error) {
	gcd, x, _ := exgcd(a, mod)
	if gcd != 1 {
		return -1, fmt.Errorf("no inverse element")
	}
	return (x%mod + mod) % mod, nil
}

func main() {
	inv, err := modInv(3, 7)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(inv)
	}
}
