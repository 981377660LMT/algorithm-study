import sys

from sortedcontainers.sortedlist import SortedList

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 英小文字からなる長さ
# N の文字列
# S と
# Q 個のクエリが与えられます。クエリを順に処理してください。

# クエリは以下の
# 2 種類です。

# 1 x c ：
# S の
# x 文字目を文字
# c に置き換える
# 2 l r ：
# S を文字の昇順に並び替えて得られる文字列を
# T とする。
# S の
# l 文字目から
# r 文字目までからなる文字列が
# T の部分文字列であるとき Yes、部分文字列でないとき No を出力する
if __name__ == "__main__":
    n = int(input())
    s = input()
    q = int(input())

    sl = SortedList(s)
    arr = list(s)
    for _ in range(q):
        op, *rest = input().split()
        op = int(op)
        if op == 1:
            a, b = rest
            a = int(a) - 1
            sl.remove(arr[a])
            sl.add(b)
            arr[a] = b
        else:
            left, right = map(int, rest)
            left -= 1
