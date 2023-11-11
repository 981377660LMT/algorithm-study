# 小根堆在线维护前k大元素

from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, k = map(int, input().split())
    perm = list(map(int, input().split()))  # 1-n的全排列

    pq = []
    res = []
    for i in range(n):
        heappush(pq, perm[i])
        if len(pq) > k:
            heappop(pq)
        res.append(pq[0])
    print(*res[k - 1 :], sep="\n")
