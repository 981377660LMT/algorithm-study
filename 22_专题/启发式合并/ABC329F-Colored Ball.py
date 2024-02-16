# ABC329F-Colored Ball
# https://atcoder.jp/contests/abc329/tasks/abc329_f
# 有 q次操作，每次将 a 号箱子里的所有球放入 b号箱子，问此时 b 箱子中球的颜色种类数。


import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    n, q = map(int, input().split())
    nums = list(map(int, input().split()))

    sets = [set([v]) for v in nums]
    for _ in range(q):
        a, b = map(int, input().split())
        a, b = a - 1, b - 1
        if len(sets[a]) > len(sets[b]):
            sets[a], sets[b] = sets[b], sets[a]
        sets[b] |= sets[a]
        sets[a] = set()
        print(len(sets[b]))
