# 给你两个整数 num1 和 num2，表示一个 闭 区间 [num1, num2]。
#
# 一个数字的 波动值 定义为该数字中 峰 和 谷 的总数：
#
# 如果一个数位 严格大于 其两个相邻数位，则该数位为 峰。
# 如果一个数位 严格小于 其两个相邻数位，则该数位为 谷。
# 数字的第一个和最后一个数位 不能 是峰或谷。
# 任何少于 3 位的数字，其波动值均为 0。
# 返回范围 [num1, num2] 内所有数字的波动值之和。


from typing import Tuple
from functools import cache


class Solution:
    def totalWaviness(self, num1: int, num2: int) -> int:
        def f(num: int) -> int:
            s = str(num)
            n = len(s)

            @cache
            def dfs(
                pos: int, pre1: int, pre2: int, isLimit: bool, hasLeadingZero: bool
            ) -> Tuple[int, int]:
                """返回(构造出的数字个数，峰谷和)."""
                if pos == n:
                    return 1, 0
                res = 0
                count = 0
                upper = int(s[pos]) if isLimit else 9
                for d in range(upper + 1):
                    nIsLimit = isLimit and (d == upper)
                    nHasLeadingZero = hasLeadingZero and (d == 0)
                    hit = 0
                    if not hasLeadingZero and pre1 != -1 and pre2 != -1:
                        hit = (pre1 > pre2 and pre1 > d) or (pre1 < pre2 and pre1 < d)
                    nPre1 = -1
                    nPre2 = -1
                    if nHasLeadingZero:
                        nPre1 = -1
                        nPre2 = -1
                    else:
                        if hasLeadingZero:
                            nPre1 = d
                            nPre2 = -1
                        else:
                            nPre1 = d
                            nPre2 = pre1
                    nCount, nRes = dfs(pos + 1, nPre1, nPre2, nIsLimit, nHasLeadingZero)
                    count += nCount
                    res += nRes + (nCount * hit)
                return count, res

            res = dfs(0, -1, -1, True, True)
            dfs.cache_clear()
            return res[1]

        return f(num2) - f(num1 - 1)
