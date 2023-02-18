from collections import defaultdict
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 1 から
# N の番号がついた、表裏が区別できるコインが
# N 枚あります。コインの表裏は長さ
# N の文字列
# S で表され、
# S の
# i 番目の文字が 1 のときコイン
# i は表を向いており、0 のときコイン
# i は裏を向いています。

# あなたは、以下の操作を
# 0 回以上好きな回数繰り返すことができます。

# 1≤i<j≤N かつ
# j−i≥2 を満たす整数組
# (i,j) を選ぶ。コイン
# i とコイン
# j を裏返す。
# 操作によって
# N 枚のコイン全てを裏向きにできるか判定し、可能な場合必要な操作の回数の最小値を求めてください。

# T 個のテストケースが与えられるので、それぞれについて答えてください。


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        n = int(input())
        s = input()
        ones = s.count("1")
        if (ones & 1) == 1:
            print(-1)
            continue
        if ones == 2 and "11" in s:
            print(-1)
            continue
        print(ones // 2)
