from collections import Counter
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    pre = dict()
    B = []
    for i in range(N):
        x = A[i]
        if x in pre:
            B.append(pre[x] + 1)
        else:
            B.append(-1)
        pre[x] = i
    print(" ".join(map(str, B)))
