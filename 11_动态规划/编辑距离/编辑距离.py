# 给你两个单词 word1 和 word2， 请返回将 word1 转换成 word2 所使用的最少操作数  。
# 你可以对一个单词进行如下三种操作：
# 插入一个字符
# 删除一个字符
# 替换一个字符
# !O(n*m)


from typing import Any, Sequence


INF = int(1e18)


def min(a: int, b: int) -> int:
    return a if a < b else b


def editDistance(word1: Sequence[Any], word2: Sequence[Any]) -> int:
    n1, n2 = len(word1), len(word2)
    dp = [[INF] * (n2 + 1) for _ in range(n1 + 1)]
    dp[0][0] = 0

    for i in range(n1 + 1):
        dp[i][0] = i
    for j in range(n2 + 1):
        dp[0][j] = j

    for i in range(1, n1 + 1):
        for j in range(1, n2 + 1):
            if word1[i - 1] == word2[j - 1]:
                dp[i][j] = dp[i - 1][j - 1]
            else:
                dp[i][j] = min(dp[i - 1][j - 1] + 1, min(dp[i - 1][j] + 1, dp[i][j - 1] + 1))

    return dp[n1][n2]
