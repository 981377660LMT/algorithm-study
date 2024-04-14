from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s，表示一个 12 小时制的时间格式，其中一些数字（可能没有）被 "?" 替换。

# 12 小时制时间格式为 "HH:MM" ，其中 HH 的取值范围为 00 至 11，MM 的取值范围为 00 至 59。最早的时间为 00:00，最晚的时间为 11:59。

# 你需要将 s 中的 所有 "?" 字符替换为数字，使得结果字符串代表的时间是一个 有效 的 12 小时制时间，并且是可能的 最晚 时间。


def max2(a: int, b: int) -> int:
    return a if a > b else b


# 返回结果字符串。
class Solution:
    def findLatestTime(self, s: str) -> str:
        def dfs(index: int, path: str):
            if index == 5:
                hh = int(path[:2])
                mm = int(path[3:])
                if 0 <= hh <= 11 and 0 <= mm <= 59:
                    nonlocal res
                    res = max2(res, hh * 60 + mm)
                return
            if s[index] == "?":
                for i in range(10):
                    dfs(index + 1, path + str(i))
            else:
                dfs(index + 1, path + s[index])

        res = 0
        dfs(0, "")

        def time2hhmm(time: int) -> str:
            hh = time // 60
            mm = time % 60
            return f"{hh:02d}:{mm:02d}"

        return time2hhmm(res)
