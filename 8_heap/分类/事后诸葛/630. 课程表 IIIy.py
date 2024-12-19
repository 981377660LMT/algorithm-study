from heapq import heappop, heappush
from typing import List


class Solution:
    def scheduleCourse(self, courses: List[List[int]]) -> int:
        """courses[i] = [durationi, lastDayi]
        Return the maximum number of courses that you can take.
        """
        courses.sort(key=lambda x: x[1])
        cost = 0
        pq = []
        for need, end in courses:
            cost += need
            heappush(pq, -need)
            if cost > end:
                cost += heappop(pq)
        return len(pq)


print(Solution().scheduleCourse([[100, 200], [200, 1300], [1000, 1250], [2000, 3200]]))
