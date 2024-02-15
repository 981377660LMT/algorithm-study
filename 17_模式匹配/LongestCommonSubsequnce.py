# 最长公共子序列/最长公共子序列求方案

from typing import Any, List, Sequence, Tuple


def max2(a: int, b: int) -> int:
    return a if a > b else b


def longestCommonSubsequnce(seq1: Sequence["Any"], seq2: Sequence["Any"]) -> int:
    n2 = len(seq2)
    dp = [0] * (n2 + 1)
    for a in seq1:
        for i in range(n2 - 1, -1, -1):
            if a == seq2[i]:
                dp[i + 1] = max2(dp[i + 1], dp[i] + 1)
        for i in range(n2):
            dp[i + 1] = max2(dp[i + 1], dp[i])
    return dp[n2]


def longestCommonSubsequnceRestore(
    seq1: Sequence["Any"], seq2: Sequence["Any"]
) -> List[Tuple[int, int]]:
    n, m = len(seq1), len(seq2)
    dp = [[0] * (m + 1) for _ in range(n + 1)]
    for i in range(n):
        pre = dp[i]
        cur = pre[:]
        dp[i + 1] = cur
        for j in range(m):
            cur[j + 1] = max2(cur[j + 1], cur[j])
            if seq1[i] == seq2[j]:
                cur[j + 1] = max2(cur[j + 1], pre[j] + 1)

    res = []
    ptr1, ptr2 = n, m
    while dp[ptr1][ptr2] > 0:
        if dp[ptr1][ptr2] == dp[ptr1 - 1][ptr2]:
            ptr1 -= 1
        elif dp[ptr1][ptr2] == dp[ptr1][ptr2 - 1]:
            ptr2 -= 1
        else:
            ptr1 -= 1
            ptr2 -= 1
            res.append((ptr1, ptr2))
    return res[::-1]


if __name__ == "__main__":
    # https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=ALDS1_10_C&lang=ja
    T = int(input())
    for i in range(T):
        s1 = input()
        s2 = input()
        res1 = longestCommonSubsequnce(s1, s2)
        res2 = longestCommonSubsequnceRestore(s1, s2)
        for i, j in res2:
            assert s1[i] == s2[j]
        for i in range(len(res2) - 1):
            assert res2[i][0] < res2[i + 1][0]
            assert res2[i][1] < res2[i + 1][1]
        print(res1)
