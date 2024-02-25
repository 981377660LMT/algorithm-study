// API:
//  RollingHash(base) : 以base为底的哈希函数.
//  Build(s) : O(n)返回s的哈希值表.
//  Query(sTable, l, r) : 返回s[l, r)的哈希值.
//  Combine(h1, h2, h2len) : 哈希值h1和h2的拼接, h2的长度为h2len. (线段树的op操作)
//  AddChar(h,c) : 哈希值h和字符c拼接形成的哈希值.
//  Lcp(aTable, l1, r1, bTable, l2, r2) : O(logn)返回a[l1, r1)和b[l2, r2)的最长公共前缀.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF1056E()
}

// 单词映射
// https://www.luogu.com.cn/problem/CF1056E
// 我们有一个由0，1组成的数字串s，还有一个字符串t。
// 我们令s中的每一个0都对应一个串s1，令s中每一个1都对应另一个串s2，s1与s2不能完全相同。
// 问有多少组不同的串对(s1,s2)，使得s被s1，s2替换后后与t完全相同。
// 保证s中至少有一个0和一个1。
// !eg: s='001',t='kokokokotlin' => 则0=>ko,1=>kotlin或者0=>koko,1=>tlin
//
// s 串的开头那个字符一定对应的是 t 串的一个前缀。因此我们只需要枚举这个字符对应串的长度就可以唯一确定它对应的串。
func CF1056E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var source, target string
	fmt.Fscan(in, &source, &target)

	zero, one := 0, 0
	for i := 0; i < len(source); i++ {
		if source[i] == '0' {
			zero++
		} else {
			one++
		}
	}

	R1 := NewRollingHash(131, 1e9+7)
	table1 := R1.Build(target)
	getHash := func(start, end int) uint {
		return R1.Query(table1, start, end)
	}

	n := len(target)
	res := 0

	check := func(x, y int) bool {
		var h1, h2 []uint
		ptr := 0
		for i := 0; i < len(source); i++ {
			if source[i] == '0' {
				curHash := getHash(ptr, ptr+x)
				if len(h1) > 0 && h1[0] != curHash {
					return false
				}
				h1 = append(h1, curHash)
				ptr += x
			} else {
				curHash := getHash(ptr, ptr+y)
				if len(h2) > 0 && h2[0] != curHash {
					return false
				}
				h2 = append(h2, curHash)
				ptr += y
			}
		}
		if ptr != n {
			panic("ptr != n")
		}
		return h1[0] != h2[0]
	}

	for x := 1; x < n; x++ { // 枚举长度
		y := (n - zero*x) / one
		if y <= 0 {
			break
		}
		if zero*x+one*y != n { // 不能整除
			continue
		}

		if check(x, y) {
			res++
		}
	}

	fmt.Fprintln(out, res)

}

// 顺序遍历每个单词,问之前是否见过类似的单词.
//
//	`类似`: 最多交换一次相邻字符,可以得到相同的单词.
//	n<=2e5 sum(len(s))<=1e6
func ConditionalReflection(words []string) []bool {

	R := NewRollingHash(37, 2102001800968)
	res := make([]bool, len(words))
	visited := make(map[uint]struct{})

	for i := 0; i < len(words); i++ {
		m := len(words[i])
		word := words[i]
		table := R.Build(word)
		hashes := []uint{R.Query(table, 0, m)} // 不交换
		for j := 0; j < m-1; j++ {             // 交换
			newWord := string(word[j+1]) + string(word[j])
			mid := R.Query(R.Build(newWord), 0, 2)
			left := R.Query(table, 0, j)
			right := R.Query(table, j+2, m)
			left = R.Combine(left, mid, 2)
			left = R.Combine(left, right, m-j-2)
			hashes = append(hashes, left)
		}

		for _, h := range hashes {
			if _, ok := visited[h]; ok {
				res[i] = true
				break
			}
		}

		visited[hashes[0]] = struct{}{}
	}

	return res
}

func yuki2102() {
	// https://yukicoder.me/problems/no/2102
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	res := ConditionalReflection(words)
	for i := 0; i < n; i++ {
		if res[i] {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

type S = string

func GetHash(s S, base uint, mod uint) uint {
	if len(s) == 0 {
		return 0
	}
	res := uint(0)
	for i := 0; i < len(s); i++ {
		res = (res*base + uint(s[i])) % mod
	}
	return res
}

type RollingHash struct {
	base  uint
	mod   uint
	power []uint
}

// eg: NewRollingHash(131, 1e9+7)
// mod: 999999937/999999929/999999893/999999797/999999761/999999757/999999751/999999739
func NewRollingHash(base uint, mod uint) *RollingHash {
	return &RollingHash{
		base:  base,
		mod:   mod,
		power: []uint{1},
	}
}

func (r *RollingHash) Build(s S) (hashTable []uint) {
	sz := len(s)
	hashTable = make([]uint, sz+1)
	for i := 0; i < sz; i++ {
		hashTable[i+1] = (hashTable[i]*r.base + uint(s[i])) % r.mod
	}
	return hashTable
}

func (r *RollingHash) Query(sTable []uint, start, end int) uint {
	r.expand(end - start)
	return (r.mod + sTable[end] - sTable[start]*r.power[end-start]%r.mod) % r.mod
}

func (r *RollingHash) Combine(h1, h2 uint, h2len int) uint {
	r.expand(h2len)
	return (h1*r.power[h2len] + h2) % r.mod
}

func (r *RollingHash) AddChar(hash uint, c byte) uint {
	return (hash*r.base + uint(c)) % r.mod
}

// 两个字符串的最长公共前缀长度.
func (r *RollingHash) LCP(sTable []uint, start1, end1 int, tTable []uint, start2, end2 int) int {
	len1 := end1 - start1
	len2 := end2 - start2
	len := min(len1, len2)
	low := 0
	high := len + 1
	for high-low > 1 {
		mid := (low + high) / 2
		if r.Query(sTable, start1, start1+mid) == r.Query(tTable, start2, start2+mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

func (r *RollingHash) expand(sz int) {
	if len(r.power) < sz+1 {
		preSz := len(r.power)
		r.power = append(r.power, make([]uint, sz+1-preSz)...)
		for i := preSz - 1; i < sz; i++ {
			r.power[i+1] = (r.power[i] * r.base) % r.mod
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
