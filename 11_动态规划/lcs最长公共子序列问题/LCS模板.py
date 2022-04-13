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


if __name__ == '__main__':
    assert LCS('abc', 'abcd') == 3

