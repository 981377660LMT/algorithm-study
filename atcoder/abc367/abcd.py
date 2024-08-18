from collections import defaultdict
from itertools import product
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N, M = map(int, input().split())
    A = list(map(int, input().split()))
    A2 = A[:] + A[:]
    N2 = 2 * N
    res, curSum, preSum = 0, 0, defaultdict(int)
    visited = []
    for right in range(N2 - 1):
        curSum += A2[right]
        mod = curSum % M
        preSum[mod] += 1
        visited.append(mod)
        if right >= N:
            preSum[visited[right - N]] -= 1
        if right >= N - 1:
            res += preSum[mod]
    print(res - N)
