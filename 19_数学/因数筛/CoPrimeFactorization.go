// CoPrimeFactorization

package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(CoPrimeFactorization([]int{21, 60, 140, 400}))
}

// 给定一组互质的数,求出它们的最大公共基底和对应分解的结果.
// eg: [21,60,140,400] -> [3,7,20], [[(0,1),(1,1)],[(0,1),(2,1)],[(1,1),(2,1)],[(2,2)]
func CoPrimeFactorization(coPrimes []int) (basis []int, exps [][][2]int) {
	for _, val := range coPrimes {
		newBasis := []int{}
		for _, x := range basis {
			if val == 1 {
				newBasis = append(newBasis, x)
				continue
			}
			dat := []int{val, x}
			for p := 1; p < len(dat); p++ {
				for i := 0; i < p; i++ {
					for {
						if dat[p] > 1 && dat[i]%dat[p] == 0 {
							dat[i] /= dat[p]
						} else if dat[i] > 1 && dat[p]%dat[i] == 0 {
							dat[p] /= dat[i]
						} else {
							break
						}
					}
					g := gcd(dat[i], dat[p])
					if g == 1 || g == dat[i] || g == dat[p] {
						continue
					}
					dat[i] /= g
					dat[p] /= g
					dat = append(dat, g)
				}
			}
			val = dat[0]
			for i := 1; i < len(dat); i++ {
				if dat[i] != 1 {
					newBasis = append(newBasis, dat[i])
				}
			}
		}
		if val > 1 {
			newBasis = append(newBasis, val)
		}
		basis = newBasis
	}

	sort.Ints(basis)

	exps = make([][][2]int, len(coPrimes))
	for i, p := range coPrimes {
		for j, b := range basis {
			e := 0
			for p%b == 0 {
				p /= b
				e++
			}
			if e > 0 {
				exps[i] = append(exps[i], [2]int{j, e})
			}
		}
	}
	return
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
