# 爬楼梯 每次跳>=3节
# 求跳到S节的方案数 (S<=2e5)
# !dp[i]=dp[i-3]+(dp[i-4]+...dp[0])
# !即 dp[i]=dp[i-3]+dp[i-1]
# 注意初始项 [1,0,0,1]


from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    dp = defaultdict(int, {0: 1})
    dpSum = defaultdict(int, {0: 1})
    for i in range(1, n + 1):
        dp[i] = dpSum[i - 3] % MOD
        dpSum[i] = (dpSum[i - 1] + dp[i]) % MOD
    print(dp[n])

    # 推导式变形
    # n = int(input())
    # dp = [1, 0, 0, 1]
    # for i in range(4, n + 1):
    #     dp.append((dp[i - 3] + dp[i - 1]) % MOD)
    # print(dp[n])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
