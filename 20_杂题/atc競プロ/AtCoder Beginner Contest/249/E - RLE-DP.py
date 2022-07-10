"""
所有长为n的字符串中
RLE编码后比原来短的字符串有多少个
n<=3000 暗示O(n^2)的dp

dp[j][i] :i文字目まで決めて、RLEの長さがjの場合の文字列の通り数
遷移：ブロックを一つずつ追加
eg: 3 => 12 => 1 对应 +2 +3 +2 方案数对应 26*25*25
dp[j+getLRELen(blockLen)][i+blockLen]+=dp[j][i]*25 (i>=1,1≤k≤N,k表示新追加的bloack的长度)

1-9 => 2
10-99 => 3
100-999 => 4
1000-9999 => 5

这个dp是O(n^3)的
因为转移过来的原字符长度是连续的,所以可以用前缀和优化到O(n^2logn)
"""


import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
POW10 = [1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000]


def getLRELen(blockLen: int) -> int:
    """LRE编码长度"""
    return len(str(blockLen)) + 1


def main() -> None:
    n, MOD = map(int, input().split())
    if n <= 2:
        print(0)
        exit(0)

    dp = [[0] * (n + 5) for _ in [0] * (n + 5)]  # 前j个字符编码长为i的方案数
    dpSum = [[0] * (n + 5) for _ in [0] * (n + 5)]
    for j in range(1, n + 1):  # 只有一个block(一种字符)
        dp[getLRELen(j)][j] = 26
        dpSum[getLRELen(j)][j] = dpSum[getLRELen(j)][j - 1] + dp[getLRELen(j)][j]

    for j in range(1, n + 1):  # 编码长度
        for i in range(1, n + 1):  # 原串长度
            for len_ in range(2, 6):  # 当前block的编码长度,计算从哪些状态转移过来
                preJ = j - len_
                right, left = max(0, i - POW10[len_ - 2]), max(0, i - POW10[len_ - 1])
                if preJ < 0 or right < left:
                    continue
                dp[j][i] += (dpSum[preJ][right] - dpSum[preJ][left]) * 25
                dp[j][i] %= MOD
            dpSum[j][i] = dpSum[j][i - 1] + dp[j][i]
            dpSum[j][i] %= MOD

    res = 0
    for i in range(n):
        res += dp[i][n]
        res %= MOD
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
