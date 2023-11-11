# !判断无向图是否同构(无向图的同构)
# 枚举排列判断是否能建立映射即可

from collections import defaultdict
from itertools import combinations, permutations
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    adjMap1 = defaultdict(set)
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap1[u].add(v)
        adjMap1[v].add(u)

    adjMap2 = defaultdict(set)
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap2[u].add(v)
        adjMap2[v].add(u)

    for perm in permutations(range(n)):
        flag = True
        for i, j in combinations(range(n), 2):
            if (j in adjMap1[i]) ^ (perm[j] in adjMap2[perm[i]]):
                flag = False
                break

        if flag:
            print("Yes")
            exit(0)

    print("No")
