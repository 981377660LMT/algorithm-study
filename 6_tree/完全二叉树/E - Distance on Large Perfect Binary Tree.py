"""
题意:有一颗2^n-1个结点的完全二叉树,边权为1,
问你有多少点对(i,j)它们的最短路径为D.注意(i,j)和(j,i)不相同。
2<=n<=1e6,D<=2e6

枚举LCA
- 我们把一个节点视作 01 序列,走左子树相当于在序列末尾添加 0 ,走右子树相当于添加 1 。
那么两个点的距离其实取决于`公共前缀的长度`。
- 相同深度的点对答案的贡献是一样的
这里必须分两种情况计算:
1.两个点在同一子树内,2^i * 2^D = 2^(i+D)
2.两个节点一个在左子树,一个在右子树,不难发现总数怕为2^(D-2)
  假设左子树深度为 i+k ，右子树深度为 i+(D-k) ，这里 1<=k<=D-1 。
  因为深度不能超过 n ，所以 i+k<=n 且 i+(D-k)<=n
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, D = map(int, input().split())
    res = 0
    for depth in range(n):
        # 1. 两个点在同一子树内
        if depth + D < n:
            res += pow(2, depth + D, MOD)
            res %= MOD

        # 2. 两个节点一个在左子树,一个在右子树
        if D >= 2:
            left = max(1, depth + D - n + 1)
            right = min(D - 1, n - 1 - depth)
            if left <= right:
                res += (right - left + 1) * pow(2, depth + D - 2, MOD)
                res %= MOD
    print(res * 2 % MOD)

# TODO
