# 给定二维平面上的n个点(xi,yi)和一个正整数k。
# !请列出所有欧几里得距离小于等于k的点对。
# !1<=n<=2e5,1<=k<=1.5e9。保证答案不超过4e5对。
# 按字典序输出所有点对

# !分桶
# !然后枚举九个相邻桶

from collections import defaultdict
import sys

DIR9 = [(1, 0), (0, 1), (-1, 0), (0, -1), (1, 1), (1, -1), (-1, 1), (-1, -1), (0, 0)]

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, k = map(int, input().split())
points = [tuple(map(int, input().split())) for _ in range(n)]
res = []  # 点对序号 (i,j) i<j

bucket = defaultdict(list)
for i, (r, c) in enumerate(points):
    bucket[(r // k, c // k)].append(i)

for i, (r1, c1) in enumerate(points):
    rid1, cid1 = r1 // k, c1 // k
    for dr, dc in DIR9:
        rid2, cid2 = rid1 + dr, cid1 + dc
        for j in bucket[(rid2, cid2)]:
            if j >= i:
                break
            r2, c2 = points[j]
            dist = (r1 - r2) * (r1 - r2) + (c1 - c2) * (c1 - c2)
            if dist <= k * k:
                res.append((j, i))


res.sort()
print(len(res))
for pair in res:
    print(pair[0] + 1, pair[1] + 1)


# !时间复杂度证明:O(n+m) m为输出答案的个数,n为点数
