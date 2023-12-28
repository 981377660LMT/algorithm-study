# 1725 琪露诺-青蛙跳到对岸的最大得分
# https://www.luogu.com.cn/problem/P1725
# 开始时，琪露诺在编号 0 的格子上.
# !每次跳跃，琪露诺可以选择向前跳 left到right格.
# 只要她下一步的位置跳出了 [0,n] 的范围，她就会停止跳跃.
# !每个格子都有一个得分，琪露诺的得分就是她经过的格子的得分之和.
# 请问琪露诺最多可以得多少分.
#
# 思路:
# dp[i] 表示跳到i时的最大得分 (0<=i<=n+right)
# dp[i] = max(dp[j]+scores[i]) | i-right<=j<=i-left
# 答案为 max(dp[n],dp[n+1],...,dp[n+right])
#
# TODO: 最好还是用线段树


from MonoQueue import MonoQueue


from typing import List, Tuple

INF = int(1e18)


def maxJumpScore(scores: List[int], left: int, right: int) -> int:
    n = len(scores)
    queue = MonoQueue[Tuple[int, int]](lambda x, y: x[0] > y[0])  # (dp[i], i)
    dp = [-INF] * (n + right + 1)
    dp[0] = scores[0]
    queue.append((dp[0], 0))
    for i in range(left, n + right + 1):
        while queue and i - queue.head()[1] > right:
            queue.popleft()
        preMax = queue.head()[0] if queue else -INF
        curScore = scores[i] if i < n else 0
        dp[i] = max(dp[i], preMax + curScore)
        # !注意left边界范围
        queue.append((dp[i - (left - 1)], i - (left - 1)))
    return max(dp[len(scores) :])


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    _, left, right = map(int, input().split())
    scores = list(map(int, input().split()))
    print(maxJumpScore(scores, left, right))
