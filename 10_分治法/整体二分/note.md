整体二分 Parallel Binary Search
https://oi-wiki.org/misc/parallel-binsearch/
https://codeforces.com/blog/entry/45578
todo 整体二分解决静态区间第 k 小的优化 https://www.luogu.com.cn/blog/2-6-5-3-5/zheng-ti-er-fen-xie-jue-jing-tai-ou-jian-di-k-xiao-di-you-hua
模板题 https://www.luogu.com.cn/problem/P3527
todo https://atcoder.jp/contests/agc002/tasks/agc002_d
https://www.hackerrank.com/contests/hourrank-23/challenges/selective-additions/problem
https://www.codechef.com/problems/MCO16504

---

将询问离线，对所有询问同时进行二分，在二分的过程中逐渐将询问分为 “答案
在 [l,mid]中的” 和 “答案在 [mid+1,r]中的” 两类，直到区间长度变为 1

时间复杂度`O(n+q)log(update)log(max)`
