// lcsbit/bitlcs
// https://atcoder.jp/contests/dp/submissions/34604402
// https://loj.ac/s/1633431

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s, &t)
	fmt.Fprintln(out, strings.Join(GetLCSString(s, t), ""))
}

func perf() {
	n1, n2 := int(1e5), int(1e5)
	nums1, nums2 := make([]int, n1), make([]int, n2)
	for i := 0; i < n1; i++ {
		nums1[i] = i
	}
	for i := 0; i < n2; i++ {
		nums2[i] = i
	}
	time1 := time.Now()
	GetLCS(nums1, nums2)
	fmt.Println(time.Since(time1))
}

func GetLCSString(s, t string) []string {
	ords1, ords2 := make([]int, len(s)), make([]int, len(t))
	for i := range s {
		ords1[i] = int(s[i])
	}
	for i := range t {
		ords2[i] = int(t[i])
	}
	ords := GetLCS(ords1, ords2)
	res := make([]string, 0, len(ords))
	for _, v := range ords {
		res = append(res, fmt.Sprintf("%c", v))
	}
	return res

}

func GetLCS(nums1, nums2 []int) []int {
	nums1, nums2 = append(nums1[:0:0], nums1...), append(nums2[:0:0], nums2...)
	id := make(map[int]int)
	for i, v := range nums1 {
		if value, ok := id[v]; ok {
			nums1[i] = value
		} else {
			id[v] = len(id)
			nums1[i] = len(id) - 1
		}
	}
	for i, v := range nums2 {
		if value, ok := id[v]; ok {
			nums2[i] = value
		} else {
			id[v] = len(id)
			nums2[i] = len(id) - 1
		}
	}
	rid := make([]int, len(id))
	for k, v := range id {
		rid[v] = k
	}

	n, m := len(nums1), len(nums2)
	sets := make([]*BS, len(id))
	for i := range sets {
		sets[i] = NewBS(n + 1)
	}
	dp := make([]*BS, m+1)
	for i := range dp {
		dp[i] = NewBS(n + 1)
	}
	tmp := NewBS(n + 1)

	for i, v := range nums1 {
		sets[v].Set(i + 1)
	}

	for i := 1; i <= m; i++ {
		dp[i].SetAll(dp[i-1])
		tmp.SetAll(dp[i])
		tmp.OrAll(sets[nums2[i-1]])
		dp[i].Shift()
		dp[i].Set(0)
		dp[i].SetAll(Minus(tmp, dp[i]))
		dp[i].XorAll(tmp)
		dp[i].AndAll(tmp)
	}

	i, j := n, m
	res := []int{}
	for i > 0 && j > 0 {
		if nums1[i-1] == nums2[j-1] {
			res = append(res, nums1[i-1])
			i--
			j--
		} else if dp[j].Get(i) == false {
			i--
		} else {
			j--
		}
	}

	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	for i, v := range res {
		res[i] = rid[v]
	}
	return res
}

type BS struct {
	data []uint64
}

func NewBS(n int) *BS {
	return &BS{data: make([]uint64, 1+n>>6)}
}

func (bs *BS) Set(i int) *BS {
	bs.data[i>>6] |= 1 << (i & 63)
	return bs
}

func (bs *BS) Get(i int) bool {
	return bs.data[i>>6]&(1<<(i&63)) != 0
}

// 所有位向高位移动一位.
//
//	[1,2,3] -> [2,3,4]
func (bs *BS) Shift() *BS {
	last := uint64(0)
	for i := range bs.data {
		cur := (bs.data[i] >> 63) // 保留最高位
		bs.data[i] <<= 1          // 所有位向高位移动一位
		bs.data[i] |= last        // 将上一位的最高位放到当前位的最低位
		last = cur
	}
	return bs
}

func (bs *BS) SetAll(rhs *BS) *BS {
	copy(bs.data, rhs.data)
	return bs
}

func (bs *BS) AndAll(rhs *BS) *BS {
	for i := range bs.data {
		bs.data[i] &= rhs.data[i]
	}
	return bs
}

func (bs *BS) OrAll(rhs *BS) *BS {
	for i := range bs.data {
		bs.data[i] |= rhs.data[i]
	}
	return bs
}

func (bs *BS) XorAll(rhs *BS) *BS {
	for i := range bs.data {
		bs.data[i] ^= rhs.data[i]
	}
	return bs
}

func (bs *BS) String() string {
	sb := strings.Builder{}
	sb.Write([]byte{'['})
	nums := strings.Builder{}
	for i := range bs.data {
		start, end := i<<6, (i+1)<<6
		for j := start; j < end; j++ {
			if bs.Get(j) {
				nums.WriteString(fmt.Sprintf("%d,", j))
			}
		}
	}
	sb.WriteString(nums.String()[:nums.Len()-1])
	sb.Write([]byte{']'})
	return sb.String()
}

// Minus 计算两个集合的差集.
//
//	注意要用大的(lhs)减小的(rhs).
func Minus(lhs, rhs *BS) *BS {
	last := uint64(0)
	res := &BS{data: make([]uint64, len(lhs.data))}
	for i := range lhs.data {
		cur := (lhs.data[i] < (rhs.data[i] + last))
		res.data[i] = lhs.data[i] - rhs.data[i] - last
		if cur {
			last = 1
		} else {
			last = 0
		}
	}

	return res
}
