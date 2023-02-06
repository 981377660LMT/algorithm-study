import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 人の人があるコンテストに参加し、
# i 位の人のハンドルネームは
# S
# i
# ​
#   でした。
# 上位
# K 人のハンドルネームを辞書順に出力してください。
if __name__ == "__main__":
    n, k = map(int, input().split())
    s = [input() for _ in range(n)]
    s = s[:k]
    s.sort()
    for i in range(k):
        print(s[i])
