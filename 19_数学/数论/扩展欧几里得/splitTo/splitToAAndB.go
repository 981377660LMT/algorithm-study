package main

func main() {

}

func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

func solveLinearDiophantineEquations(a, b, c int) (n, x1, y1, x2, y2 int) {
	g, x0, y0 := exgcd(a, b)

	// 无解
	if c%g != 0 {
		n = -1
		return
	}

	a /= g
	b /= g
	c /= g
	x0 *= c
	y0 *= c

	x1 = x0 % b
	if x1 <= 0 { // 若要求的是非负整数解，去掉等号
		x1 += b
	}
	k1 := (x1 - x0) / b
	y1 = y0 - k1*a

	y2 = y0 % a
	if y2 <= 0 { // 若要求的是非负整数解，去掉等号
		y2 += a
	}
	k2 := (y0 - y2) / a
	x2 = x0 + k2*b

	// 无正整数解
	if y1 <= 0 {
		return
	}

	// k 越大 x 越大
	n = k2 - k1 + 1
	return
}
