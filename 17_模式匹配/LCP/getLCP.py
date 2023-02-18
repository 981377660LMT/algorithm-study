"""s的任意`两个`后缀的LCP"""
from typing import List


def getLCP(s: str) -> List[List[int]]:
    """O(n^2) dp 求解 两个后缀的 LCP

    Args:
        s (str): 输入字符串
    Returns:
        List[List[int]]: LCP[i][j] 表示后缀 s[i:] 和 s[j:] 的最长公共前缀
    """
    n = len(s)
    lcp = [[0] * (n + 1) for _ in range(n + 1)]
    for i in range(n - 1, -1, -1):
        for j in range(n - 1, -1, -1):
            if s[i] == s[j]:
                lcp[i][j] = lcp[i + 1][j + 1] + 1
    return lcp


if __name__ == "__main__":
    # 求s的两个最长的不重合的子串(后缀)的长度 len(s)<=5e3
    # 枚举两个子串的起点
    # https://atcoder.jp/contests/abc141/tasks/abc141_e
    n = int(input())
    s = input()
    dp = getLCP(s)
    res = 0
    left1, left2 = -1, -1
    for i in range(n):
        for j in range(i + 1, n):
            cand = min(dp[i][j], j - i)
            if cand > res:
                res = cand
                left1, left2 = i, j
    print(res, s[left1 : left1 + res], s[left2 : left2 + res])
