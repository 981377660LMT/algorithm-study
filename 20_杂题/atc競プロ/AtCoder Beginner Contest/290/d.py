from math import gcd
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 0 から
# N−1 までの番号がつけられた
# N 個のマスが並んでいます。 今から、すぬけくんが以下の手順に従って全てのマスに印をつけていきます。

# マス
# 0 に印をつける。
# 次の i - iii の手順を
# N−1 回繰り返す。
# 最後に印をつけたマスの番号を
# A としたとき、変数
# x を
# (A+D)modN で初期化する。
# マス
# x に印が付いている限り、
# x を
# (x+1)modN に更新することを繰り返す。
# マス
# x に印をつける。
# すぬけくんが
# K 番目に印をつけるマスの番号を求めてください。

# T 個のテストケースが与えられるので、それぞれについて答えを求めてください。
if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        n, d, k = map(int, input().split())
        if k == 1:
            print(0)
            continue
        if gcd(n, d) == 1:
            print(d * (k - 1) % n)
            continue

        # 重合的次数
        # 多少个lcm重合了
        lcm = n * d // gcd(n, d)
        # 重合的次数
        count = (k - 1) * d // lcm
        print((count + d * (k - 1)) % n)
# 5 8 1
# 0
# 5 8 2
# 3
# 5 8 3
# 1
# 5 8 4
# 4
# 5 8 5
# 2


# 9
# 4 2 1
# 4 2 2
# 4 2 3
# 4 2 4
# 5 8 1
# 5 8 2
# 5 8 3
# 5 8 4
# 5 8 5
# 0 2 1 3 0 3 1 4 2
