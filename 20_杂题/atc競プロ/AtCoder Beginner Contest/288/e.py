import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# お店に
# N 個の商品が並んでおり、それらは商品
# 1 、商品
# 2 、
# … 、商品
# N と番号づけられています。
# i=1,2,…,N について、商品
# i の定価は
# A
# i
# ​
#   円です。また、各商品の在庫は
# 1 つです。

# 高橋君は、商品
# X
# 1
# ​
#   、商品
# X
# 2
# ​
#   、
# … 、商品
# X
# M
# ​
#   の
# M 個の商品が欲しいです。

# 高橋君は、欲しい商品をすべて手に入れるまで、下記の行動を繰り返します。

# 現在売れ残っている商品の個数を
# r とする。
# 1≤j≤r を満たす整数
# j を選び、現在売れ残っている商品のうち番号が
# j 番目に小さい商品を、その定価に
# C
# j
# ​
#   円だけ加えた金額で購入する。

# 高橋君が欲しい商品をすべて手に入れるまでにかかる合計費用としてあり得る最小値を出力してください。

# なお、高橋君は欲しい商品ではない商品を購入することもできます。
if __name__ == "__main__":
    n, m = map(int, input().split())
    prices = list(map(int, input().split()))
    add = list(map(int, input().split()))
    needs = list(map(int, input().split()))

    # 每个物品在这一轮买不买
    # dp[i][j]表示买前i个物品,第i个物品买的顺序位于第j时的最小代价
