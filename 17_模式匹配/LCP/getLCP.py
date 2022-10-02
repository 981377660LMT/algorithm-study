"""LCP"""
from typing import List


def getLCP(s: str) -> List[List[int]]:
    """O(n^2) dp 求解 LCP

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
