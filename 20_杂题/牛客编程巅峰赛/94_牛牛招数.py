# 找到一个大于a且为b的倍数的最小整数
class Solution:
    def findNumber(self, a, b):
        # write code here
        return a + b - a % b

