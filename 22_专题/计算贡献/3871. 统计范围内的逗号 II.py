# 3871. 统计范围内的逗号 II-横看成岭侧成峰(从逗号的视角看，它出现在多少个数中)
# https://leetcode.cn/problems/count-commas-in-range-ii/description/
# 给你一个整数 n。

# 返回将所有从 [1, n]（包含两端）范围内的整数以 标准 数字格式书写时所用到的 逗号总数。

# 在 标准 格式中：

# 从右边开始，每 三位 数字后插入一个逗号。
# 位数 少于四位 的数字不包含逗号。


class Solution:
    def countCommas(self, n: int) -> int:
        res = 0
        low = 1000
        while low <= n:
            res += n - low + 1
            low *= 1000
        return res
