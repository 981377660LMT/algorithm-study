from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def minimumTime(self, time: List[int], totalTrips: int) -> int:
        def check(mid):
            res = 0
            for num in time:
                res += mid // num
            return res >= totalTrips

        # 上界直接开大一点 注意int(1e20)都可以 这道题上界是10**14
        left, right = 1, int(1e20)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left

