from typing import Any, List, Sequence


def rangePalindrome(s: Sequence[Any]) -> List[List[int]]:
    """统计区间内回文串个数
    返回一个二维数组 dp, dp[i][j] 表示闭区间 [i,j] 内的回文串的个数
    https://codeforces.com/problemset/problem/245/H
    """
    n = len(s)
    dp = [[0] * n for _ in range(n)]
    for i in range(n):
        tmp = dp[i]
        tmp[i] = 1
        if i + 1 < n and s[i] == s[i + 1]:
            tmp[i + 1] = 1

    for i in range(n - 3, -1, -1):
        tmp1, tmp2 = dp[i], dp[i + 1]
        for j in range(i + 2, n):
            if s[i] == s[j]:
                tmp1[j] = tmp2[j - 1]

    # 到这里为止，dp[i][j] = 1 表示 s[i:j+1] 是回文串
    for i in range(n - 2, -1, -1):
        tmp1, tmp2 = dp[i], dp[i + 1]
        for j in range(i + 1, n):
            tmp1[j] += tmp1[j - 1] + tmp2[j] - tmp2[j - 1]

    return dp


if __name__ == "__main__":
    print(rangePalindrome("abab"))
