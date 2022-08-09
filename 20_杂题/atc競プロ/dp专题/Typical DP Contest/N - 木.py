# https://algo-logic.info/typical-dp-contest-n/
# 用线连接出一棵树所有的边
# 要求连接过程中所有边始终是连接着的
# 有多少种连接方法
# n<=1000

# !蚂蚁构建房间的不同顺序
# 树形dp返回(子节点个数、连接方案)
import sys
from typing import Tuple


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

fac = [1]
ifac = [1]
for i in range(1, int(1e4)):
    fac.append((fac[-1] * i) % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


n = int(input())
adjList = [[] for _ in range(n)]
for _ in range(n - 1):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    adjList[u].append(v)
    adjList[v].append(u)


def dfs(cur: int, pre: int) -> Tuple[int, int]:
    nodeCount, res = 0, 1
    for next in adjList[cur]:
        if next == pre:
            continue
        subCount, subRes = dfs(next, cur)
        nodeCount += subCount
        res = (res * subRes % MOD) * C(nodeCount, subCount) % MOD
    return (nodeCount + 1), res


res = 0
for root in range(n):  # 每个结点作为根
    res += dfs(root, -1)[1]
    res %= MOD
print((res * ifac[2]) % MOD)  # !注意不能用//2 要用取模意义下的除法
