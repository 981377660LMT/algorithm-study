# n,m,q<=2e5
# !每个查询 [xi,yi] 输出xi结点的颜色 并把xi结点和邻居全部变为yi颜色
# 最开始所有结点颜色为1
# !注意到TLE的数据: 中心一个点 周围连着很多点 (星图)

# !分块算法
# !邻居很多的点(大顶点，数量少，利用数量少的特性暴力更新):暴力更新这些大顶点的颜色(即保持大顶点的查询和更新都是及时的)
# !邻居很少的点(小顶点，数量多，利用邻居少的特性查邻居):用一个lasts数组记录每个节点最后一个询问是哪个时间 小顶点最后的颜色就是max(last[邻居])
# !小顶点无法及时更新颜色 查询颜色只能由周围的记录来查询

# !分块界限 sqrt(2*M)
# !时间复杂度O(Q*sqrt(2*M))
# https://atcoder.jp/contests/typical90/submissions/24052294

from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n, m = map(int, input().split())
adjMap = defaultdict(set)
for _ in range(m):
    x, y = map(int, input().split())
    x, y = x - 1, y - 1
    adjMap[x].add(y)
    adjMap[y].add(x)

# !1.处理大小顶点
SQRT = int(2 * m ** 0.5)
big = set()
for i in range(n):
    if len(adjMap[i]) >= SQRT:
        big.add(i)
small = set(range(n)) - big


# !2.每个点连接的大、小顶点
group = defaultdict(lambda: [set(), set()])
for i in range(n):
    for next in adjMap[i]:
        if next in big:
            group[i][0].add(next)
        else:
            group[i][1].add(next)


q = int(input())
colors = [1] * n
lasts = [-1] * n  # 每个节点最后的查询位置
history = []  # 历史颜色更新值
for i in range(q):
    node, newColor = map(int, input().split())
    node -= 1
    res = 1

    # !大顶点实时查询
    if node in big:
        res = colors[node]
        colors[node] = newColor
    else:
        # !小顶点邻居查询
        preI = lasts[node]
        for next in adjMap[node]:
            preI = max(preI, lasts[next])
        if preI == -1:
            res = 1
        else:
            res = history[preI]

    # !大顶点实时暴力更新
    for nextBig in group[node][0]:
        colors[nextBig] = newColor

    lasts[node] = i
    history.append(newColor)
    print(res)

