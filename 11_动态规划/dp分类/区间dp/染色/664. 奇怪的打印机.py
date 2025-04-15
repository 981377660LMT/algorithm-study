# 664. 奇怪的打印机
# https://leetcode.cn/problems/strange-printer/description/
#
# 有台奇怪的打印机有以下两个特殊要求：
#
# 打印机每次只能打印由 同一个字符 组成的序列。
# 每次可以在任意起始和结束位置打印新字符，并且会覆盖掉原来已有的字符。
# 给你一个字符串 s ，你的任务是计算这个打印机打印它需要的最少打印次数。
# 1 <= s.length <= 100

from functools import lru_cache
from itertools import groupby


class Solution:
    def strangePrinter(self, s: str) -> int:
        values = [a for a, _ in groupby(s)]

        @lru_cache(None)
        def dfs(start: int, end: int) -> int:
            if start >= end:
                return 0

            res = dfs(start + 1, end) + 1  # !直接移除这一段
            for mid in range(start + 1, end):
                if values[mid] == values[start]:
                    res = min(
                        res, dfs(start + 1, mid) + dfs(mid, end)
                    )  # !处理这一段后，start和mid合并
            return res

        return dfs(0, len(values))


print(Solution().strangePrinter("aaabbb"))
# 输出：2
# 解释：首先打印 "aaa" 然后打印 "bbb"。


# CABBA
#  CA | BB | A
# It was simply inserted with the cost of 1
# It was free from some previous step to the left that printed this character already (we can print extra character all the way till the end)
# 打一个新的要1快钱  cost = dfs(s[:-1]) + 1
# 打一个前面相同的不要钱   dfs(s[: i + 1]) + dfs(s[i + 1 : -1])
