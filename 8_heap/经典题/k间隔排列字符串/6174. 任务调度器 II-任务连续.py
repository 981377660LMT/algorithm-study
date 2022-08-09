from typing import List
from collections import defaultdict

INF = int(1e20)

# 6174. 任务调度器 II  模拟
# k间隔排列


class Solution:
    def taskSchedulerII(self, tasks: List[int], space: int) -> int:
        n = len(tasks)
        pre = defaultdict(lambda: -INF)
        day = 0
        res = n
        for num in tasks:
            diff = day - pre[num]
            if diff < space:
                res += space - diff  # 需要休息的天数
                day = pre[num] + space
            day += 1
            pre[num] = day
        return res


print(Solution().taskSchedulerII(tasks=[1, 2, 1, 2, 3, 1], space=3))
print(Solution().taskSchedulerII(tasks=[5, 8, 8, 5], space=2))

# class Solution:
#     def taskSchedulerII(self, tasks: List[int], space: int) -> int:
#         last = defaultdict(lambda : -inf)
#         cur = 0
#         for tp in tasks:
#             cur += 1
#             cur = max(cur, last[tp] + space + 1)
#             last[tp] = cur

#         return cur
