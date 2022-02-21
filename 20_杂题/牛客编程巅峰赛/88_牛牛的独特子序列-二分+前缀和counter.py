# 求最长的形如 aaa...bbb...ccc... 的子串长度
# abc个数要一样
# n<=10^6
class Solution:
    def Maximumlength(self, x: str) -> int:
        # write code here
        """二分 `1/3的长度`"""
