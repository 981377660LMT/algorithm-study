# https://www.luogu.com.cn/problem/P2662
# 给一堆系数，求最大的不能被线性表出的数。
# n<100,m<3000
# !从起点开始有一些点无法到达(即有一整个剩余系都不能被表出)


# 设 base是我们规定的模数,那么同余最短路求出的距离dist[i]
# 就是能拼凑出的以i为代表元的剩余类中的最小数字
# !那么dist[i]-base就是不能拼凑出的在该剩余类中的最大数字

from ModShortestPath import modShortestPath

INF = int(1e18)
n, cut = map(int, input().split())  # 木料的种类和每根木料削去的最大值
sticks = list(map(int, input().split()))  # 第i根木料的原始长度

coeffs = []
for s in sticks:
    for c in range(min(cut, s) + 1):
        coeffs.append(s - c)

base, dist = modShortestPath(coeffs)
if any(v == INF for v in dist):  # 这个剩余类不能被表出
    print(-1)
    exit(0)

print(max(dist) - base)
