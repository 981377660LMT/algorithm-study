class Solution:
    def solve(self, x):
        """判断它是不是由某个完全平方数对1000取模得到的呢。"""
        for i in range(1, 1001):
            if pow(i, 2) % 1000 == x:
                return True
        return False
