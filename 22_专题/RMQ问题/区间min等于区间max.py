# https://codeforces.com/problemset/problem/689/D

# 给你两个数组 a b，长度均为 n(n<=2e5)，元素范围 [-1e9,1e9]。
# 求所有满足 max(a[l..r]) = min(b[l..r]) 的区间 [l,r] 的个数。


# 输入
# a=[1,2,3,2,1,4]
# b=[6,7,1,2,3,2]
# 输出 2
# 解释 （下标从 0 开始的区间）[3,3] 和 [3,4]，
# 对于 [3,3] 有 max(2)=min(2)，对于 [3,4] 有 max(2,1)=min(2,3)


# https://codeforces.com/contest/689/submission/158934891

# !枚举左端点，区间 max 单调递增，区间 min 单调递减，
# !因此可以二分找 max>=min 和 max>min 的位置，
# 两者相减即为右端点的合法个数。

# 话说这题改成 max^2 = min^2 还能不能做🤔

