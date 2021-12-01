from typing import List

# 两位客户的空闲时间表：slots1 和 slots2，以及会议的预计持续时间 duration，请你为他们安排合适的会议时间
# 「会议时间」是两位客户都有空参加，并且持续时间能够满足预计时间 duration 的 最早的时间间隔。


class Solution:
    def minAvailableDuration(
        self, slots1: List[List[int]], slots2: List[List[int]], duration: int
    ) -> List[int]:
        slots1.sort()
        slots2.sort()

        n1, n2 = len(slots1), len(slots2)
        i, j = 0, 0

        while i < n1 and j < n2:
            s1, e1 = slots1[i]
            s2, e2 = slots2[j]
            s = max(s1, s2)
            e = min(e1, e2)
            if e - s >= duration:
                return [s, s + duration]
            if e1 < e2:
                i += 1
            else:
                j += 1

        return []


print(
    Solution().minAvailableDuration(
        slots1=[[10, 50], [60, 120], [140, 210]], slots2=[[0, 15], [60, 70]], duration=8
    )
)
# 输出：[60,68]
