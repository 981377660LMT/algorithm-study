package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var A, B, M int
	fmt.Fscan(in, &A, &B, &M)
	N := A*B - 1

	size_LIS := A + 1
	size_LDS := B + 1
	size_minL := N + 2
	size_maxS := N + 1

	base1 := size_LDS * size_minL * size_maxS
	base2 := size_minL * size_maxS
	base3 := size_maxS

	dp_current := make([]int, size_LIS*size_LDS*size_minL*size_maxS)
	dp_new := make([]int, size_LIS*size_LDS*size_minL*size_maxS)

	init_index := 0*base1 + 0*base2 + (N+1)*base3 + 0
	dp_current[init_index] = 1

	for total := 0; total < N; total++ {
		for i := 0; i < len(dp_new); i++ {
			dp_new[i] = 0
		}

		for lis_len := 0; lis_len <= A; lis_len++ {
			for lds_len := 0; lds_len <= B; lds_len++ {
				for min_L := 1; min_L <= N+1; min_L++ {
					for max_S := 0; max_S <= N; max_S++ {
						idx := lis_len*base1 + lds_len*base2 + min_L*base3 + max_S
						count := dp_current[idx]
						if count == 0 {
							continue
						}

						rem_numbers := N - total
						choices := rem_numbers

						if lis_len+1 <= A {
							for x := 1; x < min_L; x++ {
								new_idx := (lis_len+1)*base1 + lds_len*base2 + x*base3 + max_S
								dp_new[new_idx] = (dp_new[new_idx] + count) % M
							}
						}

						if lds_len+1 <= B {
							for x := max_S + 1; x <= N; x++ {
								new_idx := lis_len*base1 + (lds_len+1)*base2 + min_L*base3 + x
								dp_new[new_idx] = (dp_new[new_idx] + count) % M
							}
						}

						add := choices - (min_L - 1) - (N - max_S)
						if add < 0 {
							add += M
						}
						add_val := (count * add) % M
						new_idx := lis_len*base1 + lds_len*base2 + min_L*base3 + max_S
						dp_new[new_idx] = (dp_new[new_idx] + add_val) % M
					}
				}
			}
		}

		dp_current, dp_new = dp_new, dp_current
	}

	answer := 0
	for min_L := 1; min_L <= N+1; min_L++ {
		for max_S := 0; max_S <= N; max_S++ {
			if min_L > max_S {
				idx := A*base1 + B*base2 + min_L*base3 + max_S
				answer = (answer + dp_current[idx]) % M
			}
		}
	}

	fmt.Fprintln(out, answer)
}
