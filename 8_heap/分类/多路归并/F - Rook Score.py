from collections import defaultdict
from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    # 站在交叉点上

    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    S = set((x, y) for x, y, _ in points)
    rowSum, colSum = defaultdict(int), defaultdict(int)
    for x, y, w in points:
        rowSum[x] += w
        colSum[y] += w

    res = 0
    for r, c, w in points:
        res = max(res, rowSum[r] + colSum[c] - w)

    # 不站在交叉点上
    # !两个数组选数,最大化和
    row = sorted(rowSum.items(), key=lambda x: x[1], reverse=True)
    col = sorted(colSum.items(), key=lambda x: x[1], reverse=True)
    pq = [(-(row[0][1] + col[0][1]), row[0][0], col[0][0], 0, 0)]
    while pq:
        s, x, y, ptr1, ptr2 = heappop(pq)
        s = -s

        if (x, y) not in S:
            res = max(res, s)
            print(res)
            exit(0)
        else:
            if ptr1 + 1 < len(row):
                nextS = s - row[ptr1][1] + row[ptr1 + 1][1]
                heappush(pq, (-nextS, row[ptr1 + 1][0], y, ptr1 + 1, ptr2))
            if ptr2 + 1 < len(col):
                nextS = s - col[ptr2][1] + col[ptr2 + 1][1]
                heappush(pq, (-nextS, x, col[ptr2 + 1][0], ptr1, ptr2 + 1))
    print(res)
