package main

func main() {

}

/* 后缀数组

两个字符串
	长度不小于 k 的公共子串的个数 http://poj.org/problem?id=3415
		单调栈
	最短公共唯一子串 https://codeforces.com/contest/427/problem/D
		唯一性可以用 height[i] 与前后相邻值的大小来判定
	公共回文子串 http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2292
		todo
	todo http://poj.org/problem?id=3729
多个字符串
	多串最长公共子串 SPOJ LCS2 https://www.luogu.com.cn/problem/SP1812 https://loj.ac/p/171 LC1923 https://leetcode.cn/problems/longest-common-subpath/ http://poj.org/problem?id=3450
		拼接，二分答案，对 height 分组，判定组内元素对应不同字符串的个数等于字符串个数
	不小于 k 个字符串中的最长子串 http://poj.org/problem?id=3294
		拼接，二分答案，对 height 分组，判定组内元素对应不同字符串的个数不小于 k
	在每个字符串中至少出现两次且不重叠的最长子串 https://www.luogu.com.cn/problem/SP220
		拼接，二分答案，对 height 分组，判定组内元素在每个字符串中至少出现两次且 sa 的最大最小之差不小于二分值（用于判定是否重叠）
	出现或反转后出现在每个字符串中的最长子串 http://poj.org/problem?id=1226
		拼接反转后的串 s[i]+="#"+reverse(s)，拼接所有串，二分答案，对 height 分组，判定组内元素在每个字符串或其反转串中出现
	acSearch（https://www.luogu.com.cn/problem/P3796）的后缀数组做法
		拼接所有串（模式+文本，# 隔开），对每个模式 p 找其左右范围，满足该范围内 height[i] >= len(p)，这可以用 ST+二分或线段树二分求出，然后统计区间内的属于文本串的后缀
逆向
	todo 根据 sa 反推有多少个能生成 sa 的字符串 https://codeforces.com/problemset/problem/1526/E

 https://www.luogu.com.cn/problem/P5546





*/
