import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 数直線上を N 匹の魚が泳いでいます。

# 魚 i の重さは W
# i
# ​
#   であり、時刻 0 に座標 X
# i
# ​
#   にいて、正方向に速さ V
# i
# ​
#   で移動しています。

# 高橋君は 0 以上の実数 t を自由に選び、時刻 t に一度だけ以下の行動を行います。
# 行動：実数 x を自由に選ぶ。現在の座標が x 以上 x+A 以下である魚を全て捕まえる。

# 捕まえることができる魚の重さの合計の最大値を求めてください。

# 先选左端点，然后根据左端点确定时间
# !枚举左端点的鱼，现在其他鱼位于网内的时间区间就可以算出来了，然后就是区间覆盖问题了
if __name__ == "__main__":
    n, A = map(int, input().split())
    fish = []
    for _ in range(n):
        weight, pos, speed = map(int, input().split())
        fish.append((weight, pos, speed))

    # dp[i][j] 表示
