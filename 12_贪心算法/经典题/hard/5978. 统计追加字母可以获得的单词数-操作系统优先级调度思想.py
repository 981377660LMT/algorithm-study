from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits


MOD = int(1e9 + 7)
INF = 0x3F3F3F3F
EPS = int(1e-6)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]

# https://leetcode-cn.com/problems/earliest-possible-day-of-full-bloom/comments/1323899
# 把 播种时间 看作 进程使用cpu的时间 ，
# 把 生长时间 看作 I/O设备工作的时间，
# 由于I/O型设备的处理速度比CPU慢得多，只有让I/O设备尽早地进行工作，
# cpu才能调度其他进程执行，这样会提升系统的整体效率。
# cpu先调度I/O时间长的进程先执行，
# 才能在I/O设备工作的时候尽可能的调度其他程序执行，这样最终的时间会更短
# 对照计算机系统结构的流水线作业，种植过程不能并行，生长过程可以并行，让并行的时间最长，则总时间最短

# 1 <= n <= 105
# 1 <= plantTime[i], growTime[i] <= 104

# 贪心还是要考虑排序


class Solution:
    def earliestFullBloom(self, plantTime: List[int], growTime: List[int]) -> int:
        serial, res = 0, 0
        for cpu, io in sorted(zip(plantTime, growTime), key=lambda x: -x[1]):
            serial += cpu
            # 现在的串行+之后的并行 里的最大值
            res = max(res, serial + io)
        return res


# 9 9 2
print(Solution().earliestFullBloom(plantTime=[1, 4, 3], growTime=[2, 3, 1]))
print(Solution().earliestFullBloom(plantTime=[1, 2, 3, 2], growTime=[2, 1, 2, 1]))
print(Solution().earliestFullBloom(plantTime=[1], growTime=[1]))
print(
    Solution().earliestFullBloom(
        plantTime=[27, 5, 24, 17, 27, 4, 23, 16, 6, 26, 13, 17, 21, 3, 9, 10, 28, 26, 4, 10, 28, 2],
        growTime=[
            26,
            9,
            14,
            17,
            6,
            14,
            23,
            24,
            11,
            6,
            27,
            14,
            13,
            1,
            15,
            5,
            12,
            15,
            23,
            27,
            28,
            12,
        ],
    )
)
