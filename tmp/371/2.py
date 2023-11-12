from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 、下标从 0 开始的二维字符串数组 access_times 。对于每个 i（0 <= i <= n - 1 ），access_times[i][0] 表示某位员工的姓名，access_times[i][1] 表示该员工的访问时间。access_times 中的所有条目都发生在同一天内。

# 访问时间用 四位 数字表示， 符合 24 小时制 ，例如 "0800" 或 "2250" 。

# 如果员工在 同一小时内 访问系统 三次或更多 ，则称其为 高访问 员工。

# 时间间隔正好相差一小时的时间 不 被视为同一小时内。例如，"0815" 和 "0915" 不属于同一小时内。

# 一天开始和结束时的访问时间不被计算为同一小时内。例如，"0005" 和 "2350" 不属于同一小时内。


# 以列表形式，按任意顺序，返回所有 高访问 员工的姓名。


class Solution:
    def findHighAccessEmployees(self, access_times: List[List[str]]) -> List[str]:
        def toTime(s: str) -> int:
            return int(s[:2]) * 60 + int(s[2:])

        mp = defaultdict(list)
        for name, time in access_times:
            mp[name].append(toTime(time))
        for name in mp:
            mp[name].sort()
        res = set()
        for name in mp:
            for i in range(len(mp[name]) - 2):
                if mp[name][i + 2] - mp[name][i] < 60:
                    res.add(name)
                    break
        return sorted(res)
