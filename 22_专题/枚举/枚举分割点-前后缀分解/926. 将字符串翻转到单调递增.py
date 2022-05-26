# 给你一个二进制字符串 s，你可以将任何 0 翻转为 1 或者将 1 翻转为 0 。
# 返回使 s 单调递增的最小翻转次数。

# 最后肯定是前缀0+后缀1 枚举分割点即可


from itertools import accumulate


class Solution:
    def minFlipsMonoIncr(self, s: str) -> int:
        # 前缀1的个数
        L1 = list(accumulate(s, initial=0, func=lambda pre, cur: pre + int(cur == '1')))
        # 后缀0的个数
        R0 = list(accumulate(s[::-1], initial=0, func=lambda pre, cur: pre + int(cur == '0')))[::-1]
        return min(a + b for a, b in zip(L1, R0))


print(Solution().minFlipsMonoIncr("00110"))
