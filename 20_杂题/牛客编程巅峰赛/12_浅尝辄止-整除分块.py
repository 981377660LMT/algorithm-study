# 整除分块，维护left/right
class Solution:
    def work(self, n):
        # write code here
        res = 0
        left = 1
        while left <= n:
            div = n // left
            # 找到相同商数的区间右界
            right = n // div
            res += (right - left + 1) * div
            left = right + 1
        return res % 998244353
