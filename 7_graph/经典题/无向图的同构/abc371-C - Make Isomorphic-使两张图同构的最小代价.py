import sys
from itertools import permutations


input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)


# C - Make Isomorphic
# https://atcoder.jp/contests/abc371/tasks/abc371_c
# 使两张图同构的最小代价.
# N<=8.
if __name__ == "__main__":
    N = int(input())
    G1 = [[False] * N for _ in range(N)]
    M1 = int(input())
    for _ in range(M1):
        a, b = map(int, input().split())
        a -= 1
        b -= 1
        G1[a][b] = G1[b][a] = True
    G2 = [[False] * N for _ in range(N)]
    M2 = int(input())
    for _ in range(M2):
        a, b = map(int, input().split())
        a -= 1
        b -= 1
        G2[a][b] = G2[b][a] = True

    C = [[0] * N for _ in range(N)]
    for i in range(N):
        cur = list(map(int, input().split()))
        for j in range(i + 1, N):
            C[i][j] = C[j][i] = cur[j - (i + 1)]

    res = INF
    for perm in permutations(range(N)):
        cost = 0
        for i in range(N):
            for j in range(i + 1, N):
                if G1[i][j] != G2[perm[i]][perm[j]]:
                    cost += C[perm[i]][perm[j]]
        res = min(res, cost)

    print(res)
