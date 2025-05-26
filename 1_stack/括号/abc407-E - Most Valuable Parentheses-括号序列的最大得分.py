# E - Most Valuable Parentheses(括号序列的最大得分)
# https://atcoder.jp/contests/abc407/tasks/abc407_e
# 给定长度为 2N 的整数数组 A，
# 你要构造一个长度为 2N 的合法括号序列 s（即有效括号），
# 定义其得分为：
# 把所有 s[i] = ')' 的位置的 A[i] 变成 0 后，A 的元素和。
# 问所有合法括号序列中最大得分是多少。
#
# !对于前缀长度 i，最多可以跳过（置为‘)’）floor(i/2) 个位置，其余都必须选为‘(’。

from heapq import heappop, heappush


def solve():

    T = int(input())
    for _ in range(T):
        N = int(input()) * 2
        A = [0] * N
        for i in range(N):
            A[i] = int(input())

        res = 0
        pq = []
        for i in range(N):
            heappush(pq, -A[i])
            if i & 1 == 0:
                res += -heappop(pq)
        print(res)


if __name__ == "__main__":
    solve()
