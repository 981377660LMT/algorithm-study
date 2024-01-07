from re import A
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


def min2(a, b):
    return a if a < b else b


class Solution:
    def numberOfPowerfulInt(self, start: int, finish: int, limit: int, s: str) -> int:
        def cal(upper: int) -> int:
            @lru_cache(None)
            def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, realPos) -> int:
                """当前在第pos位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界"""
                if pos >= n:
                    ok = 1 if (not hasLeadingZero) else 0
                    return ok

                res = 0
                up = min2(nums[pos], limit) if isLimit else limit
                shouldMatch = pos >= offset
                for cur in range(up + 1):
                    if not shouldMatch:
                        if hasLeadingZero and cur == 0:
                            res += dfs(pos + 1, True, (isLimit and cur == nums[pos]), realPos)
                        else:
                            res += dfs(
                                pos + 1,
                                False,
                                (isLimit and cur == nums[pos]),
                                realPos,
                            )
                    else:
                        if cur == sNums[realPos]:
                            res += dfs(
                                pos + 1,
                                False,
                                (isLimit and cur == nums[pos]),
                                realPos + 1,
                            )
                return res

            if upper < int(s):
                return 0
            nums = list(map(int, str(upper)))
            sNums = list(map(int, s))
            n = len(nums)
            offset = max(0, len(nums) - len(s))
            return dfs(0, True, True, 0)

        if len(s) > len(str(finish)):
            return 0
        return cal(finish) - cal(start - 1)


# 1
# 971
# 9
# "17"
# print(Solution().numberOfPowerfulInt(1, 971, 9, "17"))
# # start = 1, finish = 6000, limit = 4, s = "124"
# print(Solution().numberOfPowerfulInt(1, 6000, 4, "124"))
# 33 111
# 82 1000 4 338
# 0
# 33 111
# -1
def check(start, finish, limit, s):
    res = 0
    for i in range(start, finish + 1):
        if str(i).endswith(s) and max(map(int, str(i))) <= limit:
            res += 1
    return res


# print(Solution().numberOfPowerfulInt(82, 1000, 4, "338"))

print(Solution().numberOfPowerfulInt(156, 229, 8, "7"))
print(check(156, 229, 8, "7"))
assert Solution().numberOfPowerfulInt(156, 229, 8, "7") == check(156, 229, 8, "7")
# 87 436 8 466
assert Solution().numberOfPowerfulInt(87, 436, 8, "466") == check(87, 436, 8, "466")
# 27 838 6 24
print(Solution().numberOfPowerfulInt(27, 838, 6, "24"))
print(check(27, 838, 6, "24"))
assert Solution().numberOfPowerfulInt(27, 838, 6, "24") == check(27, 838, 6, "24")
if __name__ == "__main__":
    import random

    for i in range(100):
        start = random.randint(1, 1000)
        finish = random.randint(start, 1000)
        limit = random.randint(1, 9)
        s = str(random.randint(1, 1000))
        if check(start, finish, limit, s) != Solution().numberOfPowerfulInt(
            start, finish, limit, s
        ):
            print(start, finish, limit, s)
            print(check(start, finish, limit, s))
            print(Solution().numberOfPowerfulInt(start, finish, limit, s))
            raise ValueError
