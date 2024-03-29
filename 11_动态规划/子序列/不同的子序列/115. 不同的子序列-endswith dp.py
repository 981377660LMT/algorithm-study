# 求子序列个数


def countSubSequence(s: str, t: str) -> int:
    """求s中有多少个子序列为t,时间复杂度O(st)"""

    if not t:
        return 0

    dp = [0] * (len(t) + 1)  # endswith dp
    dp[0] = 1

    for i in range(len(s)):
        for j in reversed(range(len(t))):
            if s[i] == t[j]:
                dp[j + 1] += dp[j]

    return dp[-1]
