// 部分可持久化KMP/动态KMP/可撤销KMP
// 这里的kmp说的是一个可以查询border的数据结构，修改就是往后添加一个字符
// 历史版本可以查询，最新版本可以修改和查询
// https://zhuanlan.zhihu.com/p/527154301

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF1721E()
}

// Prefix Function Queries
// https://www.luogu.com.cn/problem/CF1721E
// 给定字符串s以及q个串ti，求将s分别与每个ti拼接起来后，最靠右的|ti|个前缀的 border 长度。询问间相互独立。。
//
// 由于 KMP 是基于均摊的，所以显然不能每次询问暴力跑一遍 KMP.
// 考虑优化询问时 KMP 跳 next 的过程：
// 预处理时记录每种状态后面加每种字符的 next，其实就是单串的 AC 自动机，
// 当询问 KMP 时跳到原串部分后，直接返回结果。
func CF1721E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	K := NewKmpUndoable(len(s), func(i int) int { return int(s[i]) })

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var t string
		fmt.Fscan(in, &t)
		K.Snapshot()
		for i := 0; i < len(t); i++ {
			K.Append(curVersion, t[i])
			fmt.Fprint(out, K.Query(curVersion), " ")
		}
		fmt.Fprintln(out)
	}
}

// #include<bits/stdc++.h>
// using namespace std;
// /**
//  * @brief 计算st从slen开始部分到tlen的前缀函数，前提是<slen的部分已经求出
//  *
//  * @param st 字符串
//  * @param slen 原始长度
//  * @param tlen 目标长度，tlen > slen
//  * @param pi 前缀函数
//  */
// void calPi(int* st, int slen, int tlen, int* pi) {
//   int j = slen ? pi[slen - 1] : 0;
//   for (int i = max(slen, 1); i < tlen; i++) {
//     while (j && st[j] != st[i]) {
//       if (pi[j - 1] <= j / 2 || st[i] == st[pi[j - 1]])
//         j = pi[j - 1];
//       else {
//         j = pi[j % (j - pi[j - 1]) + (j - pi[j - 1]) - 1];
//       }
//     }

//     if (j || st[j] == st[i]) pi[i] = ++j;
//     else pi[i] = 0;
//   }
// }
// const int N = 2e5 + 5;
// int s[N], pi[N];
// void solve() {
//   int n;
//   cin >> n;
//   int len = 0;
//   pi[0] = 0;
//   vector<int> ans;
//   int m = n;
//   while (n--) {
//     char op;
//     cin >> op;
//     if (op == '+') {
//       cin >> s[len];
//       len++;
//       calPi(s, len - 1, len, pi);
//     }
//     else {
//       --len;
//     }
//     cout << pi[max(len - 1, 0)] << "\n";
//   }
// }

type KmpUndoable struct {
	State int // 状态即为当前字符长度
	Nums  []int
	Fail  []int
}

// Kmp 初始状态.
func NewKmpUndoable(n int, f func(i int) int) *KmpUndoable {
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = f(i)
	}
	fail := make([]int, n)
	j := 0
	for i := 1; i < n; i++ {
		vi := nums[i]
		for j > 0 && vi != nums[j] {
			j = fail[j-1]
		}
		if vi == nums[j] {
			j++
		}
		fail[i] = j
	}
	return &KmpUndoable{State: n, Nums: nums, Fail: fail}
}

// 在当前字符串末尾追加一个字符.
func (kmp *KmpUndoable) Append(c int) {
	if kmp.State < len(kmp.Nums) {
		kmp.Nums[kmp.State] = c
	} else {
		kmp.Nums = append(kmp.Nums, c)
	}
	kmp.State++
	kmp.updateFail(kmp.State-1, kmp.State)
}

// 查询当前字符串的border长度.
func (kmp *KmpUndoable) Query() int {
	return kmp.Fail[max(kmp.State-1, 0)]
}

func (kmp *KmpUndoable) GetState() int {
	return kmp.State
}

func (kmp *KmpUndoable) Undo() bool {
	if kmp.State == 0 {
		return false
	}
	kmp.State--
	return true
}

func (kmp *KmpUndoable) Rollback(state int) bool {
	if state < 0 || state > kmp.State {
		return false
	} else {
		kmp.State = state
		return true
	}
}

// 计算字符串从rawLen开始到targetLen的前缀函数，前提是<rawLen的部分已经求出.
func (kmp *KmpUndoable) updateFail(rawLen, targetLen int) {

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
