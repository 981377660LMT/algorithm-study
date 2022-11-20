# !子序列dp(默认字符集为26个小写字母)
# https://nyaannyaan.github.io/library/string/subsequence-dp.hpp


# !1.最少连接s的多少个子序列可以得到t
# 建 s 的序列自动机然后把 t 丢进去转移，每次转移不动就回到初始点转移并让答案+1。

# !2.统计s的本质不同子序列个数 (记忆化搜索)


from typing import List, Tuple


def subSequnceDp(s: str) -> int:
    """子序列dp 示例:求s的本质不同子序列个数"""

    n = len(s)
    next = buildNext(s)
    dp = [0] * (n + 1)  # 前i个字符中的本质子序列个数
    dp[0] = 1
    for i in range(n):
        for j in range(26):
            if next[i][j] == n:
                continue
            dp[next[i][j] + 1] += dp[i]
    return sum(dp) - 1


def buildNext(s: str) -> List[Tuple[int, ...]]:
    n = len(s)
    next = [None] * n
    last = [n] * 26
    for i in range(n - 1, -1, -1):
        last[ord(s[i]) - 97] = i
        next[i] = tuple(last)  # type: ignore
    return next  # type: ignore


print(subSequnceDp("abc"))
