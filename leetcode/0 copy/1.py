from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个正整数 days，表示员工可工作的总天数（从第 1 天开始）。另给你一个二维数组 meetings，长度为 n，其中 meetings[i] = [start_i, end_i] 表示第 i 次会议的开始和结束天数（包含首尾）。

# 返回员工可工作且没有安排会议的天数。


# 注意：会议时间可能会有重叠。


class Solution:
    def countDays(self, days: int, meetings: List[List[int]]) -> int:
        ...
