from typing import List


class Solution:
    def distanceBetweenBusStops(self, distance: List[int], start: int, destination: int) -> int:
        return min(
            s1 := sum(distance[min(start, destination) : max(start, destination)]),
            sum(distance) - s1,
        )


# 输入：distance = [1,2,3,4], start = 0, destination = 1
# 输出：1
# 解释：公交站 0 和 1 之间的距离是 1 或 9，最小值是 1。
print(Solution().distanceBetweenBusStops([1, 2, 3, 4], 0, 1))
