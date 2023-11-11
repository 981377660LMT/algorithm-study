# 每个树结点的值为'a'或'b'
# !在树中删去`若干条`边
# 使得剩下每个联通分量都有'a'和'b'两个值
# 求删边的方案数 (不删边视作删去0条边)

# 树形dp
import sys
from typing import Tuple

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())
values = input().split()
adjList = [[] for _ in range(n)]
for _ in range(n - 1):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    adjList[u].append(v)
    adjList[v].append(u)


def dfs(cur: int, pre: int) -> Tuple[int, int, int]:
    """后序dfs返回(只有a 只有b ab都有)的删边方案数"""
    res1, res2, res3 = int(values[cur] == "a"), int(values[cur] == "b"), 1
    for next in adjList[cur]:
        if next == pre:
            continue
        a, b, c = dfs(next, cur)
        res1 *= a + c  # 不删边 a / 删边 c
        res1 %= MOD
        res2 *= b + c  # 不删边 b / 删边 c
        res2 %= MOD

        # 不删边 a b c / 删边 c
        res3 *= (a + b + c) + c  # !注意这里多算了一些非法状态 最后再减去
        res3 %= MOD
    return res1, res2, (res3 - (res2 if values[cur] == "b" else res1)) % MOD


print(dfs(0, -1)[2])
