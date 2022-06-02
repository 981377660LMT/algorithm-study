from functools import lru_cache
from typing import Sequence, TypeVar


T = TypeVar('T', str, int)


# LCS模板
def LCS(seq1: Sequence[T], seq2: Sequence[T]) -> int:
    """返回LCS长度"""
    n1, n2 = len(seq1), len(seq2)
    res = 0
    dp = [[0] * (n2 + 1) for _ in range(n1 + 1)]
    for i in range(1, n1 + 1):
        for j in range(1, n2 + 1):
            if seq1[i - 1] == seq2[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1
                res = max(res, dp[i][j])
            else:
                dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])

    return res


def LCS2(seq1: Sequence[T], seq2: Sequence[T]) -> int:
    """返回LCS长度"""

    @lru_cache(None)
    def dfs(i: int, j: int) -> int:
        if i == len(seq1) or j == len(seq2):
            return 0
        if seq1[i] == seq2[j]:
            return dfs(i + 1, j + 1) + 1
        else:
            return max(dfs(i + 1, j), dfs(i, j + 1))

    return dfs(0, 0)


if __name__ == '__main__':
    assert LCS('abc', 'abcd') == 3
    assert LCS2('abc', 'abcd') == 3

