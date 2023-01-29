# 2481. 分割圆的最少切割次数
# !注意特判n=1的情况

# 给你一个整数 n ，请你返回将圆切割成相等的 n 等分的 最少 切割次数。
# 1 <= n <= 100


class Solution:
    def numberOfCuts(self, n: int) -> int:
        if n == 1:
            return 0
        return n if n & 1 else n // 2
