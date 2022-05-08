from itertools import groupby


MOD = int(1e9 + 7)
INF = int(1e20)

# 多写几遍groupBy
# 反思：题目没看清 该整数是 num 的一个`长度为 3` 的 子字符串
# groupby API 题


class Solution:
    def largestGoodInteger(self, num: str) -> str:
        groups = [[char, len(list(group))] for char, group in groupby(num)]
        return max((c * 3 for c, l in groups if l >= 3), default='')


print(Solution().largestGoodInteger(num="6777133339"))

print(Solution().largestGoodInteger("02844458683593127114444440593333336810269"))
