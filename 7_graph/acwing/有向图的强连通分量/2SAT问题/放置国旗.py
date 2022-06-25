# 1-N号旗设置位置
# 第i号旗可以设置在xi位置或者yi位置
# 任意两面旗距离需要大于D
# 是否可以设置旗子

# 1≤N≤1000
# D,Xi,Yi<=1e9

from TwoSAT import TwoSAT

# 放置国旗
n, d = map(int, input().split())
xy = [tuple(map(int, input().split())) for _ in range(n)]


# 构造条件:i选择yi 或者 j选择yj
ts = TwoSAT(n)
for i in range(n):
    x1, y1 = xy[i]
    for j in range(i + 1, n):
        x2, y2 = xy[j]
        if abs(x1 - x2) < d:
            ts.addEdge(i, True, j, True)
        if abs(x1 - y2) < d:
            ts.addEdge(i, True, j, False)
        if abs(y1 - x2) < d:
            ts.addEdge(i, False, j, True)
        if abs(y1 - y2) < d:
            ts.addEdge(i, False, j, False)


ts.buildGraph()
if ts.check():
    print("Yes")
    res = ts.work()
    res = [xy[i][t] for i, t in enumerate(res)]
    print(*res, sep="\n")
else:
    print("No")
