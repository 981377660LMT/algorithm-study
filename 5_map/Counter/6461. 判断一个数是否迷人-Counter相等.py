# 6461. 判断一个数是否迷人-Counter相等
# !如果经过以下修改得到的数字 恰好 包含数字 1 到 9 各一次且不包含任何 0 ，
# 那么我们称数字 n 是 迷人的 ：
# 将 n 与数字 2 * n 和 3 * n 连接 。


from collections import Counter


class Solution:
    def isFascinating(self, n: int) -> bool:
        cur = str(n) + str(n * 2) + str(n * 3)
        return Counter(cur) == Counter("123456789")
