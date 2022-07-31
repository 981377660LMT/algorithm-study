# n,m<=2e5 无向无权稀疏图
# 求结点染色的方案数(数据量暗示dp???) mod 998244353
# 使得 红色顶点有`k个` 且 顶点颜色不同的边有`偶数条`

# !考察度数 每个一对红色顶点连接 那么顶点颜色不同的贡献就会少2
# !红节点度数之和 = 顶点颜色不同的边 + 顶点全红的边*2
# !则红节点度数之和为偶数 即选k个结点度数和为偶数


import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


fac = [1]
ifac = [1]
for i in range(1, int(2e5) + 10):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def main() -> None:
    n, m, k = map(int, input().split())
    deg = [0] * n
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        deg[u] += 1
        deg[v] += 1

    odd, even = 0, 0
    for i in range(n):
        if deg[i] % 2 == 1:
            odd += 1
        else:
            even += 1

    res = 0
    for oddCount in range(0, k + 1, 2):
        evenCount = k - oddCount
        res += C(odd, oddCount) * C(even, evenCount)
        res %= MOD
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
