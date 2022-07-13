# Distance Sequence
# 类似:
# !11_动态规划\dp优化\辅助数据结构dp\相邻差绝对值不大于k的最长子序列-平衡树优化dp.py


# 求相邻差绝对值大于等于k的序列个数
# !1<=ai<=M 且 abs(ai-ai+1)>=k
# n<=1000 M<=5000

# !dp[i][val] 表示 前i个数最后一个数为val的序列个数
# !那么dp可以表示为一段和 即可以前缀和dp

# !注意k=0时不要算重 复杂度O(nm)

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:
    n, m, k = map(int, input().split())
    dp, dpSum = [0] * (m + 5), [0] * (m + 5)
    for val in range(1, m + 1):
        dp[val] = 1
        dpSum[val] = dpSum[val - 1] + dp[val]

    for _ in range(1, n):
        ndp, ndpSum = [0] * (m + 5), [0] * (m + 5)
        for val in range(1, m + 1):
            # !注意k=0时不要算重
            if k == 0:
                ndp[val] = dpSum[m] % MOD
            else:
                leftUpper, rightLower = max(0, val - k), min(m + 1, val + k)
                ndp[val] = dpSum[leftUpper] + dpSum[m] - dpSum[rightLower - 1]
                ndp[val] %= MOD

        for val in range(1, m + 1):
            ndpSum[val] = ndpSum[val - 1] + ndp[val]
            ndpSum[val] %= MOD

        dp, dpSum = ndp, ndpSum

    print(dpSum[m])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
