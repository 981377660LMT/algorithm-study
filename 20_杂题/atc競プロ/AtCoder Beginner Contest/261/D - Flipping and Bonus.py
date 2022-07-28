# 高桥君扔n次硬币，刚开始的连胜次数为0
# 每一次扔硬币时，如果是正面的话，他的连胜次数就会加1，
# 他们获得当前连胜次数的奖励(可能有，也可能没有)
# 如果是负面，他的连胜次数就会清零。第i次如果胜了也有a的奖励

# m,counti<=n<=5000
# to jump or not to jump
# dp[index][count]

from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, m = map(int, input().split())
    scores = list(map(int, input().split()))
    bonus = defaultdict(int)
    for _ in range(m):
        count, score = map(int, input().split())
        bonus[count] += score

    dp = [0] * (n + 1)
    dp[1] = scores[0] + bonus[1]
    for i in range(1, n):
        ndp = [0] * (n + 1)

        # 中或不中
        for pre in range(i + 1):
            ndp[0] = max(ndp[0], dp[pre])
            ndp[pre + 1] = max(ndp[pre + 1], dp[pre] + scores[i] + bonus[pre + 1])

        dp = ndp

    print(max(dp, default=0))

    # memo = [[-1] * (n + 1) for _ in range(n + 1)]

    # def dfs(index: int, count: int) -> int:
    #     if index == n:
    #         return 0
    #     if memo[index][count] != -1:
    #         return memo[index][count]

    #     res = dfs(index + 1, 0)
    #     res = max(res, dfs(index + 1, count + 1) + scores[index] + bonus[count + 1])

    #     memo[index][count] = res
    #     return res

    # print(dfs(0, 0))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
