// API:
//  RollingHash(base) : 以base为底的哈希函数.
//  Build(s) : O(n)返回s的哈希值表.
//  Query(sTable, l, r) : 返回s[l, r)的哈希值.
//  Combine(h1, h2, h2len) : 哈希值h1和h2的拼接, h2的长度为h2len. (线段树的op操作)
//  AddChar(h,c) : 哈希值h和字符c拼接形成的哈希值.
//  Lcp(aTable, l1, r1, bTable, l2, r2) : O(logn)返回a[l1, r1)和b[l2, r2)的最长公共前缀.
//
// !此外，常用技巧为：
// 1. 对哈希值排序，然后使用二分查找是否存在相同的哈希值。
// 2. 按照字符串长度对哈希值分组，最多存在根号种不同的长度分组.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	P3501()
	// CF514C()
	// CF1056E()
	// SP220()
}

// P3501 [POI2010] ANT-Antisymmetry (!对于一个0/1序列，求出其中异或意义下回文的子串数量。)
// https://www.luogu.com.cn/problem/P3501
// 对于一个01字符串，如果将这个字符串0和1取反后，再将整个串反过来和原串一样，就称作“反对称”字符串。
// 比如00001111和010101就是反对称的，1001就不是。
// 现在给出一个长度为N的01字符串，求它有多少个子串是反对称的。
// eg:
// 11001011
// 7个反对称子串分别是：01（出现两次），10（出现两次），0101，1100和001011
//
// 反对称一定是偶数长度的回文串.
// 对每个位置i，哈希+二分求出以i为中心的最长偶数长度回文串的长度，然后累加即可.
func P3501() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	var s01 string
	fmt.Fscan(in, &s01)

	getRev := func(s string) string {
		sb := []byte(s)
		for i := 0; i < len(sb); i++ {
			if sb[i] == '0' {
				sb[i] = '1'
			} else {
				sb[i] = '0'
			}
		}
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			sb[i], sb[j] = sb[j], sb[i]
		}
		return string(sb)
	}
	revS01 := getRev(s01)

	R := NewRollingHash(131, 999999751)
	table := R.Build(s01)
	revTable := R.Build(revS01)

	// [start,end)是否为反对称回文串.
	check := func(start, end int32) bool {
		if start < 0 || end > n {
			return false
		}
		return R.Query(table, int(start), int(end)) == R.Query(revTable, int(n-end), int(n-start))
	}

	//以center为中心的最长偶数长度回文串的长度.
	getMaxEvenRadius := func(center int32) int32 {
		left, right := int32(1), n
		for left <= right {
			mid := (left + right) / 2
			if check(center-mid+1, center+mid+1) {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		return right
	}

	res := 0
	for i := int32(0); i < n; i++ {
		res += int(getMaxEvenRadius(int32(i)))
	}
	fmt.Fprintln(out, res)
}

// Watto and Mechanism
// https://www.luogu.com.cn/problem/CF514C
// 给出 n 个已知字符串，q 次询问，每次询问给出一个字符串，
// 问上面 n 个字符串中是否有一个字符串满足恰好有一个字母不同于询问的字符串。
// n,m<=3e5, 所有串总长度<=6e5.
// 所有字符串的字符属于{'a','b','c'}
//
// 哈希+长度分组
func CF514C() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	R := NewRollingHash(37, 2102001800968)
	hashGroup := make(map[int][]uint) // 按照长度分组
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		table := R.Build(s)
		hash := R.Query(table, 0, len(s))
		hashGroup[len(s)] = append(hashGroup[len(s)], hash)
	}
	for _, group := range hashGroup {
		sort.Slice(group, func(i, j int) bool {
			return group[i] < group[j]
		})
	}

	charHash := make([]uint, 3)
	for i := 0; i < 3; i++ {
		s := string(rune('a' + i))
		table := R.Build(s)
		charHash[i] = R.Query(table, 0, 1)
	}

	// 是否存在一个字符不同
	query := func(s string) bool {
		m := len(s)
		table := R.Build(s)
		for i := 0; i < m; i++ {
			preHash := R.Query(table, 0, i)
			sufHash := R.Query(table, i+1, m)
			for j := 0; j < 3; j++ {
				if j == int(s[i]-'a') {
					continue
				}
				newHash := R.Combine(preHash, charHash[j], 1)
				newHash = R.Combine(newHash, sufHash, m-(i+1))
				if BinarySearch(hashGroup[m], newHash, true) != -1 {
					return true
				}
			}
		}
		return false
	}

	for i := 0; i < q; i++ {
		var s string
		fmt.Fscan(in, &s)

		ok := query(s)
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
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

	R1 := NewRollingHash(131, 999999929)
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

// PHRASES - Relevant Phrases of Annihilation
// https://www.luogu.com.cn/problem/SP220
// 给定n个字符串，求在每个字符串中至少出现两次且不重叠的最长子串长度.
//
// 解法1：广义SAM+维护endPos等价类最小End、最大End.但是复杂度与字符个数有关,不够好.
// !解法2：二分长度+哈希分组
func SP220() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	H := NewRollingHash(131, 999999751)

	solve := func(words []string) int32 {
		n := int32(len(words))
		tables := make([][]uint, n)
		for i := int32(0); i < n; i++ {
			tables[i] = H.Build(words[i])
		}

		check := func(mid int32) bool {
			counter := make([]map[uint]int32, n)
			for i := int32(0); i < n; i++ {
				counter[i] = make(map[uint]int32)
			}
			visited := make(map[uint]struct{})
			for i, table := range tables {
				m := int32(len(words[i]))
				last := make(map[uint]int32) // 记录每个哈希值上一次出现的位置
				for j := int32(0); j+mid <= m; j++ {
					hash_ := H.Query(table, int(j), int(j+mid))
					v, ok := last[hash_]
					if !ok || v <= j { // 不重叠出现
						last[hash_] = j + mid
						counter[i][hash_]++
						visited[hash_] = struct{}{}
					}
				}
			}
			for h := range visited {
				ok := true
				for _, c := range counter {
					if c[h] < 2 {
						ok = false
						break
					}
				}
				if ok {
					return true
				}
			}
			return false
		}

		left, right := int32(1), int32(len(words[0]))
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return right
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var n int32
		fmt.Fscan(in, &n)
		words := make([]string, n)
		for j := int32(0); j < n; j++ {
			fmt.Fscan(in, &words[j])
		}
		fmt.Fprintln(out, solve(words))
	}
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

// eg:
// NewRollingHash(37, 2102001800968)
// NewRollingHash(131, 999999751)
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

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

// 查询非递减数组中目标值的第一个或最后一个位置.
// arr: 非递减数组
// target: 目标值
// findFirst: 是否查询第一个位置. 默认为 true.
// 返回目标值的第一个或最后一个位置. 如果目标值不存在, 返回-1.
func BinarySearch[T Int](arr []T, target T, findFirst bool) int {
	if len(arr) == 0 || arr[0] > target || arr[len(arr)-1] < target {
		return -1
	}
	if findFirst {
		left, right := 0, len(arr)-1
		for left <= right {
			mid := left + (right-left)>>1
			if arr[mid] < target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		if left < len(arr) && arr[left] == target {
			return left
		}
		return -1
	} else {
		left, right := 0, len(arr)-1
		for left <= right {
			mid := left + (right-left)>>1
			if arr[mid] <= target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		if left > 0 && arr[left-1] == target {
			return left - 1
		}
		return -1
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
