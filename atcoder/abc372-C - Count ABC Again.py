# abc372-C - Count ABC Again
# https://atcoder.jp/contests/abc372/tasks/abc372_c
# 给定一个长度为N的字符串S，字符串S只包含A、B、C三种字符。
# 有Q次操作，每次操作给定一个位置i和一个字符c，将S的第i个字符替换为c。
# 每次操作后，统计S中有多少个连续的ABC子串。
# 1<=N<=10^5,1<=Q<=10^5
#
# !关键是提取出update函数.

import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")


if __name__ == "__main__":
    N, Q = map(int, input().split())
    S = input()

    res = 0
    sb = list(S)
    valid = [False] * N

    def check(i: int) -> bool:
        if not (0 <= i < N - 2):
            return False
        return sb[i] == "A" and sb[i + 1] == "B" and sb[i + 2] == "C"

    def update(i: int, c: str) -> None:
        if i < 0 or i >= N:
            return
        global res
        sb[i] = c
        for j in range(max(0, i - 2), min(N - 2, i + 1)):
            res -= valid[j]
            valid[j] = check(j)
            res += valid[j]

    for i in range(N):
        update(i, sb[i])

    for _ in range(Q):
        i, c = input().split()
        i = int(i) - 1
        update(i, c)
        print(res)
