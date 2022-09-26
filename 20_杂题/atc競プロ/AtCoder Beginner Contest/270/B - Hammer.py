import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 数直線の原点に高橋君がいます。高橋君は座標 X にあるゴールに移動しようとしています。
# 座標 Y には壁があり、最初、高橋君は壁を超えて移動することができません。
# 座標 Z にあるハンマーを拾った後でなら、壁を破壊して通過できるようになります。
# 高橋君がゴールに到達することが可能か判定し、可能であれば移動距離の最小値を求めてください。

if __name__ == "__main__":
    x, y, z = map(int, input().split())

    # 不捡锤子
    if not ((0 <= y <= x) or (x <= y <= 0)):
        print(abs(x))
        exit(0)

    # 捡锤子
    if not ((0 <= y <= z) or (z <= y <= 0)):
        print(abs(x - z) + abs(z))
        exit(0)

    # 捡不了锤子
    print(-1)
