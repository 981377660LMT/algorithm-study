from itertools import combinations
import re
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# 英小文字および (, ) からなる長さ
# N の文字列
# S が与えられます。
# 以下の操作を可能な限り繰り返したあとの
# S を出力してください。

# S の連続部分文字列であって、最初の文字が ( かつ 最後の文字が ) かつ 最初と最後以外に ( も ) も含まないものを自由に
# 1 つ選び削除する
# なお、操作を可能な限り繰り返したあとの
# S は操作の手順によらず一意に定まることが証明できます。
if __name__ == "__main__":
    n = int(input())
    s = input()

    print(re.sub(r"\(.*\)", "", s))
