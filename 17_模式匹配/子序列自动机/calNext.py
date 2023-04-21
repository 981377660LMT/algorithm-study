from typing import List

INF = int(1e18)


def calNext(s: str) -> List[List[int]]:
    """子序列自动机.这里的字符集为小写字母.
    dp[i][c] := i 第i(0<=i<=n)个字符以后第一次出现字符c的位置(不存在的话为n).
    """
    n = len(s)
    dp = [[n] * 26 for _ in range(n + 1)]
    for i in range(n - 1, -1, -1):
        dp[i] = dp[i + 1][:]
        dp[i][ord(s[i]) - 97] = i
    return dp


if __name__ == "__main__":
    # E - Don't Be a Subsequence
    # https://atcoder.jp/contests/arc081/tasks/arc081_c
    # !求出字典序最小的字符串res使得res不是S的子序列
    # n<=2e5 都是小写字母
    # !字典序最小:从后往前

    S = input()
    n = len(S)
    next_ = calNext(S)
    dp = [INF] * (n + 1)  # 每个位置结尾的最小长度(dist)
    nextState = [("", n)] * (n + 1)  # dp复元
    dp[n] = 1
    for i in range(n - 1, -1, -1):
        for j in range(26):
            # 没有下一个字符
            if next_[i][j] == n:
                if dp[i] > 1:
                    dp[i] = 1
                    nextState[i] = (chr(j + 97), n)  # type: ignore
            # 有下一个字符
            elif dp[next_[i][j] + 1] + 1 < dp[i]:
                dp[i] = dp[next_[i][j] + 1] + 1
                nextState[i] = (chr(j + 97), next_[i][j] + 1)  # type: ignore

    res = []
    pos = 0
    while pos < n:
        res.append(nextState[pos][0])
        pos = nextState[pos][1]
    print("".join(res))
