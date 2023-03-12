# https://www.acwing.com/problem/content/364/

# 给定 n 个区间 [ai,bi] 和 n 个整数 ci。
# !你需要构造一个整数集合 Z，使得 ∀i∈[1,n]，Z 中满足 ai≤x≤bi 的整数 x 不少于 ci 个。
# !求这样的整数集合 Z 最少包含多少个数。
# 1≤n≤50000,
# !0≤ai,bi≤50000

# !前缀和preSum[i]表示[0,i]选择了多少个数 题目要求preSum[50000]的最小值
# 最小值=>求最长路
# !即求在约束下0到50000的最长路
# 所有的限制要找全：
# !1. Si >= Si-1
# !2. Si - Si-1 <= 1
# !3. Sb - Sa-1 >= c

from 差分约束 import DualShortestPath


if __name__ == "__main__":
    n = int(input())
    MAX = 50000
    D = DualShortestPath(MAX + 10, min=True)
    for i in range(n):
        a, b, c = map(int, input().split())
        a, b = a + 1, b + 1  # !从1开始
        D.addEdge(a - 1, b, -c)
    for i in range(1, MAX + 5):  # 多加一点
        D.addEdge(i - 1, i, 0)
        D.addEdge(i, i - 1, 1)

    res, ok = D.run()
    print(res[MAX + 1])
