import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 非負整数を要素とする H 行 W 列の行列 A が与えられます。 1≤i≤H かつ 1≤j≤W を満たす整数の組 (i,j) について、 A の i 行目 j 列目の要素を A
# i,j
# ​
#   で表します。

# A に対して以下の手順を行います。

# まず、A の要素のうち 0 であるものそれぞれを、任意の正の整数で置き換える（ 0 である要素が複数ある場合、それぞれを異なる正の整数で置き換えることもできます）。
# その後、「下記の 2 つの操作のどちらかを行うこと」を好きな回数（ 0 回でも良い）だけ行う。

# 1≤i<j≤H を満たす整数の組 (i,j) を選び、A の i 行目と j 行目を入れ替える。
# 1≤i<j≤W を満たす整数の組 (i,j) を選び、A の i 列目と j 列目を入れ替える。
# A が次の条件を満たすようにすることができるかどうかを判定してください。

# A
# 1,1
# ​
#  ≤A
# 1,2
# ​
#  ≤⋯≤A
# 1,W
# ​
#  ≤A
# 2,1
# ​
#  ≤A
# 2,2
# ​
#  ≤⋯≤A
# 2,W
# ​
#  ≤A
# 3,1
# ​
#  ≤⋯≤A
# H,1
# ​
#  ≤A
# H,2
# ​
#  ≤⋯≤A
# H,W
# ​


# 言い換えると、1≤i,i
# ′
#  ≤H および 1≤j,j
# ′
#  ≤W を満たす任意の 2 つの整数の組 (i,j) と (i
# ′
#  ,j
# ′
#  ) について、下記の 2 つの条件がともに成り立つ。

# i<i
# ′
#   ならば A
# i,j
# ​
#  ≤A
# i
# ′
#  ,j
# ′

# ​

# 「 i=i
# ′
#   かつ j<j
# ′
#   」ならば A
# i,j
# ​
#  ≤A
# i
# ′
#  ,j
# ′

# ​

if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(ROW)]

    # !二维拓扑?
    # 确定出必须的
