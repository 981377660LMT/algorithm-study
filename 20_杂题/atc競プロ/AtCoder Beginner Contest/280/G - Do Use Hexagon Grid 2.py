import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 以下のような、無限に広い六角形のグリッドがあります。
# 六角形のマスは 2 つの整数 i,j を用いて (i,j) と表されます。
# マス (i,j) は以下の 6 つのマスと辺で隣接しています。

# (i−1,j−1)
# (i−1,j)
# (i,j−1)
# (i,j+1)
# (i+1,j)
# (i+1,j+1)
# 2 つのマス X,Y の距離を、辺で隣接しているマスをたどってマス X からマス Y まで移動するときの、移動回数の最小値と定めます。
# 例えばマス (0,0) とマス (1,1) の距離は 1、マス (2,1) とマス (−1,−1) の距離は 3 です。

# N 個のマス (X
# 1
# ​
#  ,Y
# 1
# ​
#  ),…,(X
# N
# ​
#  ,Y
# N
# ​
#  ) が与えられます。
# この N マスの中から 1 つ以上のマスを選ぶ方法のうち、選んだマスのうちどの 2 マスの距離も D 以下になるようなものは何通りありますか？
# 998244353 で割ったあまりを求めてください。


N = 400
bell = [[0] * N for _ in range(N)]
bell[1][1] = 1
for i in range(2, N):
    bell[i][1] = bell[i - 1][i - 1]
    for j in range(2, i + 1):
        bell[i][j] = bell[i - 1][j - 1] + bell[i][j - 1]
        bell[i][j] %= MOD

# 将 n 个元素划分成若干个非空子集的划分方案数


def cal(n):
    return bell[n][n]


# 一张图的生成树个数即为其基尔霍夫矩阵的行列式。
if __name__ == "__main__":
    n, d = map(int, input().split())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    matrix = [[INF] * n for _ in range(n)]
    for i in range(n):
        x1, y1 = points[i]
        for j in range(n):
            if i == j:
                matrix[i][j] = 0
                continue
            x2, y2 = points[j]
            if (x1 - x2) * (y1 - y2) < 0:
                matrix[i][j] = abs(x1 - x2) + abs(y1 - y2)
            else:
                matrix[i][j] = max(abs(x1 - x2), abs(y1 - y2))

    # 矩阵树定理?
    print(matrix)
    # 枚举距离?
