"""0-n的全排列,给出长度为n字符串,只包含“D”和“I”,问有多少种排列方式,使得排列的数字大小满足给出的字符串。"""
# 1 <= S.length <= 3000
# DI序列的有效排列 是对整数 {0, 1, ..., n} 的一个排列


# !状态可以由(当前分配到第i个数，剩下的数中有j个比右端大) 确定
# `dp[i][j] // i番目まで考えて、右端より大きいのがj個ある時`
# 例えば右端が4で"<"が書いてあったとき、手持ちが{2,5,8,10}でも{3,11,22,33}でも
# 残り使える数字は3つ(5,8,10とか11,22,33とか)なので、そこからの遷移は同じ値になる。
# https://scrapbox.io/pocala-kyopro/T_-_Permutation
# https://qiita.com/Series_205/items/7d2c57b45179006d0bc6
# https://gyazo.com/0dd00309b3586016beacaa560ffb55ad/max_size/1000

from itertools import accumulate


MOD = int(1e9 + 7)


class Solution:
    def numPermsDISequence(self, s: str) -> int:
        """优化:时间复杂度O(n^2)"""
        n = len(s)
        dp = [1] * (n + 1)  # dp[j]表示剩下的数中有j个比右端大时的排列数

        for i in range(n):
            ndp, dpSum = [0] * (n + 1), [0] + list(accumulate(dp, lambda x, y: (x + y) % MOD))

            if s[i] == "I":
                # 用去一个大的 ndp[j]=dp[j+1]+dp[j+2]+...dp[n-i]
                for j in range(n - i):
                    ndp[j] = (dpSum[n - i + 1] - dpSum[j + 1]) % MOD
            else:
                # 用去一个小的 ndp[j]=dp[0]+dp[1]+...+dp[j]
                for j in range(n - i):
                    ndp[j] = dpSum[j + 1] % MOD

            dp = ndp

        return dp[0]


if __name__ == "__main__":

    print(Solution().numPermsDISequence("DID"))
    # 输出：5
    # 解释：
    # (0, 1, 2, 3) 的五个有效排列是：
    # (1, 0, 3, 2)
    # (2, 0, 3, 1)
    # (2, 1, 3, 0)
    # (3, 0, 2, 1)
    # (3, 1, 2, 0)
    # decrease increase

    n = int(input())
    s = input()
    print(Solution().numPermsDISequence("".join(["D" if char == ">" else "I" for char in s])))
