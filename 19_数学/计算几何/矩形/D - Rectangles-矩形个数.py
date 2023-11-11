# 给出n点二维坐标，（在第一象限）,求可以组成的矩形个数
# 要求矩形的四条边都平行于坐标轴。
# n<=2000

# !枚举对角线
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    pSet = set(points)

    res = 0
    for i in range(n):
        x1, y1 = points[i]
        for j in range(i + 1, n):
            x2, y2 = points[j]
            if x1 != x2 and y1 != y2:
                p3 = (x2, y1)
                p4 = (x1, y2)
                if p3 in pSet and p4 in pSet:
                    res += 1

    print(res // 2)
