from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你三个整数 start ，finish 和 limit 。同时给你一个下标从 0 开始的字符串 s ，表示一个 正 整数。

# 如果一个 正 整数 x 末尾部分是 s （换句话说，s 是 x 的 后缀），且 x 中的每个数位至多是 limit ，那么我们称 x 是 强大的 。

# 请你返回区间 [start..finish] 内强大整数的 总数目 。


# 如果一个字符串 x 是 y 中某个下标开始（包括 0 ），到下标为 y.length - 1 结束的子字符串，那么我们称 x 是 y 的一个后缀。比方说，25 是 5125 的一个后缀，但不是 512 的后缀。


from functools import lru_cache


class Solution:
    def numberOfPowerfulInt(self, start: int, finish: int, limit: int, s: str) -> int:
        def cal(upper: int) -> int:
            @lru_cache(None)
            def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, s) -> int:
                """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界,出现次数为count"""
                if pos == n:
                    ok = 1 if not hasLeadingZero else 0
                    if ok:
                        print(s)
                    return ok

                res = 0
                up = min(nums[pos], limit) if isLimit else limit
                for cur in range(up + 1):
                    if sNums[pos] == -1 or cur == sNums[pos]:
                        if hasLeadingZero and cur == 0:
                            res += dfs(pos + 1, True, False, s + str(cur))
                        else:
                            res += dfs(pos + 1, False, (isLimit and cur == nums[pos]), s + str(cur))
                return res

            nums = list(map(int, str(upper)))

            sNums = [-1] * (len(nums) - len(s)) + list(map(int, s))
            n = len(nums)
            return dfs(0, True, True, "")

        if len(s) > len(str(finish)):
            return 0
        return cal(finish) - cal(start - 1)


# 1
# 971
# 9
# "17"
# print(Solution().numberOfPowerfulInt(1, 971, 9, "17"))
# start = 1, finish = 6000, limit = 4, s = "124"
# print(Solution().numberOfPowerfulInt(1, 6000, 4, "124"))
# start = 1000, finish = 2000, limit = 4, s = "3000"
print(Solution().numberOfPowerfulInt(1000, 2000, 4, "3000"))
