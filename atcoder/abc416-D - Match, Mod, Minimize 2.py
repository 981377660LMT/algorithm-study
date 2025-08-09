# D - Match, Mod, Minimize 2
# https://atcoder.jp/contests/abc416/tasks/abc416_d
# 给定长度为 N 的非负整数序列 A 和 B，以及正整数 M。可以任意重排 A 的元素，求在这种重排下
# ∑_{i=1..N} (A_i + B_i) mod M 的最小可能值。
#
#
# !0<=Ai,Bi<M
#
#  ∑_{i=1..N} (A_i + B_i) - C*M 最小化
# !配对数尽可能多的>=M.


def solve():
    N, M = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))

    A.sort(reverse=True)
    B.sort()
    c, ptr = 0, 0
    for a in A:
        while ptr < N and B[ptr] + a < M:
            ptr += 1
        if ptr >= N:
            break
        c += 1
        ptr += 1

    res = sum(A) + sum(B) - c * M
    print(res)


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        solve()
