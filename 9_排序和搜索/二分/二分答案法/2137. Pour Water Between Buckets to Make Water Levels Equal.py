from typing import List

EPS = 1e-5
# You want to make the amount of water in each bucket equal.
# However, every time you pour k gallons of water, you spill loss percent of k.
# Return the maximum amount of water in each bucket after making the amount of water equal. Answers within 10-5 of the actual answer will be accepted.
class Solution:
    def equalizeWater(self, buckets: List[int], loss: int) -> float:
        def check(threshold: float) -> bool:
            """最后水位可以达到threshold"""
            more, less = 0, 0
            for num in buckets:
                if num > threshold:
                    more += num - threshold
                else:
                    less += threshold - num
            return more * ((100 - loss) / 100) >= less

        left, right = min(buckets), max(buckets)
        while left <= right:
            mid = (left + right) / 2
            if check(mid):
                left = mid + EPS
            else:
                right = mid - EPS
        return right


print(Solution().equalizeWater(buckets=[1, 2, 7], loss=80))
print(Solution().equalizeWater(buckets=[2, 4, 6], loss=50))
print(Solution().equalizeWater(buckets=[3, 3, 3, 3], loss=40))
